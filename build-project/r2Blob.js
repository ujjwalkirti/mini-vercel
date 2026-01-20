import { S3Client, PutObjectCommand } from "@aws-sdk/client-s3";
import { createReadStream } from "fs";
import mime from "mime-types";

class R2BlobService {
    constructor() {
        this.s3Client = new S3Client({
            region: "auto",
            endpoint: `https://${process.env.R2_ACCOUNT_ID}.r2.cloudflarestorage.com`,
            credentials: {
                accessKeyId: process.env.R2_ACCESS_KEY_ID,
                secretAccessKey: process.env.R2_SECRET_ACCESS_KEY,
            },
        });
        this.bucketName = process.env.R2_BUCKET_NAME;
    }

    async uploadToBlob(filePath, file, projectId) {
        const key = `${projectId}/${file}`;
        const fileStream = createReadStream(filePath);
        const contentType = mime.lookup(filePath) || "application/octet-stream";

        const command = new PutObjectCommand({
            Bucket: this.bucketName,
            Key: key,
            Body: fileStream,
            ContentType: contentType,
        });

        await this.s3Client.send(command);
    }
}

export default R2BlobService;
