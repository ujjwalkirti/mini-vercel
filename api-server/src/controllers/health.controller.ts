import type { Request, Response } from 'express';
import { ApiResponseUtil } from '../utils/response.js';

export class HealthController {
    static check(_req: Request, res: Response): void {
        ApiResponseUtil.success(res, {
            status: 'ok',
            timestamp: new Date().toISOString(),
            uptime: process.uptime()
        });
    }
}
