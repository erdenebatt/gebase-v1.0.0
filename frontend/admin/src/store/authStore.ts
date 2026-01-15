import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import { User, SystemRole, TokenType } from "@/types";

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
  setTokens: (tokens: {
    platformToken?: string;
    refreshToken?: string;
  }) => void;
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
      currentTokenType: "platform",
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
          currentTokenType: "platform",
          isAuthenticated: true,
          isLoading: false,
        }),

      setSystemToken: (token) =>
        set({
          systemToken: token,
          currentTokenType: "system",
        }),

      clearSystemToken: () =>
        set({
          systemToken: null,
          currentTokenType: "platform",
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
          currentTokenType: "platform",
          availableSystems: [],
          isAuthenticated: false,
          isLoading: false,
        }),

      setLoading: (isLoading) => set({ isLoading }),

      getActiveToken: () => {
        const state = get();
        return state.currentTokenType === "system" && state.systemToken
          ? state.systemToken
          : state.platformToken;
      },

      isInSystemContext: () => {
        const state = get();
        return state.currentTokenType === "system" && !!state.systemToken;
      },
    }),
    {
      name: "auth-storage",
      storage: createJSONStorage(() => localStorage),
      onRehydrateStorage: () => (state) => state?.setLoading(false),
    }
  )
);
