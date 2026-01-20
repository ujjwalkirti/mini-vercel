import { ContainerInstanceManagementClient } from "@azure/arm-containerinstance";
import { DefaultAzureCredential } from "@azure/identity";

export interface AzureACISDKConfig {
    subscriptionId: string;
    location: string;
    osType: string;
    containerName: string;
    image: string;
    acrServer: string;
    acrUsername: string;
    acrPassword: string;
    dnsLabel: string;
}

export interface EnvVar {
    name: string;
    value: string;
}

export default class AzureACIServiceSDK {
    private subscriptionId: string;
    private location: string;
    private osType: string;
    private containerGroupName: string;
    private image: string;
    private acrServer: string;
    private acrUsername: string;
    private acrPassword: string;
    private dnsLabel: string;

    private containerInstanceManagementClient: ContainerInstanceManagementClient;

    constructor({
        subscriptionId,
        location,
        osType,
        containerName,
        image,
        acrServer,
        acrUsername,
        acrPassword,
        dnsLabel
    }: AzureACISDKConfig) {
        this.subscriptionId = subscriptionId;
        this.location = location;
        this.osType = osType;
        this.containerGroupName = containerName;
        this.image = image;
        this.acrServer = acrServer;
        this.acrUsername = acrUsername;
        this.acrPassword = acrPassword;
        this.dnsLabel = dnsLabel;

        this.containerInstanceManagementClient = new ContainerInstanceManagementClient(
            new DefaultAzureCredential(),
            subscriptionId
        );
    }

    async startACI(envVars: EnvVar[] = [], resourceGroup: string): Promise<void> {
        const containerGroup = {
            location: this.location,
            properties: {
                osType: this.osType,
                sku: "Confidential",
                restartPolicy: "Never",

                confidentialComputeProperties: {
                    ccePolicy:
                        "eyJhbGxvd19hbGwiOiB0cnVlLCAiY29udGFpbmVycyI6IHsibGVuZ3RoIjogMCwgImVsZW1lbnRzIjogbnVsbH19"
                },

                ipAddress: {
                    type: "Public",
                    ports: [
                        {
                            protocol: "TCP",
                            port: 80
                        }
                    ],
                    dnsNameLabel: this.dnsLabel || undefined
                },

                containers: [
                    {
                        name: this.containerGroupName,
                        properties: {
                            image: this.image,
                            command: [],
                            environmentVariables: envVars,
                            ports: [{ port: 80 }],
                            resources: {
                                requests: {
                                    cpu: 1,
                                    memoryInGB: 1.5
                                }
                            }
                        }
                    }
                ],

                imageRegistryCredentials: [
                    {
                        server: this.acrServer,
                        username: this.acrUsername,
                        password: this.acrPassword
                    }
                ]
            }
        };

        try {
            await this.containerInstanceManagementClient.containerGroups.beginCreateOrUpdateAndWait(
                resourceGroup,
                this.containerGroupName,
                containerGroup as any
            );
        } catch (error: any) {
            console.error("ERROR:", error.response?.data || error.message);
            throw new Error(error.response?.data || error.message || "Failed to create container group");
        }
    }
}
