import cors, { type CorsOptions } from 'cors';

const getAllowedOrigins = (): string[] => {
    const origins: string[] = [];

    if (process.env.FRONTEND_URL) {
        origins.push(process.env.FRONTEND_URL);
    }

    if (process.env.ALLOWED_ORIGINS) {
        const additionalOrigins = process.env.ALLOWED_ORIGINS.split(',').map(o => o.trim());
        origins.push(...additionalOrigins);
    }

    if (process.env.NODE_ENV === 'development') {
        origins.push('http://localhost:5173', 'http://localhost:3000', 'http://127.0.0.1:5173');
    }

    return origins;
};

const corsOptions: CorsOptions = {
    origin: (origin, callback) => {
        const allowedOrigins = getAllowedOrigins();

        if (!origin || allowedOrigins.includes(origin)) {
            callback(null, true);
        } else {
            callback(new Error('Not allowed by CORS'));
        }
    },
    credentials: true,
    methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS'],
    allowedHeaders: [
        'Content-Type',
        'Authorization',
        'X-Requested-With',
        'Accept',
        'Origin',
        'X-Request-ID'
    ],
    exposedHeaders: ['X-Request-ID', 'X-RateLimit-Limit', 'X-RateLimit-Remaining', 'X-RateLimit-Reset'],
    maxAge: 86400, // 24 hours
    preflightContinue: false,
    optionsSuccessStatus: 204
};

export const corsMiddleware = cors(corsOptions);

export const corsOptionsForSocket: CorsOptions = {
    origin: getAllowedOrigins(),
    credentials: true
};
