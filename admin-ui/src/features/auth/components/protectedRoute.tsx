import { Navigate } from "react-router-dom";
import { authService } from "@/features/auth/services/auth-service";

interface Props {
  children: JSX.Element;
}

export function ProtectedRoute({ children }: Props) {
  const isAuth = authService.isAuthenticated();
  return isAuth ? children : <Navigate to="/login" replace />;
}

export function GuestRoute({ children }: { children: JSX.Element }) {
  return authService.isAuthenticated() ? <Navigate to="/admin" replace /> : children;
}