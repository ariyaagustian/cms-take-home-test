import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { FileText, Users, Edit } from "lucide-react";
import { Link } from "react-router-dom";
import { apiClient } from "@/lib/api";

interface Article {
  id: string;
  slug: string;
  title: string;
  description: string;
  publishedAt: string;
}

const Index = () => {
  const [articles, setArticles] = useState<Article[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadArticles = async () => {
      try {
        setLoading(true);
        // fetch langsung tanpa react-query
        const res = await apiClient.get<any>("/api/public/post?limit=10");
        console.log("API response:", res);

        // handle struktur JSON
        const rawData = Array.isArray(res.data) ? res.data : res.data?.data || [];

        const mapped = rawData.map((raw: any) => ({
          id: raw.ID,
          slug: raw.Slug,
          title: raw.Data?.title ?? "Untitled",
          description: raw.Data?.description ?? "",
          publishedAt: raw.PublishedAt,
        }));

        setArticles(mapped);
      } catch (err) {
        console.error("Failed to fetch articles", err);
      } finally {
        setLoading(false);
      }
    };

    loadArticles();
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted/30 to-primary/5">
      <div className="container mx-auto px-4 py-16">
        {/* Hero Section */}
        <div className="text-center mb-16">
          <h1 className="text-5xl font-bold mb-6 bg-gradient-primary bg-clip-text text-transparent">
            Content Management System
          </h1>
          <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
            Powerful, flexible content management built for modern teams. 
            Create, manage, and publish content with ease.
          </p>
          <div className="flex gap-4 justify-center">
            <Button asChild className="bg-gradient-primary hover:opacity-90">
              <Link to="/admin">
                <Edit className="h-4 w-4 mr-2" />
                Admin Panel
              </Link>
            </Button>
            <Button variant="outline" asChild>
              <Link to="/login">Sign In</Link>
            </Button>
          </div>
        </div>

        {/* Features */}
        <div className="grid md:grid-cols-3 gap-8 mb-16">
          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <FileText className="h-10 w-10 text-primary mb-2" />
              <CardTitle>Content Management</CardTitle>
              <CardDescription>
                Create and manage structured content with custom fields and flexible layouts
              </CardDescription>
            </CardHeader>
          </Card>

          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <Users className="h-10 w-10 text-primary mb-2" />
              <CardTitle>Role-Based Access</CardTitle>
              <CardDescription>
                Secure user management with granular permissions and role-based access control
              </CardDescription>
            </CardHeader>
          </Card>

          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <Edit className="h-10 w-10 text-primary mb-2" />
              <CardTitle>Easy Editing</CardTitle>
              <CardDescription>
                Intuitive admin interface with real-time preview and collaborative editing
              </CardDescription>
            </CardHeader>
          </Card>
        </div>

        {/* Latest Articles */}
        <Card>
          <CardHeader>
            <CardTitle>Latest Articles</CardTitle>
            <CardDescription>Recent content from our CMS</CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <p className="text-muted-foreground text-sm">Loading...</p>
            ) : articles.length === 0 ? (
              <p className="text-muted-foreground text-sm">
                No articles published yet.
              </p>
            ) : (
              <div className="space-y-4">
                {articles.map((a) => (
                  <div
                    key={a.id} // ✅ gunakan id biar unik
                    className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors"
                  >
                    <div>
                      <h3 className="font-semibold">{a.title}</h3>
                      <p className="text-sm text-muted-foreground">{a.description}</p>
                    </div>
                    <span className="text-sm text-muted-foreground">
                      {new Date(a.publishedAt).toLocaleDateString()}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default Index;
