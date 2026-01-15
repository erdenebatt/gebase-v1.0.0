"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useAuthStore } from "@/store/authStore";
import { useSystemStore } from "@/store/systemStore";
import { authApi, devicesApi } from "@/lib/api";
import { toast } from "sonner";

const loginSchema = z.object({
  email: z.string().email("И-мэйл хаяг буруу байна"),
  password: z.string().min(1, "Нууц үг оруулна уу"),
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
      let deviceUID = localStorage.getItem("device_uid");
      if (!deviceUID) {
        deviceUID = `web_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        localStorage.setItem("device_uid", deviceUID);
      }

      // Register device first (creates or updates)
      await devicesApi.register({
        device_uid: deviceUID,
        name: navigator.userAgent.substring(0, 50),
        platform: "web",
        os_version: navigator.platform,
        app_version: "1.0.0",
      });

      const response = await authApi.login({
        email: data.email,
        password: data.password,
        device_uid: deviceUID,
      });

      const result = response.data.data;

      // Transform flat systems array to SystemRole format expected by frontend
      // Backend returns: { id, code, name, icon_name, ... }
      // Frontend expects: { system: {...}, roles: [] }
      const transformedSystems = (result.available_systems || []).map((sys: { id: number; code: string; name: string; description?: string; icon_name?: string; icon_url?: string; color?: string; sequence?: number }) => ({
        system: {
          id: sys.id,
          code: sys.code,
          name: sys.name,
          description: sys.description || "",
          icon_name: sys.icon_name || "Box",
          icon_url: sys.icon_url || "",
          color: sys.color || "#6366f1",
          sequence: sys.sequence || 0,
        },
        roles: [], // Roles will be populated when switching systems
      }));

      // Set auth state (backend returns access_token, not platform_token)
      setAuth({
        user: result.user,
        platformToken: result.access_token,
        refreshToken: result.refresh_token,
        availableSystems: transformedSystems,
      });

      toast.success("Амжилттай нэвтэрлээ");

      // Automatically switch to the first available system (typically "admin")
      if (transformedSystems.length > 0) {
        const firstSystem = transformedSystems[0].system;
        try {
          // Use the platform token directly instead of relying on interceptor
          const switchResponse = await fetch(
            `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000"}/api/v1/auth/switch-system`,
            {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${result.access_token}`,
                "X-Device-UID": deviceUID,
                "X-Platform": "web",
              },
              body: JSON.stringify({ system_code: firstSystem.code }),
            }
          );

          const switchResult = await switchResponse.json();

          if (switchResult.success && switchResult.data) {
            const switchData = switchResult.data;

            // Update system token and context
            useAuthStore.getState().setSystemToken(switchData.system_token);
            useSystemStore.getState().setSystemContext({
              currentSystem: switchData.current_system,
              currentRole: switchData.current_role,
              currentOrganization: switchData.current_organization,
              permissions: switchData.permissions || [],
              menus: switchData.menus || [],
            });
          }
        } catch (switchErr) {
          console.error("Failed to switch to system:", switchErr);
        }
      }

      router.push("/dashboard");
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: { message?: string } } } };
      const message =
        err.response?.data?.error?.message || "Нэвтрэхэд алдаа гарлаа";
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
                {...register("email")}
              />
              {errors.email && (
                <p className="text-sm text-destructive">
                  {errors.email.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="password">Нууц үг</Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                {...register("password")}
              />
              {errors.password && (
                <p className="text-sm text-destructive">
                  {errors.password.message}
                </p>
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
