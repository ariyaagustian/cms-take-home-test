// API Types for CMS Backend
export interface User {
  id: string;
  email: string;
  name: string;
  roles: Role[];
  created_at: string;
  updated_at: string;
}

export interface Role {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

// --- RAW TYPES (langsung dari backend)
export interface RawContentType {
  ID: string;
  Name: string;
  Slug: string;
  CreatedAt: string;
  UpdatedAt: string;
  Fields: RawContentField[];
}

export interface RawContentField {
  ID: string;
  ContentTypeID: string;
  Name: string;
  Kind: string;
  Options: Record<string, unknown>;
}

// --- DOMAIN TYPES (camelCase, dipakai frontend)
export interface ContentType {
  id: string;
  name: string;
  slug: string;
  createdAt: string;
  updatedAt: string;
  fields: ContentField[];
}

export interface ContentField {
  id: string;
  contentTypeId: string;
  name: string;
  kind: string;
  options: Record<string, unknown>;
}

// --- MAPPER
export function mapContentType(raw: RawContentType): ContentType {
  return {
    id: raw.ID,
    name: raw.Name,
    slug: raw.Slug,
    createdAt: raw.CreatedAt,
    updatedAt: raw.UpdatedAt,
    fields: Array.isArray(raw.Fields)
      ? raw.Fields.map((f) => ({
          id: f.ID,
          contentTypeId: f.ContentTypeID,
          name: f.Name,
          kind: f.Kind,
          options: f.Options,
        }))
      : [], // ✅ kalau null → jadi array kosong
  };
}

export interface Entry {
  id: string;
  content_type_id: string;
  content_type?: ContentType;
  data: Record<string, unknown>;
  status: 'draft' | 'published' | 'archived';
  version: number;
  published_at?: string;
  created_by: string;
  updated_by: string;
  created_at: string;
  updated_at: string;
}

export interface MediaAsset {
  id: string;
  filename: string;
  original_filename: string;
  mime_type: string;
  size: number;
  url: string;
  alt_text?: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface AuditLog {
  id: string;
  user_id: string;
  action: string;
  resource_type: string;
  resource_id: string;
  old_values?: Record<string, unknown>;
  new_values?: Record<string, unknown>;
  created_at: string;
}

// API Request/Response Types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
  success: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    current_page: number;
    total_pages: number;
    total_items: number;
    items_per_page: number;
  };
  success: boolean;
}

export interface CreateContentTypeRequest {
  name: string;
  slug: string;
  description: string;
}

export interface CreateFieldRequest {
  name: string;
  required: boolean;
  unique: boolean;
  default_value?: unknown;
  validation_rules?: Record<string, unknown>;
}

export interface CreateEntryRequest {
  data: Record<string, unknown>;
  status?: Entry['status'];
}

export interface UpdateEntryRequest {
  data?: Record<string, unknown>;
  status?: Entry['status'];
}

interface UserFromToken {
  id: string;
  email: string;
  name: string;
  role?: string;   // single role
  roles?: Role[];  // multi role
}