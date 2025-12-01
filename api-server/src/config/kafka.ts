import dotenv from "dotenv";

dotenv.config();

export const kafkaConfig = {
    brokers: (process.env.KAFKA_BROKERS ?? "")
        .split(",")
        .map((b) => b.trim())
        .filter((b) => b.length > 0),

    clientId: process.env.KAFKA_CLIENT_ID ?? "",
};
