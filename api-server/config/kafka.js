const dotenv = require("dotenv");

dotenv.config();

module.exports = { brokers: [process.env.KAFKA_BROKERS], clientId: process.env.KAFKA_CLIENT_ID }
