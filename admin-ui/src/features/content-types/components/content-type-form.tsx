import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { toast } from "@/hooks/use-toast";
import { contentTypesService } from "@/features/content-types/services/content-types-service";
import type { ContentType, ContentField } from "@/types/api";

interface Props {
  contentType?: ContentType | null;
  onClose: () => void;
  onSuccess?: () => void;
}

export function ContentTypeForm({ contentType, onClose, onSuccess }: Props) {
  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const [fields, setFields] = useState<ContentField[]>([]);
  const [newFieldName, setNewFieldName] = useState("");
  const [newFieldKind, setNewFieldKind] = useState("text");

  useEffect(() => {
    if (contentType) {
      setName(contentType.name);
      setSlug(contentType.slug);
      setFields(contentType.fields);
    }
  }, [contentType]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      if (contentType) {
        await contentTypesService.updateContentType(contentType.id, { name, slug });
        toast({ title: "Updated", description: "Content type updated successfully." });
      } else {
        const created = await contentTypesService.createContentType({ name, slug });
        toast({ title: "Created", description: "Content type created successfully." });

        // simpan fields yang ditahan di state
        for (const f of fields) {
          await contentTypesService.addField(created.id, { name: f.name, kind: f.kind, options: {} });
        }
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

  const handleAddField = async () => {
    if (!newFieldName.trim()) return;

    if (contentType) {
      try {
        const createdField = await contentTypesService.addField(contentType.id, {
          name: newFieldName,
          kind: newFieldKind,
          options: {},
        });
        setFields([...fields, createdField]);
        toast({ title: "Field added", description: `"${newFieldName}" added successfully.` });
      } catch (err: any) {
        toast({
          title: "Error",
          description: err.message || "Failed to add field",
          variant: "destructive",
        });
      }
    } else {
      setFields([...fields, { id: crypto.randomUUID(), contentTypeId: "", name: newFieldName, kind: newFieldKind, options: {} }]);
    }

    setNewFieldName("");
    setNewFieldKind("text");
  };

  const handleDeleteField = async (field: ContentField, index: number) => {
    if (contentType) {
      try {
        await contentTypesService.deleteField(contentType.id, field.id);
        setFields(fields.filter((_, i) => i !== index));
        toast({ title: "Deleted", description: `"${field.name}" removed.` });
      } catch (err: any) {
        toast({
          title: "Error",
          description: err.message || "Failed to delete field",
          variant: "destructive",
        });
      }
    } else {
      setFields(fields.filter((_, i) => i !== index));
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6 p-6 border rounded-md bg-background">
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
        <Input id="slug" value={slug} onChange={(e) => setSlug(e.target.value)} required />
      </div>

      <div className="space-y-2">
        <Label>Fields</Label>
        {fields.length === 0 && <p className="text-sm text-muted-foreground">No fields yet.</p>}
        {fields.map((f, idx) => (
          <div key={f.id} className="flex items-center space-x-2">
            <Input value={f.name} disabled className="flex-1" />
            <Input value={f.kind} disabled className="w-32" />
            {/* <Button type="button" variant="outline" onClick={() => handleDeleteField(f, idx)}>
              Remove
            </Button> */}
          </div>
        ))}

        <div className="flex items-center space-x-2">
          <Input
            placeholder="Field name"
            value={newFieldName}
            onChange={(e) => setNewFieldName(e.target.value)}
            className="flex-1"
          />
          <Select value={newFieldKind} onValueChange={setNewFieldKind}>
            <SelectTrigger className="w-32">
              <SelectValue placeholder="Type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="text">Text</SelectItem>
              <SelectItem value="number">Number</SelectItem>
              <SelectItem value="boolean">Boolean</SelectItem>
              <SelectItem value="date">Date</SelectItem>
              <SelectItem value="media">Media</SelectItem>
              <SelectItem value="json">JSON</SelectItem>
              <SelectItem value="wysiwyg">WYSIWYG</SelectItem>
            </SelectContent>
          </Select>
          <Button type="button" variant="outline" onClick={handleAddField}>
            Add
          </Button>
        </div>
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
