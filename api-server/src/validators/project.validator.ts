import { body, param } from 'express-validator';

const GITHUB_URL_REGEX = /^https:\/\/github\.com\/[\w-]+\/[\w.-]+(?:\.git)?$/;

export const createProjectValidator = [
    body('name')
        .trim()
        .notEmpty()
        .withMessage('Project name is required')
        .isLength({ min: 1, max: 100 })
        .withMessage('Project name must be between 1 and 100 characters')
        .matches(/^[\w\s-]+$/)
        .withMessage('Project name can only contain letters, numbers, spaces, underscores, and hyphens'),

    body('github_url')
        .trim()
        .notEmpty()
        .withMessage('GitHub URL is required')
        .isURL({ protocols: ['https'], require_protocol: true })
        .withMessage('Invalid URL format')
        .matches(GITHUB_URL_REGEX)
        .withMessage('Must be a valid GitHub repository URL (e.g., https://github.com/user/repo)')
];

export const projectIdValidator = [
    param('id')
        .trim()
        .notEmpty()
        .withMessage('Project ID is required')
        .isUUID()
        .withMessage('Invalid project ID format')
];

export const updateProjectValidator = [
    param('id')
        .trim()
        .notEmpty()
        .withMessage('Project ID is required')
        .isUUID()
        .withMessage('Invalid project ID format'),

    body('name')
        .optional()
        .trim()
        .isLength({ min: 1, max: 100 })
        .withMessage('Project name must be between 1 and 100 characters')
        .matches(/^[\w\s-]+$/)
        .withMessage('Project name can only contain letters, numbers, spaces, underscores, and hyphens'),

    body('customDomain')
        .optional()
        .trim()
        .isLength({ max: 253 })
        .withMessage('Custom domain must be at most 253 characters')
        .matches(/^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,}$/i)
        .withMessage('Invalid domain format')
];
