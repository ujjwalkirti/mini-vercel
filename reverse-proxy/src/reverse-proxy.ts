import express from 'express';
import httpProxy from 'http-proxy';
import dotenv from 'dotenv';

dotenv.config();

const app = express();
const proxy = httpProxy.createProxyServer();
const PORT = 8001;

const BASE_PATH = process.env.R2_PUBLIC_URL || 'https://pub-xxx.r2.dev';

app.use((req, res) => {
    const hostname = req.hostname;
    const subdomain = hostname.split('.')[0];
    const blobUrl = `${BASE_PATH}/${subdomain}`;

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
