import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { authService } from "../services/auth-service";
import { LoginRequest, RegisterRequest } from "@/types/api";
import { toast } from "@/hooks/use-toast";

export const useAuth = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const loginMutation = useMutation({
    mutationFn: (credentials: LoginRequest) => authService.login(credentials),
    onSuccess: (data) => {
      toast({
        title: "Welcome back!",
        description: `Logged in as ${data.user.name}`,
      });
      navigate("/admin", { replace: true });
    },
    onError: (error: Error) => {
      toast({
        title: "Login failed",
        description: error.message,
        variant: "destructive",
      });
    },
  });

  const registerMutation = useMutation({
    mutationFn: (userData: RegisterRequest) => authService.register(userData),
    onSuccess: (data) => {
      toast({
        title: "Account created!",
        description: `Welcome ${data.user.name}`,
      });
      navigate("/admin", { replace: true });
    },
    onError: (error: Error) => {
      toast({
        title: "Registration failed",
        description: error.message,
        variant: "destructive",
      });
    },
  });

  const logout = () => {
    authService.logout();
    queryClient.clear();
    toast({
      title: "Logged out",
      description: "You have been successfully logged out",
    });
    navigate("/login", { replace: true });
  };

  return {
    login: loginMutation.mutate,
    register: registerMutation.mutate,
    logout,
    isLoggingIn: loginMutation.isPending,
    isRegistering: registerMutation.isPending,
    isAuthenticated: authService.isAuthenticated(),
  };
};
