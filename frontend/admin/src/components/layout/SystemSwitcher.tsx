"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { ChevronDown, Check, LogOut, Settings } from "lucide-react";
import * as LucideIcons from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/store/authStore";
import { useSystemStore } from "@/store/systemStore";
import { authApi } from "@/lib/api";
import { toast } from "sonner";

export function SystemSwitcher() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);

  const { availableSystems, setSystemToken, clearSystemToken } = useAuthStore();
  const { currentSystem, setSystemContext, clearSystemContext } =
    useSystemStore();

  const handleSwitchSystem = async (
    systemCode: string,
    roleId?: number,
    orgId?: number
  ) => {
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

      // Navigate to dashboard after switching system
      router.push("/dashboard");
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: { message?: string } } } };
      toast.error(
        err.response?.data?.error?.message || "Систем сэлгэхэд алдаа гарлаа"
      );
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

      toast.success("Системээс гарлаа");
      router.push("/dashboard");
    } catch {
      toast.error("Алдаа гарлаа");
    } finally {
      setIsLoading(false);
    }
  };

  const getIcon = (iconName: string) => {
    const icons = LucideIcons as unknown as Record<string, React.ComponentType<{ className?: string; style?: React.CSSProperties }>>;
    return icons[iconName] || LucideIcons.Box;
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
                return (
                  <Icon
                    className="h-4 w-4"
                    style={{ color: currentSystem.color }}
                  />
                );
              })()}
              <span className="flex-1 text-left truncate">
                {currentSystem.name}
              </span>
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
              onClick={() =>
                handleSwitchSystem(
                  systemRole.system.code,
                  systemRole.roles[0]?.id,
                  systemRole.roles[0]?.organization_id
                )
              }
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
