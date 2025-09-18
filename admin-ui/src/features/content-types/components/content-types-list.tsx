import { useEffect, useState } from "react";
import { Plus, Edit, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import { toast } from "@/hooks/use-toast";
import { contentTypesService } from "@/features/content-types/services/content-types-service";
import type { ContentType } from "@/types/api";
import { ContentTypeForm } from "./content-type-form";


export function ContentTypesList() {
  const [contentTypes, setContentTypes] = useState<ContentType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isError, setIsError] = useState(false);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [editingContentType, setEditingContentType] = useState<ContentType | null>(null);

  // fetch data manual
  const fetchContentTypes = async () => {
    try {
      setIsLoading(true);
      setIsError(false);
      const data = await contentTypesService.getContentTypes();
      setContentTypes(data);
    } catch (err) {
      console.error(err);
      setIsError(true);
      toast({
        title: "Error",
        description: "Failed to load content types",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchContentTypes();
  }, []);

  const handleDelete = async (id: string) => {
    try {
      await contentTypesService.deleteContentType(id);
      toast({
        title: "Deleted",
        description: "Content type has been deleted successfully.",
      });
      fetchContentTypes(); // refetch setelah delete
    } catch (err: any) {
      toast({
        title: "Error",
        description: err.message,
        variant: "destructive",
      });
    }
  };

  if (isLoading) return <LoadingSpinner size="lg" className="mt-8" />;

  if (isError) {
    return (
      <div className="p-6 text-center text-red-500">
        Failed to load content types.{" "}
        <Button variant="outline" onClick={fetchContentTypes}>
          Retry
        </Button>
      </div>
    );
  }

  if (showCreateForm || editingContentType) {
  return (
    <ContentTypeForm
      contentType={editingContentType}
      onClose={() => {
        setShowCreateForm(false);
        setEditingContentType(null);
      }}
      onSuccess={fetchContentTypes} // biar list auto refresh
    />
  );
}

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Content Types</h1>
        <Button onClick={() => setShowCreateForm(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Create Content Type
        </Button>
      </div>

      {contentTypes.length === 0 ? (
        <div className="text-center py-10 border rounded">
          <h3 className="text-lg font-semibold mb-2">No content types yet</h3>
          <p className="text-muted-foreground mb-4">
            Create your first content type to start managing content.
          </p>
          <Button onClick={() => setShowCreateForm(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Create Content Type
          </Button>
        </div>
      ) : (
        <div className="overflow-x-auto rounded border">
          <table className="min-w-full border-collapse">
            <thead className="bg-muted">
              <tr>
                <th className="px-4 py-2 text-left text-sm font-medium">Name</th>
                <th className="px-4 py-2 text-left text-sm font-medium">Slug</th>
                <th className="px-4 py-2 text-left text-sm font-medium">Fields</th>
                <th className="px-4 py-2 text-sm font-medium text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              {contentTypes.map((ct) => (
                <tr key={ct.id} className="border-t hover:bg-muted/30">
                  <td className="px-4 py-2">{ct.name}</td>
                  <td className="px-4 py-2">{ct.slug}</td>
                  <td className="px-4 py-2">{ct.fields?.length ?? 0}</td>
                  <td className="px-4 py-2 text-right space-x-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setEditingContentType(ct)}
                    >
                      <Edit className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleDelete(ct.id)}
                    >
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
