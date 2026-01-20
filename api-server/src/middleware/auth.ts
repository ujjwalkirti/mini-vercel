import type { Response, NextFunction } from 'express';
import { supabaseAdmin } from '../config/supabase';
import type { AuthenticatedRequest } from '../types';
import { ApiResponseUtil } from '../utils/response';
import { logger } from '../utils/logger';

export const authMiddleware = async (
    req: AuthenticatedRequest,
    res: Response,
    next: NextFunction
): Promise<void> => {
    try {
        const authHeader = req.headers.authorization;

        if (!authHeader || !authHeader.startsWith('Bearer ')) {
            ApiResponseUtil.unauthorized(res, 'Unauthorized', 'Missing or invalid authorization header');
            return;
        }

        const token = authHeader.split(' ')[1];

        if (!token) {
            ApiResponseUtil.unauthorized(res, 'Unauthorized', 'Token not provided');
            return;
        }

        const { data: { user }, error } = await supabaseAdmin.auth.getUser(token);

        if (error || !user) {
            logger.warn('Authentication failed', { error: error?.message });
            ApiResponseUtil.unauthorized(res, 'Unauthorized', error?.message || 'Invalid token');
            return;
        }

        req.user = {
            id: user.id,
            email: user.email
        };

        next();
    } catch (error) {
        logger.error('Auth middleware error', error);
        ApiResponseUtil.internalError(res, error as Error);
    }
};

export const optionalAuthMiddleware = async (
    req: AuthenticatedRequest,
    _res: Response,
    next: NextFunction
): Promise<void> => {
    try {
        const authHeader = req.headers.authorization;

        if (authHeader?.startsWith('Bearer ')) {
            const token = authHeader.split(' ')[1];

            if (token) {
                const { data: { user } } = await supabaseAdmin.auth.getUser(token);

                if (user) {
                    req.user = {
                        id: user.id,
                        email: user.email
                    };
                }
            }
        }

        next();
    } catch {
        next();
    }
};
