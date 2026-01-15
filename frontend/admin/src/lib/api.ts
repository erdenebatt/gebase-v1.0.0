import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";
import { useAuthStore } from "@/store/authStore";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

export const api = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  headers: {
    "Content-Type": "application/json",
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
    const deviceUID = localStorage.getItem("device_uid") || generateDeviceUID();
    config.headers["X-Device-UID"] = deviceUID;
    config.headers["X-Platform"] = "web";

    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor - handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean;
    };

    // If 401 and not already retried, try refresh
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const { refreshToken, setTokens, logout } = useAuthStore.getState();

        if (!refreshToken) {
          logout();
          window.location.href = "/login";
          return Promise.reject(error);
        }

        const response = await axios.post(
          `${API_BASE_URL}/api/v1/auth/refresh`,
          {
            refresh_token: refreshToken,
          }
        );

        const { platform_token, refresh_token } = response.data.data;
        setTokens({
          platformToken: platform_token,
          refreshToken: refresh_token,
        });

        // Retry original request with new token
        originalRequest.headers.Authorization = `Bearer ${platform_token}`;
        return api(originalRequest);
      } catch (refreshError) {
        useAuthStore.getState().logout();
        window.location.href = "/login";
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

// Generate device UID
function generateDeviceUID(): string {
  const uid = `web_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  localStorage.setItem("device_uid", uid);
  return uid;
}

// API helper functions
export const authApi = {
  login: (data: { email: string; password: string; device_uid: string }) =>
    api.post("/auth/login", data),

  logout: () => api.post("/auth/logout"),

  refresh: (refreshToken: string) =>
    api.post("/auth/refresh", { refresh_token: refreshToken }),

  me: () => api.get("/auth/me"),

  availableSystems: () => api.get("/auth/available-systems"),

  switchSystem: (data: {
    system_code: string;
    organization_id?: number;
    role_id?: number;
  }) => api.post("/auth/switch-system", data),

  exitSystem: () => api.post("/auth/exit-system"),

  currentContext: () => api.get("/auth/current-context"),

  getPermissions: () => api.get("/auth/permissions"),

  getMenus: () => api.get("/auth/menus"),
};

export const usersApi = {
  list: (params?: Record<string, unknown>) => api.get("/users", { params }),
  get: (id: number) => api.get(`/users/${id}`),
  create: (data: unknown) => api.post("/users", data),
  update: (id: number, data: unknown) => api.put(`/users/${id}`, data),
  delete: (id: number) => api.delete(`/users/${id}`),
  getRoles: (id: number) => api.get(`/users/${id}/roles`),
  assignRoles: (id: number, data: unknown) =>
    api.put(`/users/${id}/roles`, data),
};

export const organizationsApi = {
  list: (params?: Record<string, unknown>) =>
    api.get("/organizations", { params }),
  get: (id: number) => api.get(`/organizations/${id}`),
  create: (data: unknown) => api.post("/organizations", data),
  update: (id: number, data: unknown) => api.put(`/organizations/${id}`, data),
  delete: (id: number) => api.delete(`/organizations/${id}`),
  getTypes: () => api.get("/organizations/types"),
  getChildren: (id: number) => api.get(`/organizations/${id}/children`),
  getSystems: (id: number) => api.get(`/organizations/${id}/systems`),
};

export const systemsApi = {
  list: () => api.get("/systems"),
  get: (id: number) => api.get(`/systems/${id}`),
  create: (data: unknown) => api.post("/systems", data),
  update: (id: number, data: unknown) => api.put(`/systems/${id}`, data),
  delete: (id: number) => api.delete(`/systems/${id}`),
  getModules: (id: number) => api.get(`/systems/${id}/modules`),
  getMenus: (id: number) => api.get(`/systems/${id}/menus`),
};

export const rolesApi = {
  list: (params?: Record<string, unknown>) => api.get("/roles", { params }),
  get: (id: number) => api.get(`/roles/${id}`),
  create: (data: unknown) => api.post("/roles", data),
  update: (id: number, data: unknown) => api.put(`/roles/${id}`, data),
  delete: (id: number) => api.delete(`/roles/${id}`),
  getPermissions: (id: number) => api.get(`/roles/${id}/permissions`),
  assignPermissions: (id: number, data: unknown) =>
    api.put(`/roles/${id}/permissions`, data),
  getMenus: (id: number) => api.get(`/roles/${id}/menus`),
  assignMenus: (id: number, data: unknown) =>
    api.put(`/roles/${id}/menus`, data),
};

export const devicesApi = {
  register: (data: unknown) => api.post("/devices/register", data),
  heartbeat: () => api.post("/devices/heartbeat"),
  list: (params?: Record<string, unknown>) => api.get("/devices", { params }),
  get: (id: number) => api.get(`/devices/${id}`),
  update: (id: number, data: unknown) => api.put(`/devices/${id}`, data),
  deactivate: (id: number) => api.delete(`/devices/${id}`),
};

export const menusApi = {
  list: (params?: Record<string, unknown>) => api.get("/menus", { params }),
  getTree: (systemId: number) => api.get("/menus/tree", { params: { system_id: systemId } }),
  get: (id: number) => api.get(`/menus/${id}`),
  create: (data: unknown) => api.post("/menus", data),
  update: (id: number, data: unknown) => api.put(`/menus/${id}`, data),
  delete: (id: number) => api.delete(`/menus/${id}`),
};

export const dslApi = {
  schemas: {
    list: (params?: Record<string, unknown>) =>
      api.get("/dsl/schemas", { params }),
    get: (id: number) => api.get(`/dsl/schemas/${id}`),
    create: (data: unknown) => api.post("/dsl/schemas", data),
    update: (id: number, data: unknown) => api.put(`/dsl/schemas/${id}`, data),
    delete: (id: number) => api.delete(`/dsl/schemas/${id}`),
    deploy: (id: number) => api.post(`/dsl/schemas/${id}/deploy`),
  },
  workflows: {
    list: (params?: Record<string, unknown>) =>
      api.get("/dsl/workflows", { params }),
    get: (id: number) => api.get(`/dsl/workflows/${id}`),
    create: (data: unknown) => api.post("/dsl/workflows", data),
    update: (id: number, data: unknown) =>
      api.put(`/dsl/workflows/${id}`, data),
    delete: (id: number) => api.delete(`/dsl/workflows/${id}`),
    execute: (id: number, data?: unknown) =>
      api.post(`/dsl/workflows/${id}/execute`, data),
    publish: (id: number) => api.post(`/dsl/workflows/${id}/publish`),
  },
  rules: {
    list: (params?: Record<string, unknown>) =>
      api.get("/dsl/rules", { params }),
    get: (id: number) => api.get(`/dsl/rules/${id}`),
    create: (data: unknown) => api.post("/dsl/rules", data),
    update: (id: number, data: unknown) => api.put(`/dsl/rules/${id}`, data),
    delete: (id: number) => api.delete(`/dsl/rules/${id}`),
    test: (id: number, data?: unknown) =>
      api.post(`/dsl/rules/${id}/test`, data),
  },
  templates: {
    list: (params?: Record<string, unknown>) =>
      api.get("/dsl/templates", { params }),
    get: (id: number) => api.get(`/dsl/templates/${id}`),
    create: (data: unknown) => api.post("/dsl/templates", data),
    update: (id: number, data: unknown) =>
      api.put(`/dsl/templates/${id}`, data),
    delete: (id: number) => api.delete(`/dsl/templates/${id}`),
    render: (id: number, data?: unknown) =>
      api.post(`/dsl/templates/${id}/render`, data),
  },
  functions: {
    list: (params?: Record<string, unknown>) =>
      api.get("/dsl/functions", { params }),
    get: (id: number) => api.get(`/dsl/functions/${id}`),
    create: (data: unknown) => api.post("/dsl/functions", data),
    update: (id: number, data: unknown) =>
      api.put(`/dsl/functions/${id}`, data),
    delete: (id: number) => api.delete(`/dsl/functions/${id}`),
    test: (id: number, data?: unknown) =>
      api.post(`/dsl/functions/${id}/test`, data),
  },
  variables: {
    list: (params?: Record<string, unknown>) =>
      api.get("/dsl/variables", { params }),
    get: (id: number) => api.get(`/dsl/variables/${id}`),
    create: (data: unknown) => api.post("/dsl/variables", data),
    update: (id: number, data: unknown) =>
      api.put(`/dsl/variables/${id}`, data),
    delete: (id: number) => api.delete(`/dsl/variables/${id}`),
  },
};
