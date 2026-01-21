import express from 'express';
import httpProxy from 'http-proxy';
import dotenv from 'dotenv';
import { prismaClient } from './lib/prisma.js';

dotenv.config();

const app = express();
const proxy = httpProxy.createProxyServer();
const PORT = 8001;

const BASE_PATH = process.env.R2_PUBLIC_URL || 'https://pub-xxx.r2.dev';

// Health check endpoint
app.get('/health', (_req, res) => {
    res.json({ status: 'ok' });
});

app.use(async (req, res) => {
    const hostname = req.hostname;
    const subdomain = hostname.split('.')[0];

    // Find the project by subdomain and get the latest READY deployment
    const project = await prismaClient.project.findFirst({
        where: { subDomain: subdomain },
        include: {
            Deployment: {
                where: { status: 'READY' },
                orderBy: { createdAt: 'desc' },
                take: 1
            }
        }
    });

    if (!project || project.Deployment.length === 0) {
        res.status(404).send('Deployment not found');
        return;
    }

    const blobUrl = `${BASE_PATH}/${project.id}`;

    proxy.web(req, res, { target: blobUrl, changeOrigin: true });
});

proxy.on('proxyReq', (proxyReq: any, req, res) => {
    const url = req.url;
    if (url === '/') {
        proxyReq.path += 'index.html';
    }
});


app.listen(PORT, () => {
    console.log(`Reverse-proxy server listening on port ${PORT}`);
});
