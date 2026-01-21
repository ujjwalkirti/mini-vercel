import { Router } from 'express';
import type { DeploymentController } from '../controllers/index.js';
import { authMiddleware, asyncHandler, validate, deployRateLimiter, logsRateLimiter } from '../middleware/index.js';
import { createDeploymentValidator, deploymentIdValidator, projectDeploymentsValidator } from '../validators/index.js';

export const createDeploymentRoutes = (deploymentController: DeploymentController): Router => {
    const router = Router();

    router.use(authMiddleware);

    router.get(
        '/projects/:projectId/deployments',
        validate(projectDeploymentsValidator),
        asyncHandler(deploymentController.getByProject.bind(deploymentController))
    );

    router.get(
        '/deployments/:id',
        validate(deploymentIdValidator),
        asyncHandler(deploymentController.getById.bind(deploymentController))
    );

    router.post(
        '/deploy',
        deployRateLimiter,
        validate(createDeploymentValidator),
        asyncHandler(deploymentController.create.bind(deploymentController))
    );

    router.get(
        '/deployments/:id/logs',
        logsRateLimiter,
        validate(deploymentIdValidator),
        asyncHandler(deploymentController.getLogs.bind(deploymentController))
    );

    return router;
};
