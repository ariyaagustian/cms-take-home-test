import { apiClient } from "@/lib/api";
import type { ContentType, ContentField } from "@/types/api";

export const contentTypesService = {
  async getContentTypes(): Promise<ContentType[]> {
    const res = await apiClient.get<{ data: any[] }>("/api/content-types");
    return res.data.map((raw) => ({
      id: raw.ID,
      name: raw.Name,
      slug: raw.Slug,
      createdAt: raw.CreatedAt,
      updatedAt: raw.UpdatedAt,
      fields: Array.isArray(raw.Fields)
        ? raw.Fields.map((f: any): ContentField => ({
            id: f.ID,
            contentTypeId: f.ContentTypeID,
            name: f.Name,
            kind: f.Kind,
            options: f.Options ?? {},
          }))
        : [],
    }));
  },

  async getContentType(id: string): Promise<ContentType> {
    const res = await apiClient.get<{ data: any }>(`/api/content-types/${id}`);
    const raw = res.data;
    return {
      id: raw.ID,
      name: raw.Name,
      slug: raw.Slug,
      createdAt: raw.CreatedAt,
      updatedAt: raw.UpdatedAt,
      fields: Array.isArray(raw.Fields)
        ? raw.Fields.map((f: any): ContentField => ({
            id: f.ID,
            contentTypeId: f.ContentTypeID,
            name: f.Name,
            kind: f.Kind,
            options: f.Options ?? {},
          }))
        : [],
    };
  },

  async createContentType(data: { name: string; slug: string; description?: string }) {
    const res = await apiClient.post<{ data: any }>("/api/content-types", data);
    const raw = res.data;
    return {
      id: raw.ID,
      name: raw.Name,
      slug: raw.Slug,
      createdAt: raw.CreatedAt,
      updatedAt: raw.UpdatedAt,
      fields: [],
    } as ContentType;
  },

  async updateContentType(
    id: string,
    data: Partial<{ name: string; slug: string; description?: string }>
  ) {
    const res = await apiClient.put<{ data: any }>(`/api/content-types/${id}`, data);
    const raw = res.data;
    return {
      id: raw.ID,
      name: raw.Name,
      slug: raw.Slug,
      createdAt: raw.CreatedAt,
      updatedAt: raw.UpdatedAt,
      fields: [],
    } as ContentType;
  },

  async deleteContentType(id: string): Promise<void> {
    await apiClient.delete(`/api/content-types/${id}`);
  },

  async addField(
    contentTypeId: string,
    field: { name: string; kind: string; options?: Record<string, unknown> }
  ): Promise<ContentField> {
    const res = await apiClient.post<{ data: any }>(
      `/api/content-types/${contentTypeId}/fields`,
      field
    );
    return {
      id: res.data.ID,
      contentTypeId: res.data.ContentTypeID,
      name: res.data.Name,
      kind: res.data.Kind,
      options: res.data.Options ?? {},
    };
  },

  async deleteField(contentTypeId: string, fieldId: string): Promise<void> {
    await apiClient.delete(`/api/content-types/${contentTypeId}/fields/${fieldId}`);
  },
};
