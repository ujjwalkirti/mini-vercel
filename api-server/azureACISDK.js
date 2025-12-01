const { ContainerInstanceManagementClient } = require("@azure/arm-containerinstance");
const { DefaultAzureCredential } = require("@azure/identity");

class AzureACIServiceSDK {
    /**
     * Constructs a new AzureACIServiceSDK instance.
     * @param {string} subscriptionId - The subscription ID of the Azure account.
     * @param {string} location - The location of the container instance.
     * @param {string} osType - The type of the operating system of the container instance.
     * @param {string} containerGroupName - The name of the container group.
     * @param {string} image - The image of the container instance.
     * @param {string} acrServer - The server URL of the Azure container registry.
     * @param {string} acrUsername - The username of the Azure container registry.
     * @param {string} acrPassword - The password of the Azure container registry.
     * @param {string} dnsLabel - The DNS label of the container instance.
     */
    constructor({ subscriptionId, location, osType, containerName, image, acrServer, acrUsername, acrPassword, dnsLabel }) {
        this.subscriptionId = subscriptionId;
        this.location = location;
        this.osType = osType;
        this.containerGroupName = containerName;
        this.image = image;
        this.acrServer = acrServer;
        this.acrUsername = acrUsername;
        this.acrPassword = acrPassword;
        this.dnsLabel = dnsLabel;
        this.containerInstanceManagementClient = new ContainerInstanceManagementClient(new DefaultAzureCredential(), subscriptionId);
    }

    async startACI(envVars, resourceGroup) {
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
                            environmentVariables: envVars || [],
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
        }

        try {
            await this.containerInstanceManagementClient.containerGroups.beginCreateOrUpdateAndWait(resourceGroup, this.containerGroupName, containerGroup);
        } catch (error) {
            console.error("ERROR: ", error.response?.data || error.message);
            throw new Error(error.response?.data || error.message || "Failed to create container group");
        }
    }
}

module.exports = AzureACIServiceSDK;
