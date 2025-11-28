import { exec, spawn } from "child_process";
import { existsSync, lstatSync, readdirSync } from "fs";
import path from "path";
import RedisService from "./redis.js";
import Redis from "ioredis";
import { BlobServiceClient } from "@azure/storage-blob";
import AzureBlobService from "./azureBlob.js";
import { fileURLToPath } from "url";

/* following functions need to be implemented:
1. cd into repo
2. run npm install
3. run npm build
4. push the build to azure blob storage
*/

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const outputPathDir = path.join(__dirname, "output");

const project_id = process.env.PROJECT_ID;


const redisClient = new Redis(process.env.REDIS_URL);

const redisService = new RedisService(redisClient);

const blobServiceClient = BlobServiceClient.fromConnectionString(process.env.AZURE_STORAGE_CONNECTION_STRING);

const azureBlobService = new AzureBlobService(blobServiceClient);

function runCommand(command, args, cwd) {
    return new Promise((resolve, reject) => {
        const child = spawn(command, args, {
            cwd,
            shell: true
        });

        // LIVE stdout log
        child.stdout.on("data", (data) => {
            const text = data.toString();
            console.log(text);
            redisService.publishLog(`logs:${project_id}`,text);
        });

        // LIVE stderr log
        child.stderr.on("data", (data) => {
            const text = data.toString();
            console.error(text);
            redisService.publishLog(`logs:${project_id}`,text);
        });

        child.on("close", (code) => {
            if (code !== 0) {
                reject(new Error(`${command} exited with code ${code}`));
            } else {
                resolve();
            }
        });
    });
}

async function buildProject() {
    console.log("INFO: Running npm install...");
    redisService.publishLog(`logs:${project_id}`,"INFO: Running npm install...");

    await runCommand("npm", ["install"], outputPathDir);

    console.log("INFO: Running npm run build...");
    redisService.publishLog(`logs:${project_id}`,"INFO: Running npm run build...");

    await runCommand("npm", ["run", "build"], outputPathDir);
}

async function uploadFiles() {
    const distFolderPath = path.join(outputPathDir, "dist");
    const projectId = process.env.PROJECT_ID;

    if (!existsSync(distFolderPath)) {
        throw new Error("dist folder does not exist. Build may have failed.");
    }

    const files = readdirSync(distFolderPath, { recursive: true });

    for (const file of files) {
        const filePath = path.join(distFolderPath, file);

        if (lstatSync(filePath).isDirectory()) continue;

        await azureBlobService.uploadToBlob(filePath, file, projectId);

        const msg = `Uploaded: ${file}`;
        console.log(msg);
        redisService.publishLog(`logs:${project_id}`,msg);
    }
}


async function main() {
    try {
        redisService.publishLog(`logs:${project_id}`,"INFO: Starting build pipeline...");
        console.log("INFO: Starting build pipeline...");

        await buildProject();

        redisService.publishLog(`logs:${project_id}`,"INFO: Build completed. Uploading artifacts...");
        console.log("INFO: Build completed. Uploading artifacts...");

        await uploadFiles();

        redisService.publishLog(`logs:${project_id}`,"INFO: Pipeline completed successfully.");
        console.log("INFO: Pipeline completed successfully.");

    } catch (err) {
        redisService.publishLog(`logs:${project_id}`,`ERROR: ${err.message}, Pipeline failed.`);
        console.error(`ERROR: ${err.message}, Pipeline failed.`);
    }
}

main();
