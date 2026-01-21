import type { Request, Response, NextFunction } from 'express';
import { validationResult, type ValidationChain } from 'express-validator';
import { ApiResponseUtil } from '../utils/response.js';

export const validate = (validations: ValidationChain[]) => {
    return async (req: Request, res: Response, next: NextFunction): Promise<void> => {
        await Promise.all(validations.map(validation => validation.run(req)));

        const errors = validationResult(req);

        if (errors.isEmpty()) {
            next();
            return;
        }

        const formattedErrors = errors.array().map(error => ({
            field: 'path' in error ? error.path : 'unknown',
            message: error.msg
        }));

        ApiResponseUtil.validationError(res, formattedErrors);
    };
};
