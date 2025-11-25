import { createReadStream } from "fs";

class AzureBlobService {
    constructor(blobServiceClient) {
        this.blobServiceClient = blobServiceClient;
    }

    async uploadToBlob(filePath, file, projectId) {
        const blobPath = `${projectId}/${file}`;
        const containerClient = this.blobServiceClient.getContainerClient("build-outputs");
        const blockBlobClient = containerClient.getBlockBlobClient(blobPath);
        const fileStream = createReadStream(filePath);
        await blockBlobClient.uploadStream(fileStream);
    }
}

export default AzureBlobService;
