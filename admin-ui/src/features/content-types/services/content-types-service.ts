import { apiClient } from "@/lib/api";
import { RawContentType, ContentType, mapContentType } from "@/types/api";

export const contentTypesService = {
  async getContentTypes(): Promise<ContentType[]> {
    const response = await apiClient.get<{ data: RawContentType[] }>("/api/content-types");
    return response.data.map(mapContentType);
  },

  async getContentType(id: string): Promise<ContentType> {
    const response = await apiClient.get<{ data: RawContentType }>(`/api/content-types/${id}`);
    return mapContentType(response.data);
  },

  async createContentType(data: {
    name: string;
    slug: string;
    description?: string;
  }): Promise<ContentType> {
    const response = await apiClient.post<{ data: RawContentType }>(
      "/api/content-types",
      data
    );
    return mapContentType(response.data);
  },

  async updateContentType(
    id: string,
    data: Partial<{ name: string; slug: string; description?: string }>
  ): Promise<ContentType> {
    const response = await apiClient.put<{ data: RawContentType }>(
      `/api/content-types/${id}`,
      data
    );
    return mapContentType(response.data);
  },

  // 🆕 Delete Content Type
  async deleteContentType(id: string): Promise<void> {
    await apiClient.delete(`/api/content-types/${id}`);
  },
};