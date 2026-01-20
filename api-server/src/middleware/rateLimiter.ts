import rateLimit, { type RateLimitRequestHandler, type Options } from 'express-rate-limit';
import type { Request, Response } from 'express';
import { ApiResponseUtil } from '../utils/response';
import { logger } from '../utils/logger';

const createRateLimiter = (options: Partial<Options>): RateLimitRequestHandler => {
    return rateLimit({
        standardHeaders: true,
        legacyHeaders: false,
        handler: (_req: Request, res: Response) => {
            ApiResponseUtil.tooManyRequests(res, 'Too many requests, please try again later');
        },
        skip: (_req: Request) => {
            return process.env.NODE_ENV === 'test';
        },
        keyGenerator: (req: Request) => {
            return req.ip || req.headers['x-forwarded-for']?.toString() || 'unknown';
        },
        ...options
    });
};

// General API rate limiter - 100 requests per minute
export const apiRateLimiter = createRateLimiter({
    windowMs: 60 * 1000, // 1 minute
    max: 100,
    message: 'Too many requests from this IP, please try again after a minute'
});

// Auth endpoints rate limiter - 10 requests per minute (stricter)
export const authRateLimiter = createRateLimiter({
    windowMs: 60 * 1000, // 1 minute
    max: 10,
    message: 'Too many authentication attempts, please try again after a minute'
});

// Deploy endpoint rate limiter - 10 deployments per minute
export const deployRateLimiter = createRateLimiter({
    windowMs: 60 * 1000, // 1 minute
    max: 10,
    message: 'Too many deployment requests, please try again after a minute',
    handler: (req: Request, res: Response) => {
        logger.warn('Deploy rate limit exceeded', { ip: req.ip });
        ApiResponseUtil.tooManyRequests(res, 'Too many deployment requests, please try again later');
    }
});

// Strict rate limiter for sensitive operations - 5 requests per minute
export const strictRateLimiter = createRateLimiter({
    windowMs: 60 * 1000, // 1 minute
    max: 5,
    message: 'Rate limit exceeded for this operation'
});

// Logs endpoint rate limiter - 60 requests per minute
export const logsRateLimiter = createRateLimiter({
    windowMs: 60 * 1000, // 1 minute
    max: 60,
    message: 'Too many log requests, please try again after a minute'
});
