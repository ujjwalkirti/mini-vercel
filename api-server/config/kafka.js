const dotenv = require("dotenv");

dotenv.config();

module.exports = { brokers: process.env.KAFKA_BROKERS.split(',').map(b => b.trim()), clientId: process.env.KAFKA_CLIENT_ID }
