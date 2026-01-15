// ============ Core Types ============
export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  family_name?: string;
  avatar_url?: string;
  language_code: string;
  organization_id?: number;
  default_system_id?: number;
}

export interface SystemInfo {
  id: number;
  code: string;
  name: string;
  description?: string;
  icon_name: string;
  icon_url?: string;
  color: string;
  sequence: number;
}

export interface RoleInfo {
  id: number;
  code: string;
  name: string;
  organization_id?: number;
  organization_name?: string;
}

export interface SystemRole {
  system: SystemInfo;
  roles: RoleInfo[];
}

export interface MenuItem {
  id: number;
  code: string;
  name: string;
  icon: string;
  path: string;
  sequence: number;
  is_visible: boolean;
  open_in_new_tab?: boolean;
  children?: MenuItem[];
}

export interface Organization {
  id: number;
  name: string;
  short_name?: string;
  reg_no?: string;
}

// ============ Auth Types ============
export interface LoginRequest {
  email: string;
  password: string;
  device_uid: string;
}

export interface LoginResponse {
  platform_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
  available_systems: SystemRole[];
  default_system_code?: string;
}

export interface SwitchSystemRequest {
  system_code: string;
  organization_id?: number;
  role_id?: number;
}

export interface SwitchSystemResponse {
  system_token: string;
  expires_in: number;
  current_system: SystemInfo;
  current_role: RoleInfo;
  current_organization?: Organization;
  permissions: string[];
  menus: MenuItem[];
}

export interface CurrentContextResponse {
  context_type: "platform" | "system";
  user: User;
  current_system?: SystemInfo;
  current_role?: RoleInfo;
  current_organization?: Organization;
  permissions: string[];
  menus: MenuItem[];
  available_systems: SystemRole[];
}

// ============ API Types ============
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: Array<{ field: string; message: string }>;
  };
  meta?: PaginationMeta;
}

export interface PaginationMeta {
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

// ============ Token Types ============
export type TokenType = "platform" | "system";

export interface TokenInfo {
  type: TokenType;
  token: string;
  expires_at: number;
}

// ============ Role & Permission Types ============
export interface Role {
  id: number;
  code: string;
  name: string;
  description?: string;
  system_id: number;
  is_system: boolean;
  is_active: boolean;
}

export interface Permission {
  id: number;
  code: string;
  name: string;
  description?: string;
  system_id: number;
  module_id: number;
  action_id?: number;
  is_active: boolean;
}

// ============ Device Types ============
export interface Device {
  id: number;
  device_uid: string;
  name: string;
  platform: string;
  os_version?: string;
  app_version?: string;
  is_registered: boolean;
  is_active: boolean;
  last_heartbeat?: string;
}

// ============ Session Types ============
export interface Session {
  id: number;
  session_token: string;
  user_id: number;
  device_id: number;
  system_id?: number;
  ip_address?: string;
  user_agent?: string;
  is_active: boolean;
  expires_at: string;
  last_activity?: string;
}
