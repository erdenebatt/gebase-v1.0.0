"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { ChevronRight, HelpCircle, Settings } from "lucide-react";
import * as LucideIcons from "lucide-react";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useAuthStore } from "@/store/authStore";
import { useSystemStore } from "@/store/systemStore";
import { cn } from "@/lib/utils";
import { MenuItem } from "@/types";
import { toast } from "sonner";

export function AppSidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);

  const { availableSystems, setSystemToken, clearSystemToken } = useAuthStore();
  const {
    currentSystem,
    currentRole,
    getActiveMenus,
    setSystemContext,
    clearSystemContext,
  } = useSystemStore();

  const menus = getActiveMenus();

  const getIcon = (name: string) => {
    const icons = LucideIcons as unknown as Record<
      string,
      React.ComponentType<{ className?: string; style?: React.CSSProperties }>
    >;
    return icons[name] || LucideIcons.Circle;
  };

  const handleSwitchSystem = async (systemCode: string) => {
    if (currentSystem?.code === systemCode) return;

    setIsLoading(true);
    try {
      const deviceUID = localStorage.getItem("device_uid") || "";
      const token = useAuthStore.getState().platformToken;

      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000"}/api/v1/auth/switch-system`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
            "X-Device-UID": deviceUID,
            "X-Platform": "web",
          },
          body: JSON.stringify({ system_code: systemCode }),
        }
      );

      const result = await response.json();

      if (result.success && result.data) {
        setSystemToken(result.data.system_token);
        setSystemContext({
          currentSystem: result.data.current_system,
          currentRole: result.data.current_role,
          currentOrganization: result.data.current_organization,
          permissions: result.data.permissions || [],
          menus: result.data.menus || [],
        });

        toast.success(`${result.data.current_system.name} системд нэвтэрлээ`);
        router.push("/dashboard");
      }
    } catch (error) {
      console.error("Failed to switch system:", error);
      toast.error("Систем сэлгэхэд алдаа гарлаа");
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      clearSystemToken();
      clearSystemContext();
      useAuthStore.getState().logout();
      router.push("/login");
    } catch {
      toast.error("Гарахад алдаа гарлаа");
    }
  };

  const renderMenuItem = (item: MenuItem) => {
    const Icon = getIcon(item.icon);
    const isActive =
      pathname === item.path || pathname.startsWith(item.path + "/");
    const hasChildren = item.children && item.children.length > 0;

    // Collapsible menu with children
    if (hasChildren) {
      return (
        <Collapsible key={item.id} defaultOpen={isActive}>
          <CollapsibleTrigger className="group flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-800">
            <Icon className="h-4 w-4 text-gray-500" />
            <span className="flex-1 text-left text-gray-700 dark:text-gray-300">
              {item.name}
            </span>
            <ChevronRight className="h-4 w-4 text-gray-400 transition-transform group-data-[state=open]:rotate-90" />
          </CollapsibleTrigger>
          <CollapsibleContent>
            <div className="ml-4 mt-1 space-y-1 border-l border-gray-200 dark:border-gray-700 pl-4">
              {item.children?.map((child) => (
                <Link
                  key={child.id}
                  href={child.path}
                  className={cn(
                    "block rounded-lg px-3 py-2 text-sm transition-colors",
                    pathname === child.path
                      ? "bg-primary/10 text-primary font-medium"
                      : "text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-900 dark:hover:text-gray-100"
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

    // Regular menu item (no children)
    return (
      <Link
        key={item.id}
        href={item.path}
        className={cn(
          "flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors",
          isActive
            ? "bg-primary/10 text-primary font-medium"
            : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800"
        )}
      >
        <Icon className="h-4 w-4" />
        <span>{item.name}</span>
      </Link>
    );
  };

  return (
    <aside className="fixed inset-y-0 left-0 z-30 flex">
      {/* System Rail - Narrow left bar with system icons */}
      <div className="w-16 flex flex-col bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800">
        {/* Logo */}
        <div className="h-16 flex items-center justify-center border-b border-gray-200 dark:border-gray-800">
          <div className="h-10 w-10 rounded-xl bg-primary flex items-center justify-center">
            <span className="text-xl font-bold text-white">G</span>
          </div>
        </div>

        {/* System Icons */}
        <ScrollArea className="flex-1 py-3">
          <div className="flex flex-col items-center gap-2 px-2">
            {availableSystems.map((systemRole) => {
              const Icon = getIcon(systemRole.system.icon_name);
              const isActive = currentSystem?.code === systemRole.system.code;

              return (
                <Tooltip key={systemRole.system.id} delayDuration={0}>
                  <TooltipTrigger asChild>
                    <button
                      onClick={() => handleSwitchSystem(systemRole.system.code)}
                      disabled={isLoading}
                      className={cn(
                        "relative w-12 h-12 rounded-xl flex items-center justify-center transition-all duration-200",
                        isActive
                          ? "bg-gray-100 dark:bg-white/10"
                          : "hover:bg-gray-100 dark:hover:bg-white/5 opacity-60 hover:opacity-100"
                      )}
                    >
                      <Icon
                        className="h-6 w-6"
                        style={{
                          color: isActive
                            ? systemRole.system.color
                            : "#9ca3af",
                        }}
                      />
                    </button>
                  </TooltipTrigger>
                  <TooltipContent side="right" className="font-medium">
                    {systemRole.system.name}
                  </TooltipContent>
                </Tooltip>
              );
            })}
          </div>
        </ScrollArea>

        {/* Bottom buttons */}
        <div className="flex flex-col items-center gap-2 p-2 border-t border-gray-200 dark:border-gray-800">
          <Tooltip delayDuration={0}>
            <TooltipTrigger asChild>
              <Link
                href="/settings"
                className={cn(
                  "w-10 h-10 rounded-xl flex items-center justify-center transition-colors",
                  pathname === "/settings"
                    ? "bg-gray-100 dark:bg-white/10 text-gray-900 dark:text-white"
                    : "text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-white/5"
                )}
              >
                <Settings className="h-5 w-5" />
              </Link>
            </TooltipTrigger>
            <TooltipContent side="right">Тохиргоо</TooltipContent>
          </Tooltip>
          <Tooltip delayDuration={0}>
            <TooltipTrigger asChild>
              <Link
                href="/help"
                className={cn(
                  "w-10 h-10 rounded-xl flex items-center justify-center transition-colors",
                  pathname === "/help"
                    ? "bg-gray-100 dark:bg-white/10 text-gray-900 dark:text-white"
                    : "text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-white/5"
                )}
              >
                <HelpCircle className="h-5 w-5" />
              </Link>
            </TooltipTrigger>
            <TooltipContent side="right">Тусламж</TooltipContent>
          </Tooltip>
        </div>
      </div>

      {/* Menu Panel - Wider right bar with menus */}
      <div className="w-60 flex flex-col bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800">
        {/* Current System Header */}
        <div className="h-16 flex items-center px-4 border-b border-gray-200 dark:border-gray-800">
          {currentSystem ? (
            <div>
              <h2 className="font-semibold text-gray-900 dark:text-white">
                {currentSystem.name}
              </h2>
              <p className="text-xs text-gray-500 dark:text-gray-400">
                {currentRole?.name}
              </p>
            </div>
          ) : (
            <h2 className="font-semibold text-gray-900 dark:text-white">
              Gebase Platform
            </h2>
          )}
        </div>

        {/* Menu Navigation */}
        <ScrollArea className="flex-1 py-4">
          <nav className="space-y-1 px-3">
            {menus
              .filter((item) => item.is_visible !== false)
              .sort((a, b) => a.sequence - b.sequence)
              .map(renderMenuItem)}
          </nav>
        </ScrollArea>

        {/* Footer */}
        <div className="border-t border-gray-200 dark:border-gray-800 p-4">
          <p className="text-xs text-gray-400 text-center">
            Gebase Platform v1.0
          </p>
        </div>
      </div>
    </aside>
  );
}
