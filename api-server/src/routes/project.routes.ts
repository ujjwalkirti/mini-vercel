import { Router } from 'express';
import { ProjectController } from '../controllers/index.js';
import { authMiddleware, asyncHandler, validate } from '../middleware/index.js';
import { createProjectValidator, projectIdValidator } from '../validators/index.js';

const router = Router();

router.use(authMiddleware);

router.get(
    '/',
    asyncHandler(ProjectController.getAll)
);

router.get(
    '/:id',
    validate(projectIdValidator),
    asyncHandler(ProjectController.getById)
);

router.post(
    '/',
    validate(createProjectValidator),
    asyncHandler(ProjectController.create)
);

router.delete(
    '/:id',
    validate(projectIdValidator),
    asyncHandler(ProjectController.delete)
);

export default router;
