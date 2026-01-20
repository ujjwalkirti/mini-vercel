import express from "express";
import dotenv from "dotenv";
import { Server } from "socket.io";
import { Kafka } from "kafkajs";
import KafkaConsumerService from "./kafkaConsumer";
import { createClient } from "@clickhouse/client";
import ClickHouseService from "./clickhouse";
import { generateSlug } from "random-word-slugs";
import { clickhouseConfig } from "./config/clickhouse";
import { kafkaConfig } from "./config/kafka";
import { prismaClient } from "./lib/prisma";
import { ECSClient } from "@aws-sdk/client-ecs";
import { awsConfig } from "./config/aws";
import AWSECSService from "./awsECS";
import { ecsConfig } from "./config/ecs";

dotenv.config();

const app = express();
const PORT = 9000;

const kafkaClient = new Kafka(kafkaConfig);

const kafkaConsumer = new KafkaConsumerService(kafkaClient, "mini-vercel-build-logs");

const ecsClient = new ECSClient({
    credentials: awsConfig
});

const awsECSService = new AWSECSService(ecsClient);

const io = new Server({
    cors: "*",
} as any);

io.on("connection", (socket) => {
    socket.on("subscribe", (channel: string) => {
        socket.join(channel);
        socket.emit("message", `Joined channel: ${channel}`);
    });
});

io.listen(9090, (() => {
    console.log("Socket server listening on port 9001");
}) as any);

console.log(clickhouseConfig.url);
const clickHouseClient = createClient(clickhouseConfig);

const clickHouseService = new ClickHouseService(clickHouseClient);

app.use(express.json());

app.post("/add-project", async (req, res) => {
    const { name, github_url } = req.body;
    if (!name || !github_url)
        return res
            .status(400)
            .send({
                success: false,
                message: "Missing name or github_url",
                error: "Missing name or github_url",
            });

    try {
        const project = await prismaClient.project.create({
            data: {
                name: name,
                gitURL: github_url,
                subDomain: generateSlug(3)
            }
        });
        res.status(201).json({
            success: true,
            message: "Project created successfully",
            data: project,
        });
    } catch (error: any) {
        console.error(error);
        res
            .status(500)
            .send({
                success: false,
                message: "Internal Server Error",
                error: error.message || "Internal Server Error",
            });
    }
});

app.post("/deploy", async (req, res) => {
    const { project_id } = req.body;

    if (!project_id)
        return res.status(400).send({ error: "Missing github_url" });

    const project = await prismaClient.project.findUnique({
        where: {
            id: project_id
        }
    });

    if (!project) return res.status(404).send({ error: "Project not found" });

    try {
        const deployment = await prismaClient.deployment.create({
            data: {
                status: "QUEUED",
                project: { connect: { id: project_id } }
            }
        })

        const envVars = [
            { name: "PROJECT_ID", value: project_id },
            { name: "GIT_REPOSITORY_URL", value: project.gitURL },
            { name: "KAFKA_BROKERS", value: process.env.KAFKA_BROKERS },
            { name: "KAFKA_CLIENT_ID", value: process.env.KAFKA_CLIENT_ID },
            { name: "KAFKA_USERNAME", value: process.env.KAFKA_USERNAME },
            { name: "KAFKA_PASSWORD", value: process.env.KAFKA_PASSWORD },
            { name: "R2_ACCOUNT_ID", value: process.env.R2_ACCOUNT_ID },
            { name: "R2_ACCESS_KEY_ID", value: process.env.R2_ACCESS_KEY_ID },
            { name: "R2_SECRET_ACCESS_KEY", value: process.env.R2_SECRET_ACCESS_KEY },
            { name: "R2_BUCKET_NAME", value: process.env.R2_BUCKET_NAME },
            { name: "DEPLOYMENT_ID", value: deployment.id },
        ];

        const ecsTaskProps = {
            cluster: ecsConfig.cluster,
            taskDefinition: ecsConfig.taskDefinition,
            image: ecsConfig.imageName,
            envVars: envVars,
            subnets: ecsConfig.subnets,
            securityGroups: ecsConfig.securityGroups,
            assignPublicIp: ecsConfig.assignPublicIp,
            launchType: ecsConfig.launchType,
            count: ecsConfig.count
        }

        await awsECSService.runTask(ecsTaskProps);

        res.status(200).send({
            success: true,
            message: "Build queued successfully",
            data: {
                status: "Queued",
                url: `${project_id}.localhost:8001`,
            },
        });
    } catch (error: any) {
        console.error(error);
        res
            .status(500)
            .send({
                success: false,
                message: "Internal Server Error",
                error: error.message || "Internal Server Error",
            });
    }
});

kafkaConsumer.listenForMessagesInBatch('mini-vercel-build-logs', async (message) => {
    const { key, value } = message;
    if (!key || !value) return;
    try {
        const { project_id, deployment_id, log } = JSON.parse(value.toString());
        console.log({ project_id, deployment_id, log });
        const { query_id } = await clickHouseService.insertLog('log_events', { deployment_id, log });

        console.log(query_id)
    } catch (error) {
        console.error(error);
    }
})

app.listen(PORT, () =>
    console.log(`API server listening on port ${PORT}`)
);
