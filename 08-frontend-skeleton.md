# Frontend Skeleton - Dashboard Template

Энэ skeleton-ийг баримтлан frontend-ээ хөгжүүлнэ. System switching, multi-system menu support нэмэгдсэн.

## 📁 File Structure

```
src/
├── app/
│   ├── layout.tsx
│   ├── globals.css
│   ├── (auth)/
│   │   ├── layout.tsx
│   │   └── login/page.tsx
│   └── (dashboard)/
│       ├── layout.tsx
│       ├── dashboard/page.tsx
│       ├── (platform)/              # Platform management pages
│       │   ├── users/page.tsx
│       │   ├── organizations/page.tsx
│       │   └── ...
│       ├── (dsl)/                   # DSL system pages
│       │   ├── schemas/page.tsx
│       │   ├── workflows/page.tsx
│       │   └── ...
│       └── (gateway)/               # Gateway system pages
│           ├── clients/page.tsx
│           ├── integrations/page.tsx
│           └── ...
├── components/
│   ├── layout/
│   │   ├── AppSidebar.tsx
│   │   ├── Header.tsx
│   │   └── SystemSwitcher.tsx       # System switching dropdown
│   └── ui/                          # shadcn components
├── store/
│   ├── authStore.ts
│   ├── systemStore.ts
│   └── uiStore.ts
├── lib/
│   ├── utils.ts
│   └── api.ts                       # API client with token handling
├── hooks/
│   └── usePermission.ts
└── types/
    └── index.ts
```

---

## 1. Types (src/types/index.ts)

```typescript
// ============ Core Types ============
export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  avatar_url?: string;
  language_code: string;
  default_system_id?: number;
}

export interface SystemInfo {
  id: number;
  code: string;
  name: string;
  description?: string;
  icon_name: string;
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
  children?: MenuItem[];
}

export interface Organization {
  id: number;
  name: string;
  short_name?: string;
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
  context_type: 'platform' | 'system';
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
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

// ============ Token Types ============
export type TokenType = 'platform' | 'system';

export interface TokenInfo {
  type: TokenType;
  token: string;
  expires_at: number;
}
```

---

## 2. Stores

### src/store/authStore.ts

```typescript
import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { User, SystemRole, TokenType } from '@/types';

interface AuthState {
  // User
  user: User | null;
  
  // Tokens
  platformToken: string | null;
  systemToken: string | null;
  refreshToken: string | null;
  currentTokenType: TokenType;
  
  // Available systems from login
  availableSystems: SystemRole[];
  
  // State
  isAuthenticated: boolean;
  isLoading: boolean;
  
  // Actions
  setAuth: (data: {
    user: User;
    platformToken: string;
    refreshToken: string;
    availableSystems: SystemRole[];
  }) => void;
  setSystemToken: (token: string) => void;
  clearSystemToken: () => void;
  setTokens: (tokens: { platformToken?: string; refreshToken?: string }) => void;
  logout: () => void;
  setLoading: (loading: boolean) => void;
  
  // Helpers
  getActiveToken: () => string | null;
  isInSystemContext: () => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      platformToken: null,
      systemToken: null,
      refreshToken: null,
      currentTokenType: 'platform',
      availableSystems: [],
      isAuthenticated: false,
      isLoading: true,

      setAuth: ({ user, platformToken, refreshToken, availableSystems }) =>
        set({
          user,
          platformToken,
          refreshToken,
          availableSystems,
          systemToken: null,
          currentTokenType: 'platform',
          isAuthenticated: true,
          isLoading: false,
        }),

      setSystemToken: (token) =>
        set({
          systemToken: token,
          currentTokenType: 'system',
        }),

      clearSystemToken: () =>
        set({
          systemToken: null,
          currentTokenType: 'platform',
        }),

      setTokens: (tokens) =>
        set((state) => ({
          ...state,
          ...tokens,
        })),

      logout: () =>
        set({
          user: null,
          platformToken: null,
          systemToken: null,
          refreshToken: null,
          currentTokenType: 'platform',
          availableSystems: [],
          isAuthenticated: false,
          isLoading: false,
        }),

      setLoading: (isLoading) => set({ isLoading }),

      getActiveToken: () => {
        const state = get();
        return state.currentTokenType === 'system' && state.systemToken
          ? state.systemToken
          : state.platformToken;
      },

      isInSystemContext: () => {
        const state = get();
        return state.currentTokenType === 'system' && !!state.systemToken;
      },
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => localStorage),
      onRehydrateStorage: () => (state) => state?.setLoading(false),
    }
  )
);
```

### src/store/systemStore.ts

```typescript
import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { SystemInfo, RoleInfo, MenuItem, Organization } from '@/types';

interface SystemState {
  // Current context
  currentSystem: SystemInfo | null;
  currentRole: RoleInfo | null;
  currentOrganization: Organization | null;
  
  // Permissions & Menus
  permissions: string[];
  menus: MenuItem[];
  
  // Platform menus (when not in system context)
  platformMenus: MenuItem[];
  platformPermissions: string[];
  
  // Actions
  setSystemContext: (data: {
    currentSystem: SystemInfo;
    currentRole: RoleInfo;
    currentOrganization?: Organization;
    permissions: string[];
    menus: MenuItem[];
  }) => void;
  
  setPlatformContext: (data: {
    permissions: string[];
    menus: MenuItem[];
  }) => void;
  
  clearSystemContext: () => void;
  
  // Helpers
  hasPermission: (code: string) => boolean;
  getActiveMenus: () => MenuItem[];
}

export const useSystemStore = create<SystemState>()(
  persist(
    (set, get) => ({
      currentSystem: null,
      currentRole: null,
      currentOrganization: null,
      permissions: [],
      menus: [],
      platformMenus: [],
      platformPermissions: [],

      setSystemContext: (data) =>
        set({
          currentSystem: data.currentSystem,
          currentRole: data.currentRole,
          currentOrganization: data.currentOrganization || null,
          permissions: data.permissions,
          menus: data.menus,
        }),

      setPlatformContext: (data) =>
        set({
          platformMenus: data.menus,
          platformPermissions: data.permissions,
        }),

      clearSystemContext: () =>
        set({
          currentSystem: null,
          currentRole: null,
          currentOrganization: null,
          permissions: [],
          menus: [],
        }),

      hasPermission: (code) => {
        const state = get();
        // Check both system and platform permissions
        return state.permissions.includes(code) || state.platformPermissions.includes(code);
      },

      getActiveMenus: () => {
        const state = get();
        // If in system context, return system menus; otherwise platform menus
        return state.currentSystem ? state.menus : state.platformMenus;
      },
    }),
    {
      name: 'system-storage',
      storage: createJSONStorage(() => localStorage),
    }
  )
);
```

### src/store/uiStore.ts

```typescript
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface UIState {
  sidebarOpen: boolean;
  systemSwitcherOpen: boolean;
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  toggleSystemSwitcher: () => void;
  setSystemSwitcherOpen: (open: boolean) => void;
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      sidebarOpen: true,
      systemSwitcherOpen: false,
      
      toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
      setSidebarOpen: (open) => set({ sidebarOpen: open }),
      toggleSystemSwitcher: () => set((state) => ({ systemSwitcherOpen: !state.systemSwitcherOpen })),
      setSystemSwitcherOpen: (open) => set({ systemSwitcherOpen: open }),
    }),
    { name: 'ui-storage' }
  )
);
```

---

## 3. API Client (src/lib/api.ts)

```typescript
import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/store/authStore';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

export const api = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - add token and device headers
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const { getActiveToken } = useAuthStore.getState();
    const token = getActiveToken();
    
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    // Add device headers
    const deviceUID = localStorage.getItem('device_uid') || generateDeviceUID();
    config.headers['X-Device-UID'] = deviceUID;
    config.headers['X-Platform'] = 'web';
    
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor - handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };
    
    // If 401 and not already retried, try refresh
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        const { refreshToken, setTokens, logout } = useAuthStore.getState();
        
        if (!refreshToken) {
          logout();
          window.location.href = '/login';
          return Promise.reject(error);
        }
        
        const response = await axios.post(`${API_BASE_URL}/api/v1/auth/refresh`, {
          refresh_token: refreshToken,
        });
        
        const { platform_token, refresh_token } = response.data.data;
        setTokens({ platformToken: platform_token, refreshToken: refresh_token });
        
        // Retry original request with new token
        originalRequest.headers.Authorization = `Bearer ${platform_token}`;
        return api(originalRequest);
      } catch (refreshError) {
        useAuthStore.getState().logout();
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }
    
    return Promise.reject(error);
  }
);

// Generate device UID
function generateDeviceUID(): string {
  const uid = `web_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  localStorage.setItem('device_uid', uid);
  return uid;
}

// API helper functions
export const authApi = {
  login: (data: { email: string; password: string; device_uid: string }) =>
    api.post('/auth/login', data),
  
  logout: () => api.post('/auth/logout'),
  
  refresh: (refreshToken: string) =>
    api.post('/auth/refresh', { refresh_token: refreshToken }),
  
  me: () => api.get('/auth/me'),
  
  availableSystems: () => api.get('/auth/available-systems'),
  
  switchSystem: (data: { system_code: string; organization_id?: number; role_id?: number }) =>
    api.post('/auth/switch-system', data),
  
  exitSystem: () => api.post('/systems/exit'),
  
  currentContext: () => api.get('/auth/current-context'),
};

export const platformApi = {
  // Users
  users: {
    list: (params?: Record<string, any>) => api.get('/platform/users', { params }),
    get: (id: number) => api.get(`/platform/users/${id}`),
    create: (data: any) => api.post('/platform/users', data),
    update: (id: number, data: any) => api.put(`/platform/users/${id}`, data),
    delete: (id: number) => api.delete(`/platform/users/${id}`),
  },
  
  // Organizations
  organizations: {
    list: (params?: Record<string, any>) => api.get('/platform/organizations', { params }),
    get: (id: number) => api.get(`/platform/organizations/${id}`),
    create: (data: any) => api.post('/platform/organizations', data),
    update: (id: number, data: any) => api.put(`/platform/organizations/${id}`, data),
    delete: (id: number) => api.delete(`/platform/organizations/${id}`),
  },
  
  // Roles
  roles: {
    list: (params?: Record<string, any>) => api.get('/platform/roles', { params }),
    get: (id: number) => api.get(`/platform/roles/${id}`),
    create: (data: any) => api.post('/platform/roles', data),
    update: (id: number, data: any) => api.put(`/platform/roles/${id}`, data),
  },
};

export const dslApi = {
  schemas: {
    list: (params?: Record<string, any>) => api.get('/systems/dsl/schemas', { params }),
    get: (id: number) => api.get(`/systems/dsl/schemas/${id}`),
    create: (data: any) => api.post('/systems/dsl/schemas', data),
    update: (id: number, data: any) => api.put(`/systems/dsl/schemas/${id}`, data),
    delete: (id: number) => api.delete(`/systems/dsl/schemas/${id}`),
    deploy: (id: number) => api.post(`/systems/dsl/schemas/${id}/deploy`),
  },
  
  workflows: {
    list: (params?: Record<string, any>) => api.get('/systems/dsl/workflows', { params }),
    get: (id: number) => api.get(`/systems/dsl/workflows/${id}`),
    create: (data: any) => api.post('/systems/dsl/workflows', data),
    execute: (id: number, data?: any) => api.post(`/systems/dsl/workflows/${id}/execute`, data),
    publish: (id: number) => api.post(`/systems/dsl/workflows/${id}/publish`),
  },
  
  // ... other DSL resources
};

export const gatewayApi = {
  clients: {
    list: (params?: Record<string, any>) => api.get('/systems/gateway/clients', { params }),
    get: (id: number) => api.get(`/systems/gateway/clients/${id}`),
    create: (data: any) => api.post('/systems/gateway/clients', data),
    regenerateSecret: (id: number) => api.post(`/systems/gateway/clients/${id}/regenerate-secret`),
    revoke: (id: number) => api.post(`/systems/gateway/clients/${id}/revoke`),
  },
  
  integrations: {
    list: (params?: Record<string, any>) => api.get('/systems/gateway/integrations', { params }),
    get: (id: number) => api.get(`/systems/gateway/integrations/${id}`),
    create: (data: any) => api.post('/systems/gateway/integrations', data),
    test: (id: number) => api.post(`/systems/gateway/integrations/${id}/test`),
  },
  
  // ... other Gateway resources
};
```

---

## 4. Permission Hook (src/hooks/usePermission.ts)

```typescript
import { useSystemStore } from '@/store/systemStore';

export function usePermission(permissionCode: string): boolean {
  const { hasPermission } = useSystemStore();
  return hasPermission(permissionCode);
}

export function usePermissions(permissionCodes: string[]): Record<string, boolean> {
  const { hasPermission } = useSystemStore();
  return permissionCodes.reduce((acc, code) => {
    acc[code] = hasPermission(code);
    return acc;
  }, {} as Record<string, boolean>);
}

// Component wrapper for permission-based rendering
export function PermissionGate({
  permission,
  children,
  fallback = null,
}: {
  permission: string;
  children: React.ReactNode;
  fallback?: React.ReactNode;
}) {
  const hasPermission = usePermission(permission);
  return hasPermission ? <>{children}</> : <>{fallback}</>;
}
```

---

## 5. System Switcher Component (src/components/layout/SystemSwitcher.tsx)

```typescript
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { ChevronDown, Check, LogOut, Settings } from 'lucide-react';
import * as LucideIcons from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/store/authStore';
import { useSystemStore } from '@/store/systemStore';
import { authApi } from '@/lib/api';
import { toast } from 'sonner';

export function SystemSwitcher() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  
  const { availableSystems, setSystemToken, clearSystemToken } = useAuthStore();
  const { currentSystem, currentRole, setSystemContext, clearSystemContext } = useSystemStore();

  const handleSwitchSystem = async (systemCode: string, roleId?: number, orgId?: number) => {
    if (currentSystem?.code === systemCode) return;
    
    setIsLoading(true);
    try {
      const response = await authApi.switchSystem({
        system_code: systemCode,
        organization_id: orgId,
        role_id: roleId,
      });
      
      const data = response.data.data;
      
      // Update tokens and context
      setSystemToken(data.system_token);
      setSystemContext({
        currentSystem: data.current_system,
        currentRole: data.current_role,
        currentOrganization: data.current_organization,
        permissions: data.permissions,
        menus: data.menus,
      });
      
      toast.success(`${data.current_system.name} системд нэвтэрлээ`);
      
      // Navigate to system dashboard
      router.push(`/${systemCode}`);
    } catch (error: any) {
      toast.error(error.response?.data?.error?.message || 'Систем сэлгэхэд алдаа гарлаа');
    } finally {
      setIsLoading(false);
    }
  };

  const handleExitSystem = async () => {
    setIsLoading(true);
    try {
      await authApi.exitSystem();
      
      clearSystemToken();
      clearSystemContext();
      
      toast.success('Системээс гарлаа');
      router.push('/platform');
    } catch (error: any) {
      toast.error('Алдаа гарлаа');
    } finally {
      setIsLoading(false);
    }
  };

  const getIcon = (iconName: string) => {
    const Icon = (LucideIcons as any)[iconName] || LucideIcons.Box;
    return Icon;
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button 
          variant="outline" 
          className="gap-2 min-w-[180px] justify-between"
          disabled={isLoading}
        >
          {currentSystem ? (
            <>
              {(() => {
                const Icon = getIcon(currentSystem.icon_name);
                return <Icon className="h-4 w-4" style={{ color: currentSystem.color }} />;
              })()}
              <span className="flex-1 text-left truncate">{currentSystem.name}</span>
            </>
          ) : (
            <>
              <Settings className="h-4 w-4" />
              <span className="flex-1 text-left">Platform</span>
            </>
          )}
          <ChevronDown className="h-4 w-4 opacity-50" />
        </Button>
      </DropdownMenuTrigger>
      
      <DropdownMenuContent align="start" className="w-[280px]">
        <DropdownMenuLabel className="text-xs text-muted-foreground">
          Систем сонгох
        </DropdownMenuLabel>
        
        {/* Platform option */}
        <DropdownMenuItem
          onClick={() => currentSystem && handleExitSystem()}
          className="gap-3"
        >
          <Settings className="h-4 w-4" />
          <div className="flex-1">
            <p className="font-medium">Platform</p>
            <p className="text-xs text-muted-foreground">Үндсэн удирдлага</p>
          </div>
          {!currentSystem && <Check className="h-4 w-4 text-primary" />}
        </DropdownMenuItem>
        
        <DropdownMenuSeparator />
        
        {/* Available systems */}
        {availableSystems.map((systemRole) => {
          const Icon = getIcon(systemRole.system.icon_name);
          const isActive = currentSystem?.code === systemRole.system.code;
          
          return (
            <DropdownMenuItem
              key={systemRole.system.id}
              onClick={() => handleSwitchSystem(
                systemRole.system.code,
                systemRole.roles[0]?.id,
                systemRole.roles[0]?.organization_id
              )}
              className="gap-3"
            >
              <Icon 
                className="h-4 w-4" 
                style={{ color: systemRole.system.color }} 
              />
              <div className="flex-1">
                <p className="font-medium">{systemRole.system.name}</p>
                <p className="text-xs text-muted-foreground">
                  {systemRole.roles[0]?.name}
                  {systemRole.roles[0]?.organization_name && (
                    <span> • {systemRole.roles[0].organization_name}</span>
                  )}
                </p>
              </div>
              {isActive && <Check className="h-4 w-4 text-primary" />}
            </DropdownMenuItem>
          );
        })}
        
        {/* Exit system option (only when in system context) */}
        {currentSystem && (
          <>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={handleExitSystem}
              className="gap-3 text-muted-foreground"
            >
              <LogOut className="h-4 w-4" />
              <span>Системээс гарах</span>
            </DropdownMenuItem>
          </>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
```

---

## 6. Updated Header with System Switcher (src/components/layout/Header.tsx)

```typescript
'use client';

import { Moon, Sun, PanelLeftClose, PanelLeft, LogOut, User, Bell } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { useAuthStore } from '@/store/authStore';
import { useSystemStore } from '@/store/systemStore';
import { useUIStore } from '@/store/uiStore';
import { SystemSwitcher } from './SystemSwitcher';

export function Header() {
  const [theme, setTheme] = useState<'light' | 'dark'>('light');
  const { user, logout } = useAuthStore();
  const { currentSystem, currentRole, clearSystemContext } = useSystemStore();
  const { sidebarOpen, toggleSidebar } = useUIStore();

  const handleLogout = () => {
    logout();
    clearSystemContext();
    window.location.href = '/login';
  };

  const getInitials = () => {
    const f = user?.first_name?.[0] || '';
    const l = user?.last_name?.[0] || '';
    return (f + l).toUpperCase() || 'U';
  };

  useEffect(() => {
    const saved = localStorage.getItem('theme') as 'light' | 'dark' | null;
    if (saved) {
      setTheme(saved);
      document.documentElement.classList.toggle('dark', saved === 'dark');
    }
  }, []);

  const toggleTheme = () => {
    const next = theme === 'light' ? 'dark' : 'light';
    setTheme(next);
    localStorage.setItem('theme', next);
    document.documentElement.classList.toggle('dark', next === 'dark');
  };

  return (
    <header className="sticky top-0 z-20 flex h-16 items-center justify-between border-b bg-white/80 px-4 backdrop-blur dark:bg-gray-900/80">
      {/* Left side */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={toggleSidebar}>
          {sidebarOpen ? <PanelLeftClose className="h-5 w-5" /> : <PanelLeft className="h-5 w-5" />}
        </Button>
        
        {/* System Switcher */}
        <SystemSwitcher />
        
        {/* Current context indicator */}
        {currentSystem && currentRole && (
          <div className="hidden md:block text-sm">
            <span className="text-muted-foreground">{currentRole.name}</span>
          </div>
        )}
      </div>

      {/* Right side */}
      <div className="flex items-center gap-2">
        {/* Notifications */}
        <Button variant="ghost" size="icon">
          <Bell className="h-5 w-5" />
        </Button>
        
        {/* Theme toggle */}
        <Button variant="ghost" size="icon" onClick={toggleTheme}>
          {theme === 'light' ? <Moon className="h-5 w-5" /> : <Sun className="h-5 w-5" />}
        </Button>

        {/* User menu */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="gap-2 px-2">
              <Avatar className="h-7 w-7">
                <AvatarImage src={user?.avatar_url} />
                <AvatarFallback className="bg-primary/10 text-xs text-primary">
                  {getInitials()}
                </AvatarFallback>
              </Avatar>
              <span className="hidden text-sm md:inline">{user?.first_name}</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56">
            <DropdownMenuLabel>
              <div className="flex flex-col">
                <span>{user?.first_name} {user?.last_name}</span>
                <span className="text-xs font-normal text-muted-foreground">{user?.email}</span>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <User className="mr-2 h-4 w-4" />
              Профайл
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleLogout} className="text-destructive">
              <LogOut className="mr-2 h-4 w-4" />
              Гарах
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
```

---

## 7. AppSidebar with Dynamic Menus (src/components/layout/AppSidebar.tsx)

```typescript
'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ChevronRight } from 'lucide-react';
import * as LucideIcons from 'lucide-react';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { ScrollArea } from '@/components/ui/scroll-area';
import { useSystemStore } from '@/store/systemStore';
import { useUIStore } from '@/store/uiStore';
import { cn } from '@/lib/utils';
import { MenuItem } from '@/types';

export function AppSidebar() {
  const pathname = usePathname();
  const { currentSystem, currentRole, getActiveMenus } = useSystemStore();
  const { sidebarOpen } = useUIStore();
  
  const menus = getActiveMenus();

  const getIcon = (name: string) => {
    return (LucideIcons as any)[name] || LucideIcons.Circle;
  };

  const renderMenuItem = (item: MenuItem) => {
    const Icon = getIcon(item.icon);
    const isActive = pathname === item.path || pathname.startsWith(item.path + '/');
    const hasChildren = item.children && item.children.length > 0;

    // Collapsible menu with children
    if (hasChildren && sidebarOpen) {
      return (
        <Collapsible key={item.id} defaultOpen={isActive}>
          <CollapsibleTrigger className="group flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-800">
            <Icon className="h-5 w-5 text-gray-500" />
            <span className="flex-1 text-left">{item.name}</span>
            <ChevronRight className="h-4 w-4 transition-transform group-data-[state=open]:rotate-90" />
          </CollapsibleTrigger>
          <CollapsibleContent>
            <div className="ml-4 mt-1 space-y-1 border-l pl-4">
              {item.children?.map((child) => (
                <Link
                  key={child.id}
                  href={child.path}
                  target={child.open_in_new_tab ? '_blank' : undefined}
                  className={cn(
                    'block rounded-lg px-3 py-2 text-sm',
                    pathname === child.path
                      ? 'bg-primary/10 text-primary'
                      : 'hover:bg-gray-100 dark:hover:bg-gray-800'
                  )}
                >
                  {child.name}
                </Link>
              ))}
            </div>
          </CollapsibleContent>
        </Collapsible>
      );
    }

    // Regular menu item
    const menuLink = (
      <Link
        href={item.path}
        target={item.open_in_new_tab ? '_blank' : undefined}
        className={cn(
          'flex items-center rounded-lg transition-colors',
          sidebarOpen ? 'gap-3 px-3 py-2' : 'justify-center p-2',
          isActive
            ? 'bg-primary/10 text-primary'
            : 'hover:bg-gray-100 dark:hover:bg-gray-800'
        )}
      >
        <Icon className="h-5 w-5" />
        {sidebarOpen && <span className="text-sm">{item.name}</span>}
      </Link>
    );

    // With tooltip when collapsed
    if (!sidebarOpen) {
      return (
        <Tooltip key={item.id} delayDuration={0}>
          <TooltipTrigger asChild>{menuLink}</TooltipTrigger>
          <TooltipContent side="right">{item.name}</TooltipContent>
        </Tooltip>
      );
    }

    return <div key={item.id}>{menuLink}</div>;
  };

  return (
    <aside
      className={cn(
        'fixed inset-y-0 left-0 z-30 flex flex-col border-r bg-white dark:bg-gray-900 transition-all duration-300',
        sidebarOpen ? 'w-[280px]' : 'w-16'
      )}
    >
      {/* Logo */}
      <div className="flex h-16 items-center border-b px-4">
        <Link href="/dashboard" className="flex items-center gap-3">
          <div 
            className="h-8 w-8 rounded-lg flex items-center justify-center"
            style={{ backgroundColor: currentSystem?.color || '#6366f1' }}
          >
            {currentSystem ? (
              (() => {
                const Icon = getIcon(currentSystem.icon_name);
                return <Icon className="h-5 w-5 text-white" />;
              })()
            ) : (
              <LucideIcons.Box className="h-5 w-5 text-white" />
            )}
          </div>
          {sidebarOpen && (
            <span className="font-semibold">
              {currentSystem?.name || 'Gebase'}
            </span>
          )}
        </Link>
      </div>

      {/* System & Role Info */}
      {sidebarOpen && currentSystem && (
        <div className="border-b px-4 py-3">
          <p className="text-sm font-medium">{currentSystem.name}</p>
          <p className="text-xs text-muted-foreground">{currentRole?.name}</p>
        </div>
      )}

      {/* Navigation */}
      <ScrollArea className="flex-1 py-4">
        <nav className={cn('space-y-1', sidebarOpen ? 'px-3' : 'px-2')}>
          {menus
            .filter((item) => item.is_visible !== false)
            .sort((a, b) => a.sequence - b.sequence)
            .map(renderMenuItem)}
        </nav>
      </ScrollArea>

      {/* Footer */}
      {sidebarOpen && (
        <div className="border-t p-4">
          <p className="text-xs text-muted-foreground text-center">
            Gebase Platform v1.0
          </p>
        </div>
      )}
    </aside>
  );
}
```

---

## 8. Login Page (src/app/(auth)/login/page.tsx)

```typescript
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useAuthStore } from '@/store/authStore';
import { useSystemStore } from '@/store/systemStore';
import { authApi } from '@/lib/api';
import { toast } from 'sonner';

const loginSchema = z.object({
  email: z.string().email('И-мэйл хаяг буруу байна'),
  password: z.string().min(1, 'Нууц үг оруулна уу'),
});

type LoginForm = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const { setAuth } = useAuthStore();
  const { setPlatformContext } = useSystemStore();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (data: LoginForm) => {
    setIsLoading(true);
    
    try {
      // Get or generate device UID
      let deviceUID = localStorage.getItem('device_uid');
      if (!deviceUID) {
        deviceUID = `web_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        localStorage.setItem('device_uid', deviceUID);
      }

      const response = await authApi.login({
        email: data.email,
        password: data.password,
        device_uid: deviceUID,
      });

      const result = response.data.data;

      // Set auth state
      setAuth({
        user: result.user,
        platformToken: result.platform_token,
        refreshToken: result.refresh_token,
        availableSystems: result.available_systems,
      });

      toast.success('Амжилттай нэвтэрлээ');

      // If user has default system, switch to it
      if (result.default_system_code) {
        router.push(`/${result.default_system_code}`);
      } else {
        router.push('/dashboard');
      }
    } catch (error: any) {
      const message = error.response?.data?.error?.message || 'Нэвтрэхэд алдаа гарлаа';
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-950 px-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 h-12 w-12 rounded-xl bg-primary flex items-center justify-center">
            <span className="text-2xl font-bold text-white">G</span>
          </div>
          <CardTitle className="text-2xl">Нэвтрэх</CardTitle>
          <CardDescription>Gebase Platform-д нэвтрэх</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email">И-мэйл</Label>
              <Input
                id="email"
                type="email"
                placeholder="example@email.com"
                {...register('email')}
              />
              {errors.email && (
                <p className="text-sm text-destructive">{errors.email.message}</p>
              )}
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="password">Нууц үг</Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                {...register('password')}
              />
              {errors.password && (
                <p className="text-sm text-destructive">{errors.password.message}</p>
              )}
            </div>

            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Нэвтрэх
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
```

---

## 🚀 Setup Commands

```bash
# 1. Create Next.js project
npx create-next-app@latest gebase-admin --typescript --tailwind --eslint --app --src-dir

# 2. Install dependencies
cd gebase-admin
npm install zustand axios sonner lucide-react clsx tailwind-merge @hookform/resolvers zod react-hook-form

# 3. Init shadcn
npx shadcn@latest init

# 4. Add required components
npx shadcn@latest add button input label card avatar dropdown-menu tooltip scroll-area collapsible sonner

# 5. Create folders
mkdir -p src/{store,types,components/layout,hooks,lib}

# 6. Copy files and start
npm run dev
```

---

## Key Features

1. **Two-Token Strategy** - Platform token + System token
2. **System Switcher** - Header дахь dropdown-оор систем сэлгэх
3. **Dynamic Menus** - Системээс хамаарч sidebar меню өөрчлөгдөнө
4. **Permission Hook** - `usePermission('dsl.schema.view')` хэлбэрээр эрх шалгах
5. **API Client** - Token auto-refresh, device headers
6. **Zustand Stores** - Auth, System, UI state management
