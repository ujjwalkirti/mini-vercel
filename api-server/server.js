const express = require("express");
const dotenv = require("dotenv");
const AzureACIServiceREST = require("./azureACIREST.js");
const { Server } = require("socket.io");
const azureACIConfig = require("./config/azure.js");
// const { Kafka } = require('kafkajs');
// const KafkaConsumerService = require("./kafkaConsumer.js");
// const kafkaConfig = require("./config/kafka.js");
const { createClient } = require("@clickhouse/client");
const ClickHouseService = require("./clickhouse.js");
const clickHouseConfig = require("./config/clickhouse.js");
// const { PrismaClient } = require("@prisma/client");
const { generateSlug } = require("random-word-slugs");
const AzureACIServiceSDK = require("./azureACISDK.js");


dotenv.config();

const app = express();
const PORT = 9000;

// const kafkaClient = new Kafka(kafkaConfig);

// const kafkaConsumer = new KafkaConsumerService(kafkaClient, "mini-vercel-build-logs");

const io = new Server({
    cors: '*',
});

io.on('connection', (socket) => {
    socket.on('subscribe', (channel) => {
        socket.join(channel);
        socket.emit('message', `Joined channel: ${channel}`);
    });
});

io.listen(9090, () => console.log("Socket server listening on port 9001"));

const clickHouseClient = createClient(clickHouseConfig);

const clickHouseService = new ClickHouseService(clickHouseClient);

// const prismaClient = new PrismaClient({
//     datasources: {
//         db: {
//             url: process.env.DATABASE_URL
//         }
//     }
// });

app.use(express.json());

console.log(azureACIConfig);

// const azureACIServiceREST = new AzureACIServiceREST(azureACIConfig);
const azureACIServiceSDK = new AzureACIServiceSDK(azureACIConfig);

app.post('/add-project', async (req, res) => {
    const { name, github_url } = req.body;
    if (!name || !github_url) return res.status(400).send({ success: false, message: "Missing name or github_url", error: "Missing name or github_url" });

    try {
        // const project = await prismaClient.project.create({
        //     data: {
        //         name: name,
        //         git_url: github_url,
        //         subdomain: generateSlug(3)
        //     }
        // });
        res.status(201).json({
            success: true,
            message: "Project created successfully",
            data: project
        });
    } catch (error) {
        console.error(error);
        res.status(500).send({ success: false, message: "Internal Server Error", error: error.message || "Internal Server Error" });
    }
})

app.post("/deploy", async (req, res) => {
    const { project_id } = req.body;

    if (!project_id) return res.status(400).send({ error: "Missing github_url" });

    // const project = await prismaClient.project.findUnique({
    //     where: {
    //         id: project_id
    //     }
    // });

    // if (!project) return res.status(404).send({ error: "Project not found" });

    try {
        // create a deployment record
        // const deployment = await prismaClient.deployement.create({
        //     data: {
        //         project: {
        //             connect: {
        //                 id: project_id
        //             }
        //         },
        //         status: "QUEUED"
        //     }
        // });

        const envVars = [
            { name: "PROJECT_ID", value: project_id },
            { name: "GIT_REPOSITORY_URL", value: process.env.GIT_REPOSITORY_URL },
            { name: "AZURE_STORAGE_CONNECTION_STRING", value: azureACIConfig.storageConnectionString },
            { name: "KAFKA_BROKERS", value: process.env.KAFKA_BROKERS },
            { name: "KAFKA_CLIENT_ID", value: process.env.KAFKA_CLIENT_ID },
            // { name: "DEPLOYMENT_ID", value: deployment.id },
        ];

        // await azureACIServiceREST.startACI(envVars, resourceGroup);

        await azureACIServiceSDK.startACI(envVars, process.env.AZURE_RESOURCE_GROUP);

        res.status(200).send({
            success: true, message: "Build queued successfully", data: {
                status: "Queued",
                "url": `${project_id}.localhost:8000`
            }
        });
    } catch (error) {
        console.error(error);
        res.status(500).send({ success: false, message: "Internal Server Error", error: error.message || "Internal Server Error" });
    }
});

// kafkaConsumer.listenForMessagesInBatch('mini-vercel-build-logs', async (message) => {
//     const { key, value } = message;
//     if (!key || !value) return;
//     try {
//         const data = JSON.parse(value);
//         io.to(key).emit('message', data);
//     } catch (error) {
//         console.error(error);
//     }
// })

app.listen(PORT, () => console.log(`API server listening on port ${PORT}`));
