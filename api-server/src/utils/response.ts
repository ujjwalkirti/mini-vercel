import type { Response } from 'express';
import type { ApiResponse } from '../types/index.js';

export class ApiResponseUtil {
    static success<T>(res: Response, data: T, message?: string, statusCode = 200): Response {
        const response: ApiResponse<T> = {
            success: true,
            data,
            ...(message && { message })
        };
        return res.status(statusCode).json(response);
    }

    static created<T>(res: Response, data: T, message = 'Created successfully'): Response {
        return this.success(res, data, message, 201);
    }

    static error(
        res: Response,
        message: string,
        statusCode = 500,
        error?: string
    ): Response {
        const response: ApiResponse = {
            success: false,
            message,
            ...(error && { error })
        };
        return res.status(statusCode).json(response);
    }

    static badRequest(res: Response, message: string, error?: string): Response {
        return this.error(res, message, 400, error);
    }

    static unauthorized(res: Response, message = 'Unauthorized', error?: string): Response {
        return this.error(res, message, 401, error);
    }

    static forbidden(res: Response, message = 'Access denied', error?: string): Response {
        return this.error(res, message, 403, error);
    }

    static notFound(res: Response, message = 'Resource not found', error?: string): Response {
        return this.error(res, message, 404, error);
    }

    static validationError(res: Response, errors: Array<{ field: string; message: string }>): Response {
        const response: ApiResponse = {
            success: false,
            message: 'Validation failed',
            errors
        };
        return res.status(400).json(response);
    }

    static internalError(res: Response, error?: Error): Response {
        const message = 'Internal Server Error';
        const errorMessage = process.env.NODE_ENV === 'development' ? error?.message : undefined;
        return this.error(res, message, 500, errorMessage);
    }

    static tooManyRequests(res: Response, message = 'Too many requests, please try again later'): Response {
        return this.error(res, message, 429);
    }
}
