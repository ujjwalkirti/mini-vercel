import type { Request, Response, NextFunction } from 'express';
import { logger } from '../utils/logger.js';
import { randomUUID } from 'crypto';

declare global {
    namespace Express {
        interface Request {
            requestId?: string;
        }
    }
}

export const requestLogger = (req: Request, res: Response, next: NextFunction): void => {
    const requestId = req.headers['x-request-id']?.toString() || randomUUID();
    req.requestId = requestId;
    res.setHeader('X-Request-ID', requestId);

    const startTime = Date.now();

    res.on('finish', () => {
        const duration = Date.now() - startTime;
        const logContext = {
            requestId,
            method: req.method,
            url: req.originalUrl,
            statusCode: res.statusCode,
            duration: `${duration}ms`,
            ip: req.ip,
            userAgent: req.headers['user-agent']
        };

        if (res.statusCode >= 400) {
            logger.warn('Request completed with error', logContext);
        } else {
            logger.info('Request completed', logContext);
        }
    });

    next();
};
