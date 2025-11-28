const express = require("express");
const dotenv = require("dotenv");
const AzureACIServiceREST = require("./azureACI");
const { Server } = require("socket.io");
const RedisService = require("./redis");
const Redis = require("ioredis");

dotenv.config();

const app = express();
const PORT = 9000;

const redisSubscriber = new Redis(process.env.REDIS_URL);
const redisService = new RedisService(redisSubscriber);

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

const subscriptionId = process.env.AZURE_SUBSCRIPTION_ID;
const resourceGroup = process.env.AZURE_RESOURCE_GROUP;
const containerName = process.env.AZURE_CONTAINER_NAME;
const location = process.env.AZURE_LOCATION;
const image = process.env.AZURE_CONTAINER_IMAGE;
const osType = process.env.AZURE_OS_TYPE;
const dnsLabelName = process.env.AZURE_DNS_LABEL_NAME;
const acrServer = process.env.AZURE_ACR_SERVER;
const acrUsername = process.env.AZURE_ACR_USERNAME;
const acrPassword = process.env.AZURE_ACR_PASSWORD;
const storageConnectionString = process.env.AZURE_STORAGE_CONNECTION_STRING;

const azureACIService = new AzureACIServiceREST(subscriptionId, location, osType, containerName, image, acrServer, acrUsername, acrPassword, dnsLabelName);

const redisUrl = process.env.REDIS_URL;

app.post("/build", async (req, res) => {
    const { github_url, project_id } = req.body;

    if (!github_url) return res.status(400).send({ error: "Missing github_url" });

    try {
        const envVars = [{ name: "PROJECT_ID", value: project_id }, { name: "GIT_REPOSITORY_URL", value: github_url }, {
            name: "AZURE_STORAGE_CONNECTION_STRING",
            value: storageConnectionString
        }, { name: "REDIS_URL", value: redisUrl }];

        const response = await azureACIService.startACI(envVars, resourceGroup);

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

redisService.subscribeLog((pattern, channel, message) => { io.to(channel).emit('message', message); });

app.listen(PORT, () => console.log(`API server listening on port ${PORT}`));
