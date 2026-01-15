"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { portalApi } from "@/lib/api";
import { toast } from "sonner";
import { Loader2, Shield, Key, Bell, Globe } from "lucide-react";

const passwordSchema = z
  .object({
    old_password: z.string().min(1, "Хуучин нууц үгээ оруулна уу"),
    new_password: z.string().min(8, "Нууц үг хамгийн багадаа 8 тэмдэгт байна"),
    confirm_password: z.string(),
  })
  .refine((data) => data.new_password === data.confirm_password, {
    message: "Нууц үг таарахгүй байна",
    path: ["confirm_password"],
  });

type PasswordForm = z.infer<typeof passwordSchema>;

export default function SettingsPage() {
  const [isLoading, setIsLoading] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<PasswordForm>({
    resolver: zodResolver(passwordSchema),
  });

  const onSubmit = async (data: PasswordForm) => {
    setIsLoading(true);
    try {
      await portalApi.profile.changePassword({
        old_password: data.old_password,
        new_password: data.new_password,
      });
      toast.success("Нууц үг амжилттай солигдлоо");
      reset();
    } catch (error) {
      toast.error("Алдаа гарлаа");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container py-8 max-w-2xl">
      <div className="mb-8">
        <h1 className="text-2xl font-bold tracking-tight">Тохиргоо</h1>
        <p className="text-muted-foreground">Аккаунтын тохиргоо</p>
      </div>

      {/* Security Settings */}
      <Card className="mb-6">
        <CardHeader>
          <div className="flex items-center gap-2">
            <Shield className="h-5 w-5 text-primary" />
            <CardTitle>Нууцлал</CardTitle>
          </div>
          <CardDescription>
            Нууц үг болон нууцлалын тохиргоо
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="old_password">Хуучин нууц үг</Label>
              <Input
                id="old_password"
                type="password"
                {...register("old_password")}
              />
              {errors.old_password && (
                <p className="text-sm text-destructive">
                  {errors.old_password.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="new_password">Шинэ нууц үг</Label>
              <Input
                id="new_password"
                type="password"
                {...register("new_password")}
              />
              {errors.new_password && (
                <p className="text-sm text-destructive">
                  {errors.new_password.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="confirm_password">Шинэ нууц үг давтах</Label>
              <Input
                id="confirm_password"
                type="password"
                {...register("confirm_password")}
              />
              {errors.confirm_password && (
                <p className="text-sm text-destructive">
                  {errors.confirm_password.message}
                </p>
              )}
            </div>

            <div className="flex justify-end">
              <Button type="submit" disabled={isLoading}>
                {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Нууц үг солих
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {/* Notification Settings */}
      <Card className="mb-6">
        <CardHeader>
          <div className="flex items-center gap-2">
            <Bell className="h-5 w-5 text-primary" />
            <CardTitle>Мэдэгдэл</CardTitle>
          </div>
          <CardDescription>
            Мэдэгдлийн тохиргоо
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium">И-мэйл мэдэгдэл</p>
                <p className="text-sm text-muted-foreground">
                  Чухал мэдээллүүдийг и-мэйлээр авах
                </p>
              </div>
              <Button variant="outline" size="sm">
                Идэвхжүүлсэн
              </Button>
            </div>
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium">Push мэдэгдэл</p>
                <p className="text-sm text-muted-foreground">
                  Браузер мэдэгдэл авах
                </p>
              </div>
              <Button variant="outline" size="sm">
                Идэвхгүй
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Language Settings */}
      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <Globe className="h-5 w-5 text-primary" />
            <CardTitle>Хэл</CardTitle>
          </div>
          <CardDescription>Интерфэйсийн хэл сонгох</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">Хэл</p>
              <p className="text-sm text-muted-foreground">
                Одоо: Монгол
              </p>
            </div>
            <Button variant="outline" size="sm">
              Солих
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
