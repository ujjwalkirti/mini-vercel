import express from 'express';
import dotenv from 'dotenv';
import { Kafka } from 'kafkajs';
import { createClient } from '@clickhouse/client';
import { ECSClient } from '@aws-sdk/client-ecs';

import { createRouter } from './routes/index.js';
import {
    corsMiddleware,
    errorHandler,
    notFoundHandler,
    requestLogger
} from './middleware/index.js';
import { logger } from './utils/logger.js';

import KafkaConsumerService from './kafkaConsumer.js';
import ClickHouseService from './clickhouse.js';
import AWSECSService from './awsECS.js';

import { clickhouseConfig } from './config/clickhouse.js';
import { kafkaConfig } from './config/kafka.js';
import { awsConfig } from './config/aws.js';
import { prismaClient } from './lib/prisma.js';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 9000;

// Trust proxy for rate limiting behind reverse proxies
app.set('trust proxy', 1);

// Request logging
app.use(requestLogger);

// CORS
app.use(corsMiddleware);

// Body parsing
app.use(express.json({ limit: '10kb' }));
app.use(express.urlencoded({ extended: true, limit: '10kb' }));

// Security headers
app.use((_req, res, next) => {
    res.setHeader('X-Content-Type-Options', 'nosniff');
    res.setHeader('X-Frame-Options', 'DENY');
    res.setHeader('X-XSS-Protection', '1; mode=block');
    res.removeHeader('X-Powered-By');
    next();
});

// Initialize external services
const kafkaClient = new Kafka(kafkaConfig);
const kafkaConsumer = new KafkaConsumerService(kafkaClient, 'mini-vercel-build-logs');

const ecsClient = new ECSClient({
    credentials: awsConfig,
    region: awsConfig.region
});
const awsECSService = new AWSECSService(ecsClient);

logger.debug('ClickHouse config', { url: clickhouseConfig.url });
const clickHouseClient = createClient(clickhouseConfig);
const clickHouseService = new ClickHouseService(clickHouseClient);

// Mount routes
const router = createRouter({
    awsECSService,
    clickHouseClient
});
app.use(router);

// 404 handler
app.use(notFoundHandler);

// Global error handler
app.use(errorHandler);

// Kafka consumer for build logs
kafkaConsumer.listenForMessagesInBatch('mini-vercel-build-logs', async (message) => {
    const { key, value } = message;
    if (!key || !value) return;

    try {
        const { project_id, deployment_id, log } = JSON.parse(value.toString());

        if (log && typeof log === 'string') {
            if (log.toLowerCase() === "INFO: Starting build pipeline...".toLowerCase()) {
                await prismaClient.deployment.update({
                    where: { id: deployment_id },
                    data: { status: 'IN_PROGRESS' }
                });
                logger.info('Deployment marked as IN_PROGRESS', { deployment_id });
            }
            // Check for successful build completion
            if (log.toLowerCase() === "info: pipeline completed successfully.") {
                await prismaClient.deployment.update({
                    where: { id: deployment_id },
                    data: { status: 'READY' }
                });
                logger.info('Deployment marked as READY', { deployment_id });
            }

            // Check for build failure
            if (log.toLowerCase().startsWith('error:') && log.toLowerCase().includes('pipeline failed')) {
                await prismaClient.deployment.update({
                    where: { id: deployment_id },
                    data: { status: 'FAIL' }
                });
                logger.warn('Deployment marked as FAIL', { deployment_id, log });
            }
        }

        logger.debug('Received build log', { project_id, deployment_id });

        const { query_id } = await clickHouseService.insertLog('log_events', { deployment_id, log });
        logger.debug('Log inserted to ClickHouse', { query_id });
    } catch (error) {
        logger.error('Failed to process Kafka message', error);
    }
});

// Start server
const httpServer = app.listen(PORT, () => {
    logger.info(`API server started`, { port: PORT, env: process.env.NODE_ENV || 'development' });
});

// Graceful shutdown
const gracefulShutdown = async (signal: string): Promise<void> => {
    logger.info(`${signal} received, starting graceful shutdown`);

    httpServer.close(() => {
        logger.info('HTTP server closed');
    });

    try {
        await kafkaConsumer.disconnect();
        logger.info('Kafka consumer disconnected');
    } catch (error) {
        logger.error('Error disconnecting Kafka consumer', error);
    }

    try {
        await clickHouseClient.close();
        logger.info('ClickHouse client closed');
    } catch (error) {
        logger.error('Error closing ClickHouse client', error);
    }

    process.exit(0);
};

process.on('SIGTERM', () => gracefulShutdown('SIGTERM'));
process.on('SIGINT', () => gracefulShutdown('SIGINT'));

// Unhandled rejection handler
process.on('unhandledRejection', (reason, promise) => {
    logger.error('Unhandled Rejection', reason as Error, { promise: String(promise) });
});

// Uncaught exception handler
process.on('uncaughtException', (error) => {
    logger.error('Uncaught Exception', error);
    process.exit(1);
});

export default app;
