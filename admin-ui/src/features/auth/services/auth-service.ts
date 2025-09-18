import { apiClient } from "@/lib/api";
import { AuthResponse, LoginRequest, RegisterRequest, ApiResponse } from "@/types/api";

export const authService = {
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>(
      '/api/auth/login',
      credentials
    );

    if (response.token) {
      localStorage.setItem('cms_token', response.token);
    }

    return response;
  },

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post<ApiResponse<AuthResponse>>('/api/auth/register', userData);
    
    if (response.success && response.data.token) {
      localStorage.setItem('cms_token', response.data.token);
    }
    
    return response.data;
  },

  logout(): void {
    localStorage.removeItem('cms_token');
  },

  getToken(): string | null {
    return localStorage.getItem('cms_token');
  },

  isAuthenticated(): boolean {
    return !!this.getToken();
  }
};