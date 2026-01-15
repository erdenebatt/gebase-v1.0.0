import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import { SystemInfo, RoleInfo, MenuItem, Organization } from "@/types";

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
        return (
          state.permissions.includes(code) ||
          state.platformPermissions.includes(code)
        );
      },

      getActiveMenus: () => {
        const state = get();
        return state.currentSystem ? state.menus : state.platformMenus;
      },
    }),
    {
      name: "portal-system-storage",
      storage: createJSONStorage(() => localStorage),
    }
  )
);
