class AzureBlobService {
    constructor(blobServiceClient) {
        this.blobServiceClient = blobServiceClient;
    }

    async uploadToBlob(filePath, file) {
        const blockBlobClient = this.blobServiceClient.getBlockBlobClient(file);
        const fileStream = createReadStream(filePath);
        await blockBlobClient.uploadStream(fileStream);
    }
}

export default AzureBlobService;
