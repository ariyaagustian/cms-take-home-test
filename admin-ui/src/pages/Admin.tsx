import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { FileText, Layers, Upload, Users } from "lucide-react";
import { useDashboardStats } from "@/features/dashboard/hooks/use-dashboard";

const Admin = () => {
  const { contentTypes, media, users, isLoading } = useDashboardStats();

  const stats = [
    { title: "Content Types", value: contentTypes, icon: Layers, color: "text-blue-600" },
    { title: "Media Files", value: media, icon: Upload, color: "text-purple-600" },
    { title: "Active Users", value: users, icon: Users, color: "text-orange-600" },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-foreground">Dashboard</h1>
        <p className="text-muted-foreground">
          Welcome to your content management system
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {stats.map((stat) => (
          <Card key={stat.title} className="hover:shadow-md transition-shadow">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                {stat.title}
              </CardTitle>
              <stat.icon className={`h-4 w-4 ${stat.color}`} />
            </CardHeader>
            <CardContent>
              {isLoading ? (
                <div className="text-muted-foreground text-sm">Loading...</div>
              ) : (
                <div className="text-2xl font-bold">{stat.value}</div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default Admin;
