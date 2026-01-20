import { Router } from 'express';
import healthRoutes from './health.routes';
import projectRoutes from './project.routes';
import { createDeploymentRoutes } from './deployment.routes';
import { DeploymentController } from '../controllers';
import { apiRateLimiter } from '../middleware';
import type AWSECSService from '../awsECS';
import type { ClickHouseClient } from '@clickhouse/client';

interface RouteDependencies {
    awsECSService: AWSECSService;
    clickHouseClient: ClickHouseClient;
}

export const createRouter = (deps: RouteDependencies): Router => {
    const router = Router();

    const deploymentController = new DeploymentController({
        awsECSService: deps.awsECSService,
        clickHouseClient: deps.clickHouseClient
    });

    router.use(apiRateLimiter);

    router.use('/health', healthRoutes);

    router.use('/projects', projectRoutes);

    const deploymentRoutes = createDeploymentRoutes(deploymentController);
    router.use('/', deploymentRoutes);

    return router;
};

export default createRouter;
