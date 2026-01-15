import { useSystemStore } from "@/store/systemStore";
import React from "react";

export function usePermission(permissionCode: string): boolean {
  const { hasPermission } = useSystemStore();
  return hasPermission(permissionCode);
}

export function usePermissions(
  permissionCodes: string[]
): Record<string, boolean> {
  const { hasPermission } = useSystemStore();
  return permissionCodes.reduce(
    (acc, code) => {
      acc[code] = hasPermission(code);
      return acc;
    },
    {} as Record<string, boolean>
  );
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
  return hasPermission ? React.createElement(React.Fragment, null, children) : React.createElement(React.Fragment, null, fallback);
}
