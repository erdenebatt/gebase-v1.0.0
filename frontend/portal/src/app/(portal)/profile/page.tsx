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
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useAuthStore } from "@/store/authStore";
import { portalApi } from "@/lib/api";
import { toast } from "sonner";
import { Loader2, Camera } from "lucide-react";

const profileSchema = z.object({
  first_name: z.string().min(1, "Нэр оруулна уу"),
  last_name: z.string().min(1, "Овог оруулна уу"),
  email: z.string().email("И-мэйл хаяг буруу байна"),
  phone_no: z.string().optional(),
});

type ProfileForm = z.infer<typeof profileSchema>;

export default function ProfilePage() {
  const { user } = useAuthStore();
  const [isLoading, setIsLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ProfileForm>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      first_name: user?.first_name || "",
      last_name: user?.last_name || "",
      email: user?.email || "",
    },
  });

  const getInitials = () => {
    const f = user?.first_name?.[0] || "";
    const l = user?.last_name?.[0] || "";
    return (f + l).toUpperCase() || "U";
  };

  const onSubmit = async (data: ProfileForm) => {
    setIsLoading(true);
    try {
      await portalApi.profile.update(data);
      toast.success("Профайл амжилттай шинэчлэгдлээ");
    } catch (error) {
      toast.error("Алдаа гарлаа");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container py-8 max-w-2xl">
      <div className="mb-8">
        <h1 className="text-2xl font-bold tracking-tight">Профайл</h1>
        <p className="text-muted-foreground">
          Хувийн мэдээллээ удирдах
        </p>
      </div>

      {/* Avatar */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Профайл зураг</CardTitle>
          <CardDescription>
            Таны зураг бусад хэрэглэгчдэд харагдана
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-6">
            <div className="relative">
              <Avatar className="h-24 w-24">
                <AvatarImage src={user?.avatar_url} />
                <AvatarFallback className="text-2xl bg-primary/10 text-primary">
                  {getInitials()}
                </AvatarFallback>
              </Avatar>
              <Button
                size="icon"
                variant="secondary"
                className="absolute bottom-0 right-0 h-8 w-8 rounded-full"
              >
                <Camera className="h-4 w-4" />
              </Button>
            </div>
            <div>
              <p className="font-medium">
                {user?.first_name} {user?.last_name}
              </p>
              <p className="text-sm text-muted-foreground">{user?.email}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Profile Form */}
      <Card>
        <CardHeader>
          <CardTitle>Хувийн мэдээлэл</CardTitle>
          <CardDescription>
            Хувийн мэдээллээ шинэчлэх
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="last_name">Овог</Label>
                <Input id="last_name" {...register("last_name")} />
                {errors.last_name && (
                  <p className="text-sm text-destructive">
                    {errors.last_name.message}
                  </p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="first_name">Нэр</Label>
                <Input id="first_name" {...register("first_name")} />
                {errors.first_name && (
                  <p className="text-sm text-destructive">
                    {errors.first_name.message}
                  </p>
                )}
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="email">И-мэйл</Label>
              <Input id="email" type="email" {...register("email")} />
              {errors.email && (
                <p className="text-sm text-destructive">
                  {errors.email.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="phone_no">Утасны дугаар</Label>
              <Input
                id="phone_no"
                type="tel"
                placeholder="99999999"
                {...register("phone_no")}
              />
            </div>

            <div className="flex justify-end">
              <Button type="submit" disabled={isLoading}>
                {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Хадгалах
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
