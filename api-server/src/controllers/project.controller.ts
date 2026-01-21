import type { Response } from 'express';
import { prismaClient } from '../lib/prisma';
import { generateSlug } from 'random-word-slugs';
import type { AuthenticatedRequest } from '../types';
import { ApiResponseUtil } from '../utils/response';
import { logger } from '../utils/logger';
import { NotFoundError } from '../middleware/errorHandler';

export class ProjectController {
    static async getAll(req: AuthenticatedRequest, res: Response): Promise<void> {
        const projects = await prismaClient.project.findMany({
            where: { userId: req.user!.id },
            include: {
                Deployment: {
                    orderBy: { createdAt: 'desc' },
                    take: 1
                }
            },
            orderBy: { createdAt: 'desc' }
        });

        ApiResponseUtil.success(res, projects);
    }

    static async getById(req: AuthenticatedRequest, res: Response): Promise<void> {
        const id = req.params.id as string;

        const project = await prismaClient.project.findFirst({
            where: {
                id,
                userId: req.user!.id
            },
            include: {
                Deployment: {
                    orderBy: { createdAt: 'desc' }
                }
            }
        });

        if (!project) {
            throw new NotFoundError('Project not found');
        }

        ApiResponseUtil.success(res, project);
    }

    static async create(req: AuthenticatedRequest, res: Response): Promise<void> {
        const { name, github_url } = req.body;

        const project = await prismaClient.project.create({
            data: {
                name,
                gitURL: github_url,
                subDomain: generateSlug(3),
                userId: req.user!.id
            }
        });

        logger.info('Project created', { projectId: project.id, userId: req.user!.id });

        ApiResponseUtil.created(res, project, 'Project created successfully');
    }

    static async delete(req: AuthenticatedRequest, res: Response): Promise<void> {
        const id = req.params.id as string;

        const project = await prismaClient.project.findFirst({
            where: {
                id,
                userId: req.user!.id
            }
        });

        if (!project) {
            throw new NotFoundError('Project not found');
        }

        await prismaClient.deployment.deleteMany({
            where: { projectId: id }
        });

        await prismaClient.project.delete({
            where: { id }
        });

        logger.info('Project deleted', { projectId: id, userId: req.user!.id });

        ApiResponseUtil.success(res, null, 'Project deleted successfully');
    }
}
