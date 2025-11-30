const dotenv = require("dotenv");

dotenv.config();

module.exports = { host: process.env.CLICKHOUSE_HOST, database: process.env.CLICKHOUSE_DATABASE, user: process.env.CLICKHOUSE_USER, password: process.env.CLICKHOUSE_PASSWORD }
