import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { FileText, Users, Edit } from "lucide-react";
import { Link } from "react-router-dom";

const Index = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted/30 to-primary/5">
      {/* Hero Section */}
      <div className="container mx-auto px-4 py-16">
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

        {/* Recent Content */}
        <Card>
          <CardHeader>
            <CardTitle>Latest Articles</CardTitle>
            <CardDescription>
              Recent content from our CMS
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors">
                <div>
                  <h3 className="font-semibold">Getting Started with CMS</h3>
                  <p className="text-sm text-muted-foreground">Learn how to use this content management system</p>
                </div>
                <span className="text-sm text-muted-foreground">2 days ago</span>
              </div>
              <div className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors">
                <div>
                  <h3 className="font-semibold">Building Custom Content Types</h3>
                  <p className="text-sm text-muted-foreground">Create flexible content structures for your needs</p>
                </div>
                <span className="text-sm text-muted-foreground">1 week ago</span>
              </div>
              <div className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors">
                <div>
                  <h3 className="font-semibold">Managing Users and Permissions</h3>
                  <p className="text-sm text-muted-foreground">Set up roles and manage access control</p>
                </div>
                <span className="text-sm text-muted-foreground">2 weeks ago</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default Index;
