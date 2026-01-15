import { create } from "zustand";
import { persist } from "zustand/middleware";

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

      toggleSidebar: () =>
        set((state) => ({ sidebarOpen: !state.sidebarOpen })),
      setSidebarOpen: (open) => set({ sidebarOpen: open }),
      toggleSystemSwitcher: () =>
        set((state) => ({ systemSwitcherOpen: !state.systemSwitcherOpen })),
      setSystemSwitcherOpen: (open) => set({ systemSwitcherOpen: open }),
    }),
    { name: "ui-storage" }
  )
);
