"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
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
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useAuthStore } from "@/store/authStore";
import { useSystemStore } from "@/store/systemStore";
import { authApi } from "@/lib/api";
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

      // Set platform context if provided
      if (result.permissions && result.menus) {
        setPlatformContext({
          permissions: result.permissions,
          menus: result.menus,
        });
      }

      toast.success("Амжилттай нэвтэрлээ");
      router.push("/dashboard");
    } catch (error: unknown) {
      const err = error as {
        response?: { data?: { error?: { message?: string } } };
      };
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
          <Link href="/" className="inline-block mx-auto mb-4">
            <div className="h-12 w-12 rounded-xl bg-primary flex items-center justify-center">
              <span className="text-2xl font-bold text-white">G</span>
            </div>
          </Link>
          <CardTitle className="text-2xl">Нэвтрэх</CardTitle>
          <CardDescription>
            Gebase Portal-д нэвтрэхийн тулд мэдээллээ оруулна уу
          </CardDescription>
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
              <div className="flex items-center justify-between">
                <Label htmlFor="password">Нууц үг</Label>
                <Link
                  href="/forgot-password"
                  className="text-sm text-primary hover:underline"
                >
                  Нууц үгээ мартсан?
                </Link>
              </div>
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
        <CardFooter className="justify-center">
          <p className="text-sm text-muted-foreground">
            Бүртгэлгүй юу?{" "}
            <Link href="/register" className="text-primary hover:underline">
              Бүртгүүлэх
            </Link>
          </p>
        </CardFooter>
      </Card>
    </div>
  );
}
