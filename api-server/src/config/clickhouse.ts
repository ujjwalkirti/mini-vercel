import dotenv from "dotenv";

dotenv.config();

export const clickhouseConfig = {
    url: process.env.CLICKHOUSE_HOST ?? "",
    username: process.env.CLICKHOUSE_USERNAME ?? "",
    password: process.env.CLICKHOUSE_PASSWORD ?? "",
};
