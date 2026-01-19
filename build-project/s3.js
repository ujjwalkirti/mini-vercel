import { PutObjectCommand } from "@aws-sdk/client-s3";
import { createReadStream } from "fs";
import { lookup } from "mime-types";


class S3Service {
    constructor(s3Client) {
        this.s3Client = s3Client
    }
    project_id = process.env.PROJECT_ID;
    deployment_id = process.env.DEPLOYMENT_ID

    async uploadToS3(filePath, file) {
        const command = new PutObjectCommand({
            Bucket: 'mini-vercel-outputs',
            Key: `__outputs/${project_id}/${deployment_id}/${file}`,
            Body: createReadStream(filePath),
            ContentType: lookup(filePath) || 'application/octet-stream'
        })

        await this.s3Client.send(command)
    }
}

export default S3Service;
