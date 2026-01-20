import { body, param } from 'express-validator';

export const createDeploymentValidator = [
    body('project_id')
        .trim()
        .notEmpty()
        .withMessage('Project ID is required')
        .isUUID()
        .withMessage('Invalid project ID format')
];

export const deploymentIdValidator = [
    param('id')
        .trim()
        .notEmpty()
        .withMessage('Deployment ID is required')
        .isUUID()
        .withMessage('Invalid deployment ID format')
];

export const projectDeploymentsValidator = [
    param('projectId')
        .trim()
        .notEmpty()
        .withMessage('Project ID is required')
        .isUUID()
        .withMessage('Invalid project ID format')
];
