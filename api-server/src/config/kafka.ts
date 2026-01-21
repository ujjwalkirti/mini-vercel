import dotenv from "dotenv";
import { logLevel } from "kafkajs";
import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import path from "path";

dotenv.config();
const __filename = fileURLToPath(import.meta.url);

const kafkaCA = readFileSync(
    "/secrets/kafka-consumer-ca",
    "utf-8"
);

export const kafkaConfig = {
    clientId: `api-server`,
    brokers: [process.env.KAFKA_BROKERS ?? ""],
    ssl: {
        ca: [kafkaCA],
    },
    sasl: {
        mechanism: "plain" as const,
        username: process.env.KAFKA_USERNAME ?? "",
        password: process.env.KAFKA_PASSWORD ?? "",
    },
    logLevel: logLevel.ERROR,
    connectionTimeout: 30000,
    requestTimeout: 30000
}
