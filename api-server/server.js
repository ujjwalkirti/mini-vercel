const express = require("express");
const dotenv = require("dotenv");
const AzureACIServiceREST = require("./azureACI");
const { Server } = require("socket.io");
const azureACIConfig = require("./config/azure.js");
const { Kafka } = require('kafkajs');
const KafkaConsumerService = require("./kafkaConsumer.js");
dotenv.config();

const app = express();
const PORT = 9000;

const kafkaClient = new Kafka({
    brokers: [process.env.KAFKA_BROKERS],
    clientId: process.env.KAFKA_CLIENT_ID
});

const kafkaConsumer = new KafkaConsumerService(kafkaClient, "mini-vercel-build-logs");

const io = new Server({
    cors: '*',
});

io.on('connection', (socket) => {
    socket.on('subscribe', (channel) => {
        socket.join(channel);
        socket.emit('message', `Joined channel: ${channel}`);
    });
});

io.listen(9001, () => console.log("Socket server listening on port 9001"));

app.use(express.json());

const azureACIService = new AzureACIServiceREST(azureACIConfig);

app.post("/build", async (req, res) => {
    const { github_url, project_id } = req.body;

    if (!github_url) return res.status(400).send({ error: "Missing github_url" });

    try {
        const envVars = [
            { name: "PROJECT_ID", value: project_id },
            { name: "GIT_REPOSITORY_URL", value: github_url },
            { name: "AZURE_STORAGE_CONNECTION_STRING", value: storageConnectionString },
            { name: "KAFKA_BROKERS", value: process.env.KAFKA_BROKERS },
            { name: "KAFKA_CLIENT_ID", value: process.env.KAFKA_CLIENT_ID },
        ];

        await azureACIService.startACI(envVars, resourceGroup);

        res.status(200).send({
            status: "Queued", message: "Build queued successfully", data: {
                "url": `${project_id}.localhost:8000`
            }
        });
    } catch (error) {
        console.error(error);
        res.status(500).send({ error: error.message || "Internal Server Error" });
    }
});

kafkaConsumer.listenForMessagesInBatch('mini-vercel-build-logs', async (message) => {
    const { key, value } = message;
    if (!key || !value) return;
    try {
        const data = JSON.parse(value);
        io.to(key).emit('message', data);
    } catch (error) {
        console.error(error);
    }
})

app.listen(PORT, () => console.log(`API server listening on port ${PORT}`));
