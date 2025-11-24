import { exec } from "child_process";
import { lstatSync, readdirSync } from "fs";
import path from "path";
import S3Service from "./s3";
import { S3Client } from "@aws-sdk/client-s3";
import RedisService from "./redis";
import Redis from "ioredis";

/* following functions need to be implemented:
1. cd into repo
2. run npm install
3. run npm build
4. push the build to s3
*/

const s3Client = new S3Client({
    region: process.env.AWS_REGION,
    credentials: {
        accessKeyId: process.env.AWS_ACCESS_KEY_ID,
        secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY
    }
})

const s3Service = new S3Service(s3Client);

const redisClient = new Redis(process.env.REDIS_URL);

const redisService = new RedisService(redisClient);

async function buildProject() {
    try {
        const outputPathDir = path.join(__dirname, "output");

        const p = exec(`cd ${outputPathDir} && npm install && npm run build`, (error, stdout, stderr) => {
            if (error) {
                redisService.publishLog(error.message);
                return;
            }
            if (stderr) {
                redisService.publishLog(stderr);
                return;
            }

            redisService.publishLog(stdout);
        });

        p.on("exit", async (code) => {
            const distFolderPath = path.join(outputPathDir, "dist");

            const distFolderContents = readdirSync(distFolderPath, { recursive: true });

            for (let i = 0; i < distFolderContents.length; i++) {
                const file = distFolderContents[i];
                const filePath = path.join(distFolderPath, file);
                if (lstatSync(filePath).isDirectory()) {
                    continue;
                }

                await s3Service.uploadToS3(filePath, file);

            }
        })
    } catch (error) {
        redisService.publishLog(error.message);
    }
}



async function main() {
    console.log("Building project...");
    await buildProject();
    console.log("Project built successfully");
}


main();
