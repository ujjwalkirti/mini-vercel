import { Router } from 'express';
import { HealthController } from '../controllers/index.js';

const router = Router();

router.get('/', HealthController.check);

export default router;
