const express = require('express');
const httpProxy = require('http-proxy');


const app = express();
const proxy = httpProxy.createProxyServer();
const PORT = 8000;

const BASE_PATH = 'https://bpaprod.blob.core.windows.net/build-outputs';

app.use((req, res) => {
    const hostname = req.hostname;
    const subdomain = hostname.split('.')[0];
    const blobUrl = `${BASE_PATH}/${subdomain}`;

    proxy.web(req, res, { target: blobUrl, changeOrigin: true });

});

proxy.on('proxyReq', (proxyReq, req, res) => {
    const url = req.url;
    if (url === '/') {
        proxyReq.path += 'index.html';
    }
});


app.listen(PORT, () => {
    console.log(`Reverse-proxy server listening on port ${PORT}`);
});
