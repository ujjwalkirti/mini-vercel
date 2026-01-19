import dotenv from "dotenv";
import { logLevel } from "kafkajs";
import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import path from "path";

dotenv.config();
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export const kafkaConfig = {
    clientId: `api-server`,
    brokers: [process.env.KAFKA_BROKERS ?? ""],
    ssl: {
        ca: [readFileSync(path.join(__dirname, "ca.pem"), "utf-8")]
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
