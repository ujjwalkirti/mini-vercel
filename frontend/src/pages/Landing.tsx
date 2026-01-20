import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';
import { Rocket, GitBranch, Globe, Zap, Shield, BarChart3 } from 'lucide-react';

export function Landing() {
  const { user } = useAuth();

  return (
    <div className="flex flex-col">
      {/* Hero Section */}
      <section className="container mx-auto px-4 py-24 text-center">
        <div className="flex justify-center mb-6">
          <div className="p-3 bg-primary/10 rounded-full">
            <Rocket className="h-12 w-12 text-primary" />
          </div>
        </div>
        <h1 className="text-4xl md:text-6xl font-bold tracking-tight mb-6">
          Deploy your projects
          <br />
          <span className="text-primary">in seconds</span>
        </h1>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto mb-8">
          Mini Vercel is a self-hosted deployment platform that lets you deploy your
          static sites and frontend applications with ease.
        </p>
        <div className="flex gap-4 justify-center">
          {user ? (
            <Link to="/projects">
              <Button size="lg" className="gap-2">
                <Rocket className="h-5 w-5" />
                Go to Dashboard
              </Button>
            </Link>
          ) : (
            <>
              <Link to="/signup">
                <Button size="lg" className="gap-2">
                  <Rocket className="h-5 w-5" />
                  Start Deploying
                </Button>
              </Link>
              <Link to="/login">
                <Button size="lg" variant="outline">
                  Sign In
                </Button>
              </Link>
            </>
          )}
        </div>
      </section>

      {/* Features Section */}
      <section className="bg-muted/50 py-24">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center mb-12">
            Everything you need to deploy
          </h2>
          <div className="grid md:grid-cols-3 gap-8">
            <FeatureCard
              icon={<GitBranch className="h-8 w-8" />}
              title="Git Integration"
              description="Connect your GitHub repository and deploy with a single click. Automatic builds on every push."
            />
            <FeatureCard
              icon={<Globe className="h-8 w-8" />}
              title="Instant Deployments"
              description="Your projects go live instantly with unique URLs. Share your work with the world in seconds."
            />
            <FeatureCard
              icon={<Zap className="h-8 w-8" />}
              title="Fast Builds"
              description="Powered by containerized build runners for fast and reliable builds every time."
            />
            <FeatureCard
              icon={<Shield className="h-8 w-8" />}
              title="Secure by Default"
              description="All deployments are secured with authentication. Your projects are safe with us."
            />
            <FeatureCard
              icon={<BarChart3 className="h-8 w-8" />}
              title="Real-time Logs"
              description="Watch your builds in real-time with detailed logs. Debug issues as they happen."
            />
            <FeatureCard
              icon={<Rocket className="h-8 w-8" />}
              title="Self-Hosted"
              description="Run on your own infrastructure. Full control over your deployment pipeline."
            />
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="container mx-auto px-4 py-24 text-center">
        <h2 className="text-3xl font-bold mb-4">Ready to get started?</h2>
        <p className="text-muted-foreground mb-8 max-w-xl mx-auto">
          Create your account and start deploying your projects today.
          No credit card required.
        </p>
        {!user && (
          <Link to="/signup">
            <Button size="lg">Create Free Account</Button>
          </Link>
        )}
      </section>

      {/* Footer */}
      <footer className="border-t py-8">
        <div className="container mx-auto px-4 text-center text-muted-foreground">
          <p>Mini Vercel - Self-hosted deployment platform</p>
        </div>
      </footer>
    </div>
  );
}

function FeatureCard({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
}) {
  return (
    <div className="p-6 rounded-lg border bg-card">
      <div className="text-primary mb-4">{icon}</div>
      <h3 className="text-xl font-semibold mb-2">{title}</h3>
      <p className="text-muted-foreground">{description}</p>
    </div>
  );
}
