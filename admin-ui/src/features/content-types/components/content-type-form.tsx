import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "@/hooks/use-toast";
import { contentTypesService } from "@/features/content-types/services/content-types-service";
import type { ContentType } from "@/types/api";

interface Props {
  contentType?: ContentType | null; // null → create mode
  onClose: () => void;
  onSuccess?: () => void; // callback untuk refetch
}

export function ContentTypeForm({ contentType, onClose, onSuccess }: Props) {
  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Prefill kalau edit mode
  useEffect(() => {
    if (contentType) {
      setName(contentType.name);
      setSlug(contentType.slug);
      // backend belum ada description → opsional
    }
  }, [contentType]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      if (contentType) {
        // update
        await contentTypesService.updateContentType(contentType.id, {
          name,
          slug,
        });
        toast({ title: "Updated", description: "Content type updated successfully." });
      } else {
        // create
        await contentTypesService.createContentType({
          name,
          slug,
        });
        toast({ title: "Created", description: "Content type created successfully." });
      }

      if (onSuccess) onSuccess();
      onClose();
    } catch (err: any) {
      toast({
        title: "Error",
        description: err.message || "Something went wrong",
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4 p-6 border rounded-md bg-background">
      <h2 className="text-xl font-semibold">
        {contentType ? "Edit Content Type" : "Create Content Type"}
      </h2>

      <div>
        <Label htmlFor="name">Name</Label>
        <Input
          id="name"
          value={name}
          onChange={(e) => {
            setName(e.target.value);
            if (!contentType) {
              setSlug(e.target.value.toLowerCase().replace(/\s+/g, "-"));
            }
          }}
          required
        />
      </div>

      <div>
        <Label htmlFor="slug">Slug</Label>
        <Input
          id="slug"
          value={slug}
          onChange={(e) => setSlug(e.target.value)}
          required
        />
      </div>

      <div className="flex justify-end space-x-2">
        <Button type="button" variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? "Saving..." : contentType ? "Update" : "Create"}
        </Button>
      </div>
    </form>
  );
}
