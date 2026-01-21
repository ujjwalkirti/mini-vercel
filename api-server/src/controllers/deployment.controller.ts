import type { Response } from 'express';
import { prismaClient } from '../lib/prisma';
import type { AuthenticatedRequest } from '../types';
import { ApiResponseUtil } from '../utils/response';
import { logger } from '../utils/logger';
import { NotFoundError, ForbiddenError } from '../middleware/errorHandler';
import { ecsConfig } from '../config/ecs';
import type AWSECSService from '../awsECS';
import type { ClickHouseClient } from '@clickhouse/client';

interface DeploymentControllerDependencies {
    awsECSService: AWSECSService;
    clickHouseClient: ClickHouseClient;
}

export class DeploymentController {
    private awsECSService: AWSECSService;
    private clickHouseClient: ClickHouseClient;

    constructor(deps: DeploymentControllerDependencies) {
        this.awsECSService = deps.awsECSService;
        this.clickHouseClient = deps.clickHouseClient;
    }

    async getByProject(req: AuthenticatedRequest, res: Response): Promise<void> {
        const projectId = req.params.projectId as string;

        const project = await prismaClient.project.findFirst({
            where: {
                id: projectId,
                userId: req.user!.id
            }
        });

        if (!project) {
            throw new NotFoundError('Project not found');
        }

        const deployments = await prismaClient.deployment.findMany({
            where: { projectId },
            orderBy: { createdAt: 'desc' }
        });

        ApiResponseUtil.success(res, deployments);
    }

    async getById(req: AuthenticatedRequest, res: Response): Promise<void> {
        const id = req.params.id as string;

        const deployment = await prismaClient.deployment.findUnique({
            where: { id },
            include: {
                project: true
            }
        });

        if (!deployment) {
            throw new NotFoundError('Deployment not found');
        }

        if (deployment.project.userId !== req.user!.id) {
            throw new ForbiddenError('Access denied');
        }

        ApiResponseUtil.success(res, deployment);
    }

    async create(req: AuthenticatedRequest, res: Response): Promise<void> {
        const { project_id } = req.body;

        const project = await prismaClient.project.findFirst({
            where: {
                id: project_id,
                userId: req.user!.id
            }
        });

        if (!project) {
            throw new NotFoundError('Project not found');
        }

        const deployment = await prismaClient.deployment.create({
            data: {
                status: 'QUEUED',
                project: { connect: { id: project_id } }
            }
        });

        const envVars = [
            { name: 'PROJECT_ID', value: project_id },
            { name: 'GIT_REPOSITORY_URL', value: project.gitURL },
            { name: 'KAFKA_BROKERS', value: process.env.KAFKA_BROKERS },
            { name: 'KAFKA_CLIENT_ID', value: process.env.KAFKA_CLIENT_ID },
            { name: 'KAFKA_USERNAME', value: process.env.KAFKA_USERNAME },
            { name: 'KAFKA_PASSWORD', value: process.env.KAFKA_PASSWORD },
            { name: 'R2_ACCOUNT_ID', value: process.env.R2_ACCOUNT_ID },
            { name: 'R2_ACCESS_KEY_ID', value: process.env.R2_ACCESS_KEY_ID },
            { name: 'R2_SECRET_ACCESS_KEY', value: process.env.R2_SECRET_ACCESS_KEY },
            { name: 'R2_BUCKET_NAME', value: process.env.R2_BUCKET_NAME },
            { name: 'DEPLOYMENT_ID', value: deployment.id }
        ];

        const ecsTaskProps = {
            cluster: ecsConfig.cluster,
            taskDefinition: ecsConfig.taskDefinition,
            image: ecsConfig.imageName,
            envVars,
            subnets: ecsConfig.subnets,
            securityGroups: ecsConfig.securityGroups,
            assignPublicIp: ecsConfig.assignPublicIp,
            launchType: ecsConfig.launchType,
            count: ecsConfig.count
        };

        await this.awsECSService.runTask(ecsTaskProps);

        logger.info('Deployment queued', {
            deploymentId: deployment.id,
            projectId: project_id,
            userId: req.user!.id
        });

        ApiResponseUtil.success(res, {
            deploymentId: deployment.id,
            status: 'Queued',
            url: `${project.subDomain}.localhost:8001`
        }, 'Build queued successfully');
    }

    async getLogs(req: AuthenticatedRequest, res: Response): Promise<void> {
        const id = req.params.id as string;

        const deployment = await prismaClient.deployment.findUnique({
            where: { id },
            include: {
                project: true
            }
        });

        if (!deployment) {
            throw new NotFoundError('Deployment not found');
        }

        if (deployment.project.userId !== req.user!.id) {
            throw new ForbiddenError('Access denied');
        }

        const result = await this.clickHouseClient.query({
            query: `SELECT event_id, deployment_id, log, timestamp FROM log_events WHERE deployment_id = {deployment_id:String} ORDER BY timestamp ASC`,
            query_params: {
                deployment_id: id
            },
            format: 'JSONEachRow'
        });

        const logs = await result.json();

        ApiResponseUtil.success(res, {
            deployment,
            logs
        });
    }
}
