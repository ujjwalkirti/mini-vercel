import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { api, type Project, type Deployment } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import {
  ArrowLeft,
  Loader2,
  GitBranch,
  ExternalLink,
  Rocket,
  Clock,
  CheckCircle,
  XCircle,
  RefreshCw,
} from 'lucide-react';

export function ProjectDetail() {
  const { projectId } = useParams<{ projectId: string }>();
  const [project, setProject] = useState<Project | null>(null);
  const [deployments, setDeployments] = useState<Deployment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [deploying, setDeploying] = useState(false);

  useEffect(() => {
    if (projectId) {
      fetchProject();
    }
  }, [projectId]);

  const fetchProject = async () => {
    try {
      const response = await api.getProject(projectId!);
      setProject(response.data);
      setDeployments(response.data.Deployment || []);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDeploy = async () => {
    setDeploying(true);
    try {
      await api.triggerDeployment(projectId!);
      fetchProject();
    } catch (err: any) {
      alert('Failed to trigger deployment: ' + err.message);
    } finally {
      setDeploying(false);
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'READY':
        return <CheckCircle className="h-4 w-4 text-green-500" />;
      case 'FAIL':
        return <XCircle className="h-4 w-4 text-destructive" />;
      case 'QUEUED':
      case 'IN_PROGRESS':
        return <RefreshCw className="h-4 w-4 text-yellow-500 animate-spin" />;
      default:
        return <Clock className="h-4 w-4 text-muted-foreground" />;
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'READY':
        return <Badge variant="success">Ready</Badge>;
      case 'QUEUED':
        return <Badge variant="warning">Queued</Badge>;
      case 'IN_PROGRESS':
        return <Badge variant="warning">Building</Badge>;
      case 'FAIL':
        return <Badge variant="destructive">Failed</Badge>;
      default:
        return <Badge variant="secondary">{status}</Badge>;
    }
  };

  const formatDate = (date: string) => {
    return new Date(date).toLocaleString();
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8 flex items-center justify-center min-h-[60vh]">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !project) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="p-4 bg-destructive/10 text-destructive rounded-md">
          Failed to load project: {error || 'Project not found'}
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-6">
        <Link to="/projects" className="inline-flex items-center text-sm text-muted-foreground hover:text-foreground">
          <ArrowLeft className="h-4 w-4 mr-1" />
          Back to Projects
        </Link>
      </div>

      <div className="flex items-start justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold">{project.name}</h1>
          <div className="flex items-center gap-4 mt-2 text-muted-foreground">
            <span className="flex items-center gap-1">
              <GitBranch className="h-4 w-4" />
              {project.gitURL.replace('https://github.com/', '')}
            </span>
            <a
              href={project.gitURL}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-1 hover:text-foreground"
            >
              <ExternalLink className="h-4 w-4" />
              View Repository
            </a>
          </div>
        </div>
        <div className="flex items-center gap-2">
          {deployments[0]?.status === 'READY' && (
            <a
              href={`http://${project.subDomain}.localhost:8001`}
              target="_blank"
              rel="noopener noreferrer"
            >
              <Button variant="outline" className="gap-2">
                <ExternalLink className="h-4 w-4" />
                Visit Site
              </Button>
            </a>
          )}
          <Button onClick={handleDeploy} disabled={deploying} className="gap-2">
            {deploying ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Rocket className="h-4 w-4" />
            )}
            Deploy
          </Button>
        </div>
      </div>

      <div className="grid lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle>Deployments</CardTitle>
              <CardDescription>History of all deployments for this project</CardDescription>
            </CardHeader>
            <CardContent>
              {deployments.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <Rocket className="h-8 w-8 mx-auto mb-2 opacity-50" />
                  <p>No deployments yet</p>
                  <p className="text-sm">Click "Deploy" to create your first deployment</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {deployments.map((deployment, index) => (
                    <div key={deployment.id}>
                      <Link
                        to={`/deployments/${deployment.id}`}
                        className="flex items-center justify-between p-4 rounded-lg border hover:bg-muted/50 transition-colors"
                      >
                        <div className="flex items-center gap-4">
                          {getStatusIcon(deployment.status)}
                          <div>
                            <div className="font-medium flex items-center gap-2">
                              Deployment #{deployments.length - index}
                              {getStatusBadge(deployment.status)}
                            </div>
                            <div className="text-sm text-muted-foreground">
                              {formatDate(deployment.createdAt)}
                            </div>
                          </div>
                        </div>
                        <Button variant="ghost" size="sm">
                          View Logs
                        </Button>
                      </Link>
                      {index < deployments.length - 1 && <Separator className="my-2" />}
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        <div>
          <Card>
            <CardHeader>
              <CardTitle>Project Info</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <div className="text-sm text-muted-foreground">Subdomain</div>
                <div className="font-mono text-sm">{project.subDomain}.localhost:8001</div>
              </div>
              <Separator />
              <div>
                <div className="text-sm text-muted-foreground">Created</div>
                <div className="text-sm">{formatDate(project.createdAt)}</div>
              </div>
              <Separator />
              <div>
                <div className="text-sm text-muted-foreground">Last Updated</div>
                <div className="text-sm">{formatDate(project.updatedAt)}</div>
              </div>
              <Separator />
              <div>
                <div className="text-sm text-muted-foreground">Total Deployments</div>
                <div className="text-sm">{deployments.length}</div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
