import { spawn } from "child_process";
import { existsSync, lstatSync, readdirSync, readFileSync } from "fs";
import path from "path";
import R2BlobService from "./r2Blob.js";
import { fileURLToPath } from "url";
import { Kafka, logLevel } from "kafkajs";
import KafkaProducerService from "./kafkaProducer.js";

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
const deployment_id = process.env.DEPLOYMENT_ID;


const kafkaClient = new Kafka({
    clientId: `docker-build-server-${deployment_id}`,
    brokers: [process.env.KAFKA_BROKERS],
    ssl: {
        ca: [readFileSync(path.join(__dirname, "ca.pem"), "utf-8")]
    },
    sasl: {
        mechanism: "plain",
        username: process.env.KAFKA_USERNAME,
        password: process.env.KAFKA_PASSWORD
    },
    logLevel: logLevel.DEBUG,
    connectionTimeout: 30000,
    requestTimeout: 30000
}
);

const kafkaProducer = new KafkaProducerService(kafkaClient);

const r2BlobService = new R2BlobService();

function runCommand(command, args, cwd) {
    return new Promise((resolve, reject) => {
        const child = spawn(command, args, {
            cwd,
            shell: true
        });

        // LIVE stdout log
        child.stdout.on("data", async (data) => {
            const text = data.toString();
            console.log(text);
            await kafkaProducer.generateMessage('mini-vercel-build-logs', { project_id, deployment_id }, text);
        });

        // LIVE stderr log
        child.stderr.on("data", async (data) => {
            const text = data.toString();
            console.error(text);
            await kafkaProducer.generateMessage('mini-vercel-build-logs', { project_id, deployment_id }, text);
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
    await kafkaProducer.generateMessage('mini-vercel-build-logs', { project_id, deployment_id }, "INFO: Running npm install...");

    await runCommand("npm", ["install"], outputPathDir);

    console.log("INFO: Running npm run build...");
    await kafkaProducer.generateMessage('mini-vercel-build-logs', { project_id, deployment_id }, "INFO: Running npm run build...");

    await runCommand("npm", ["run", "build"], outputPathDir);
}

async function uploadFiles() {
    const distFolderPath = path.join(outputPathDir, "dist");

    if (!existsSync(distFolderPath)) {
        throw new Error("dist folder does not exist. Build may have failed.");
    }

    const files = readdirSync(distFolderPath, { recursive: true });

    for (const file of files) {
        const filePath = path.join(distFolderPath, file);

        if (lstatSync(filePath).isDirectory()) continue;

        await r2BlobService.uploadToBlob(filePath, file, project_id);

        const msg = `Uploaded: ${file}`;
        console.log(msg);
        await kafkaProducer.generateMessage('mini-vercel-build-logs', { project_id, deployment_id }, msg);
    }
}

async function main() {
    await kafkaProducer.connect();
    try {
        await kafkaProducer.generateMessage('mini-vercel-build-logs', { project_id, deployment_id }, "INFO: Starting build pipeline...");
        console.log("INFO: Starting build pipeline...");

        await buildProject();

        await kafkaProducer.generateMessage('mini-vercel-build-logs', { project_id, deployment_id }, "INFO: Build completed. Uploading artifacts...");
        console.log("INFO: Build completed. Uploading artifacts...");

        await uploadFiles();

        await kafkaProducer.generateMessage('mini-vercel-build-logs', { project_id, deployment_id }, "INFO: Pipeline completed successfully.");
        console.log("INFO: Pipeline completed successfully.");

    } catch (err) {
        await kafkaProducer.generateMessage('mini-vercel-build-logs', { project_id, deployment_id }, `ERROR: ${err.message}, Pipeline failed.`);
        console.error(`ERROR: ${err.message}, Pipeline failed.`);
    } finally {
        await kafkaProducer.producer.disconnect();
        console.log("Kafka producer disconnected.");

        process.exit(0);
    }
}

main();
