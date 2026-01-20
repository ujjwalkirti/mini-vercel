import dotenv from "dotenv";

dotenv.config();

export const clickhouseConfig = {
    url: `https://${process.env.CLICKHOUSE_USERNAME}:${process.env.CLICKHOUSE_PASSWORD}@${process.env.CLICKHOUSE_HOST}:${process.env.CLICKHOUSE_PORT ?? '8443'}/${process.env.CLICKHOUSE_DATABASE}`
};
