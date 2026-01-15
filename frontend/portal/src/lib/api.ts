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

  getPermissions: () => api.get("/auth/permissions"),

  getMenus: () => api.get("/auth/menus"),
};

// Portal-specific APIs
export const portalApi = {
  // Profile
  profile: {
    get: () => api.get("/profile"),
    update: (data: unknown) => api.put("/profile", data),
    changePassword: (data: { old_password: string; new_password: string }) =>
      api.post("/profile/change-password", data),
  },

  // Notifications
  notifications: {
    list: (params?: Record<string, unknown>) =>
      api.get("/notifications", { params }),
    markAsRead: (id: number) => api.put(`/notifications/${id}/read`),
    markAllAsRead: () => api.put("/notifications/read-all"),
  },

  // Sessions
  sessions: {
    list: () => api.get("/sessions"),
    revoke: (id: number) => api.delete(`/sessions/${id}`),
    revokeAll: () => api.delete("/sessions/all"),
  },

  // Devices
  devices: {
    list: () => api.get("/devices"),
    deactivate: (id: number) => api.delete(`/devices/${id}`),
  },
};
