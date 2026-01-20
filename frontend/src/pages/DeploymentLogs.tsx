import { useState, useEffect, useRef } from 'react';
import { useParams, Link } from 'react-router-dom';
import { api, type Deployment, type LogEntry } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import {
  ArrowLeft,
  Loader2,
  CheckCircle,
  XCircle,
  RefreshCw,
  Clock,
  Terminal,
  ExternalLink,
} from 'lucide-react';

export function DeploymentLogs() {
  const { deploymentId } = useParams<{ deploymentId: string }>();
  const [deployment, setDeployment] = useState<Deployment | null>(null);
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const logsEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (deploymentId) {
      fetchLogs();
    }
  }, [deploymentId]);

  useEffect(() => {
    // Auto-refresh logs if deployment is in progress
    if (deployment && ['QUEUED', 'IN_PROGRESS'].includes(deployment.status)) {
      const interval = setInterval(fetchLogs, 7000);
      return () => clearInterval(interval);
    }
  }, [deployment?.status]);

  useEffect(() => {
    // Scroll to bottom when new logs arrive
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [logs]);

  const fetchLogs = async () => {
    try {
      const response = await api.getDeploymentLogs(deploymentId!);
      setDeployment(response.data.deployment);
      setLogs(response.data.logs);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'READY':
        return <CheckCircle className="h-5 w-5 text-green-500" />;
      case 'FAIL':
        return <XCircle className="h-5 w-5 text-destructive" />;
      case 'QUEUED':
      case 'IN_PROGRESS':
        return <RefreshCw className="h-5 w-5 text-yellow-500 animate-spin" />;
      default:
        return <Clock className="h-5 w-5 text-muted-foreground" />;
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

  if (error || !deployment) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="p-4 bg-destructive/10 text-destructive rounded-md">
          Failed to load deployment: {error || 'Deployment not found'}
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-6">
        <Link
          to={`/projects/${deployment.projectId}`}
          className="inline-flex items-center text-sm text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="h-4 w-4 mr-1" />
          Back to Project
        </Link>
      </div>

      <div className="flex items-start justify-between mb-8">
        <div className="flex items-center gap-4">
          {getStatusIcon(deployment.status)}
          <div>
            <h1 className="text-2xl font-bold flex items-center gap-3">
              Deployment Logs
              {getStatusBadge(deployment.status)}
            </h1>
            <p className="text-muted-foreground mt-1">
              Started {formatDate(deployment.createdAt)}
            </p>
          </div>
        </div>
        {deployment.status === 'READY' && deployment.project && (
          <a
            href={`http://${deployment.project.subDomain}.localhost:8001`}
            target="_blank"
            rel="noopener noreferrer"
          >
            <Button variant="outline" className="gap-2">
              <ExternalLink className="h-4 w-4" />
              Visit Site
            </Button>
          </a>
        )}
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between py-3 px-4 border-b">
          <CardTitle className="text-sm font-medium flex items-center gap-2">
            <Terminal className="h-4 w-4" />
            Build Output
          </CardTitle>
          {['QUEUED', 'IN_PROGRESS'].includes(deployment.status) && (
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Loader2 className="h-4 w-4 animate-spin" />
              Live
            </div>
          )}
        </CardHeader>
        <CardContent className="p-0">
          <ScrollArea className="h-[600px]">
            <div className="p-4 font-mono text-sm bg-zinc-950 text-zinc-100 min-h-full">
              {logs.length === 0 ? (
                <div className="text-zinc-500 text-center py-8">
                  {['QUEUED', 'IN_PROGRESS'].includes(deployment.status) ? (
                    <>
                      <Loader2 className="h-6 w-6 animate-spin mx-auto mb-2" />
                      Waiting for logs...
                    </>
                  ) : (
                    'No logs available'
                  )}
                </div>
              ) : (
                <div className="space-y-1">
                  {logs.map((log, index) => (
                    <div key={log.event_id || index} className="flex">
                      <span className="text-zinc-500 select-none w-12 flex-shrink-0">
                        {String(index + 1).padStart(3, ' ')}
                      </span>
                      <span className="whitespace-pre-wrap break-all">{log.log}</span>
                    </div>
                  ))}
                  <div ref={logsEndRef} />
                </div>
              )}
            </div>
          </ScrollArea>
        </CardContent>
      </Card>

      <div className="mt-4 flex items-center justify-between text-sm text-muted-foreground">
        <div>
          Deployment ID: <span className="font-mono">{deployment.id}</span>
        </div>
        <div>
          {logs.length} log {logs.length === 1 ? 'entry' : 'entries'}
        </div>
      </div>
    </div>
  );
}
