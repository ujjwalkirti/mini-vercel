const { DefaultAzureCredential } = require("@azure/identity");
const axios = require("axios");

class AzureACIServiceREST {
    constructor(subscriptionId, location, osType, containerGroupName, image, acrServer, acrUsername, acrPassword, dnsLabel) {
        this.subscriptionId = subscriptionId;
        this.location = location;
        this.osType = osType;
        this.containerGroupName = containerGroupName;
        this.image = image;
        this.acrServer = acrServer;
        this.acrUsername = acrUsername;
        this.acrPassword = acrPassword;
        this.dnsLabel = dnsLabel;
        this.credential = new DefaultAzureCredential();
    }

    async getToken() {
        const token = await this.credential.getToken("https://management.azure.com/.default");
        return token.token;
    }

    async startACI(envVars, resourceGroup) {
        const token = await this.getToken();

        const url = `https://management.azure.com/subscriptions/${this.subscriptionId}/resourceGroups/${resourceGroup}/providers/Microsoft.ContainerInstance/containerGroups/${this.containerGroupName}?api-version=2023-05-01`;

        const body = {
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
        };

        const response = await axios.put(url, body, {
            headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "application/json"
            },
        });
        if(response.status === 200 || response.status === 201) {
            return response.data;
        } else {
            console.error("PUT ERROR:", response.data);
            throw new Error("Failed to start ACI");
        }
    }

    async deleteACI(resourceGroup) {
        const token = await this.getToken();
        const url = `https://management.azure.com/subscriptions/${this.subscriptionId}/resourceGroups/${resourceGroup}/providers/Microsoft.ContainerInstance/containerGroups/${this.containerGroupName}?api-version=2023-05-01`;

        try {
            const response = await axios.delete(url, {
                headers: { Authorization: `Bearer ${token}` }
            });

            return response.data;
        } catch (err) {
            console.error("DELETE ERROR:", err.response?.data || err.message);
            throw err;
        }
    }

    async getLogs(resourceGroup, containerName = this.containerGroupName) {
        const token = await this.getToken();
        const url = `https://management.azure.com/subscriptions/${this.subscriptionId}/resourceGroups/${resourceGroup}/providers/Microsoft.ContainerInstance/containerGroups/${this.containerGroupName}/containers/${containerName}/logs?api-version=2023-05-01`;

        try {
            const response = await axios.get(url, {
                headers: { Authorization: `Bearer ${token}` }
            });

            return response.data;
        } catch (err) {
            console.error("LOGS ERROR:", err.response?.data || err.message);
            throw err;
        }
    }
}

module.exports = AzureACIServiceREST;
