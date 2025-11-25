const express = require('express');
const httpProxy = require('http-proxy');


const app = express();
const proxy = httpProxy.createProxyServer();
const PORT = 8000;

app.use((req, res) => {
    const hostname = req.hostname;
    const subdomain = hostname.split('.')[0];
    const resolvesTo = `https://bpaprod.blob.core.windows.net/build-outputs/${subdomain}/index.html`;

    proxy.web(req, res, { target: resolvesTo });
});

app.listen(PORT, () => {
    console.log(`Reverse-proxy server listening on port ${PORT}`);
});
