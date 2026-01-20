import { supabase } from './supabase';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:9000';

async function getAuthHeaders() {
  const { data: { session } } = await supabase.auth.getSession();
  return {
    'Content-Type': 'application/json',
    ...(session?.access_token && { Authorization: `Bearer ${session.access_token}` }),
  };
}

async function fetchWithAuth(endpoint: string, options: RequestInit = {}) {
  const headers = await getAuthHeaders();
  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers: {
      ...headers,
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: 'Request failed' }));
    throw new Error(error.message || error.error || 'Request failed');
  }

  return response.json();
}

export const api = {
  // Projects
  getProjects: () => fetchWithAuth('/projects'),

  getProject: (id: string) => fetchWithAuth(`/projects/${id}`),

  createProject: (data: { name: string; github_url: string }) =>
    fetchWithAuth('/projects', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  deleteProject: (id: string) =>
    fetchWithAuth(`/projects/${id}`, {
      method: 'DELETE',
    }),

  // Deployments
  getDeployments: (projectId: string) =>
    fetchWithAuth(`/projects/${projectId}/deployments`),

  getDeployment: (id: string) => fetchWithAuth(`/deployments/${id}`),

  triggerDeployment: (projectId: string) =>
    fetchWithAuth('/deploy', {
      method: 'POST',
      body: JSON.stringify({ project_id: projectId }),
    }),

  // Logs
  getDeploymentLogs: (deploymentId: string) =>
    fetchWithAuth(`/deployments/${deploymentId}/logs`),
};

export type Project = {
  id: string;
  name: string;
  gitURL: string;
  subDomain: string;
  customDomain: string | null;
  userId: string;
  createdAt: string;
  updatedAt: string;
  Deployment?: Deployment[];
};

export type DeploymentStatus = 'NOT_STARTED' | 'QUEUED' | 'IN_PROGRESS' | 'READY' | 'FAIL';

export type Deployment = {
  id: string;
  projectId: string;
  status: DeploymentStatus;
  createdAt: string;
  updatedAt: string;
  project?: Project;
};

export type LogEntry = {
  event_id: string;
  deployment_id: string;
  log: string;
  timestamp?: string;
};
