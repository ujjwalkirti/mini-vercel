import type { Request, Response, NextFunction } from 'express';
import { logger } from '../utils/logger.js';
import { ApiResponseUtil } from '../utils/response.js';

export class AppError extends Error {
    public statusCode: number;
    public isOperational: boolean;

    constructor(message: string, statusCode: number, isOperational = true) {
        super(message);
        this.statusCode = statusCode;
        this.isOperational = isOperational;

        Error.captureStackTrace(this, this.constructor);
    }
}

export class BadRequestError extends AppError {
    constructor(message = 'Bad request') {
        super(message, 400);
    }
}

export class UnauthorizedError extends AppError {
    constructor(message = 'Unauthorized') {
        super(message, 401);
    }
}

export class ForbiddenError extends AppError {
    constructor(message = 'Access denied') {
        super(message, 403);
    }
}

export class NotFoundError extends AppError {
    constructor(message = 'Resource not found') {
        super(message, 404);
    }
}

export class ConflictError extends AppError {
    constructor(message = 'Resource already exists') {
        super(message, 409);
    }
}

export class TooManyRequestsError extends AppError {
    constructor(message = 'Too many requests') {
        super(message, 429);
    }
}

export class InternalServerError extends AppError {
    constructor(message = 'Internal server error') {
        super(message, 500, false);
    }
}

export const errorHandler = (
    err: Error | AppError,
    _req: Request,
    res: Response,
    _next: NextFunction
): Response => {
    if (err instanceof AppError) {
        logger.error(`AppError: ${err.message}`, err, {
            statusCode: err.statusCode,
            isOperational: err.isOperational
        });

        return ApiResponseUtil.error(
            res,
            err.message,
            err.statusCode,
            process.env.NODE_ENV === 'development' ? err.stack : undefined
        );
    }

    logger.error('Unhandled error', err);

    return ApiResponseUtil.internalError(res, err);
};

export const notFoundHandler = (_req: Request, res: Response): Response => {
    return ApiResponseUtil.notFound(res, 'Endpoint not found');
};

export const asyncHandler = <T>(
    fn: (req: Request, res: Response, next: NextFunction) => Promise<T>
) => {
    return (req: Request, res: Response, next: NextFunction): void => {
        Promise.resolve(fn(req, res, next)).catch(next);
    };
};
