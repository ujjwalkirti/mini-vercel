const dotenv = require("dotenv");

dotenv.config();

module.exports = {
    subscriptionId: process.env.AZURE_SUBSCRIPTION_ID,
    resourceGroup: process.env.AZURE_RESOURCE_GROUP,
    containerName: process.env.AZURE_CONTAINER_NAME,
    location: process.env.AZURE_LOCATION,
    image: process.env.AZURE_CONTAINER_IMAGE,
    osType: process.env.AZURE_OS_TYPE,
    dnsLabelName: process.env.AZURE_DNS_LABEL_NAME,
    acrServer: process.env.AZURE_ACR_SERVER,
    acrUsername: process.env.AZURE_ACR_USERNAME,
    acrPassword: process.env.AZURE_ACR_PASSWORD,
    storageConnectionString: process.env.AZURE_STORAGE_CONNECTION_STRING,
}
