import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api";

// wrapper kecil buat fetch count
async function fetchCount(endpoint: string): Promise<number> {
  const res = await apiClient.get<{ data: unknown[] }>(endpoint);
  return res.data.length;
}

export function useDashboardStats() {
  const contentTypesQuery = useQuery({
    queryKey: ["content-types"],
    queryFn: () => fetchCount("/api/content-types"),
  });


  const mediaQuery = useQuery({
    queryKey: ["media"],
    queryFn: () => fetchCount("/api/media"),
  });

  const usersQuery = useQuery({
    queryKey: ["users"],
    queryFn: () => fetchCount("/api/admin/users"),
  });

  return {
    contentTypes: contentTypesQuery.data ?? 0,
    media: mediaQuery.data ?? 0,
    users: usersQuery.data ?? 0,
    isLoading:
      contentTypesQuery.isLoading ||
      mediaQuery.isLoading ||
      usersQuery.isLoading,
  };
}
