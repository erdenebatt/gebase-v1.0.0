"use client";

import { useState } from "react";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAuthStore } from "@/store/authStore";
import { useSystemStore } from "@/store/systemStore";
import { toast } from "sonner";
import {
  User,
  Mail,
  Phone,
  Building,
  Shield,
  Bell,
  Globe,
  Key,
  LogOut,
  Camera,
  Pencil,
  Check,
  X,
  Lock,
  Smartphone,
  Activity,
  Calendar,
  MapPin,
} from "lucide-react";
import { useRouter } from "next/navigation";

export default function ProfilePage() {
  const router = useRouter();
  const { user, clearSystemToken } = useAuthStore();
  const { currentSystem, currentRole, currentOrganization, clearSystemContext } = useSystemStore();

  const [isEditing, setIsEditing] = useState(false);
  const [profileData, setProfileData] = useState({
    firstName: user?.first_name || "",
    lastName: user?.last_name || "",
    familyName: user?.family_name || "",
    email: user?.email || "",
    phone: "",
  });

  const [passwordData, setPasswordData] = useState({
    currentPassword: "",
    newPassword: "",
    confirmPassword: "",
  });

  const [language, setLanguage] = useState(user?.language_code || "mn");

  const [notifications, setNotifications] = useState({
    emailNotifications: true,
    pushNotifications: true,
    loginAlerts: true,
  });

  const handleSaveProfile = () => {
    toast.success("Профайл амжилттай хадгалагдлаа");
    setIsEditing(false);
  };

  const handleChangePassword = () => {
    if (passwordData.newPassword !== passwordData.confirmPassword) {
      toast.error("Шинэ нууц үг таарахгүй байна");
      return;
    }
    if (passwordData.newPassword.length < 8) {
      toast.error("Нууц үг хамгийн багадаа 8 тэмдэгт байх ёстой");
      return;
    }
    toast.success("Нууц үг амжилттай солигдлоо");
    setPasswordData({ currentPassword: "", newPassword: "", confirmPassword: "" });
  };

  const handleSaveNotifications = () => {
    toast.success("Мэдэгдлийн тохиргоо хадгалагдлаа");
  };

  const handleLogout = async () => {
    try {
      clearSystemToken();
      clearSystemContext();
      useAuthStore.getState().logout();
      router.push("/login");
      toast.success("Амжилттай гарлаа");
    } catch {
      toast.error("Гарахад алдаа гарлаа");
    }
  };

  const securityItems = [
    {
      label: "Сүүлд нэвтэрсэн",
      value: "Өнөөдөр, 10:30",
      icon: Calendar,
      color: "text-blue-500",
      bgColor: "bg-blue-50 dark:bg-blue-500/10",
    },
    {
      label: "Идэвхтэй төхөөрөмж",
      value: "2 төхөөрөмж",
      icon: Smartphone,
      color: "text-emerald-500",
      bgColor: "bg-emerald-50 dark:bg-emerald-500/10",
    },
    {
      label: "2FA баталгаажуулалт",
      value: "Идэвхгүй",
      icon: Shield,
      color: "text-amber-500",
      bgColor: "bg-amber-50 dark:bg-amber-500/10",
    },
    {
      label: "Нууц үг солисон",
      value: "30 хоногийн өмнө",
      icon: Key,
      color: "text-violet-500",
      bgColor: "bg-violet-50 dark:bg-violet-500/10",
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <div className="h-12 w-12 rounded-2xl bg-gradient-to-br from-pink-500 to-pink-600 flex items-center justify-center shadow-lg shadow-pink-500/25">
          <User className="h-6 w-6 text-white" />
        </div>
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Миний профайл</h1>
          <p className="text-muted-foreground">
            Хувийн мэдээлэл болон тохиргоо
          </p>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Profile Card */}
        <Card className="lg:col-span-1 border-0 shadow-lg">
          <CardContent className="pt-6">
            <div className="flex flex-col items-center text-center">
              <div className="relative">
                <div className="h-28 w-28 rounded-2xl bg-gradient-to-br from-pink-500 to-violet-500 flex items-center justify-center shadow-xl">
                  {user?.avatar_url ? (
                    <img
                      src={user.avatar_url}
                      alt="Avatar"
                      className="h-28 w-28 rounded-2xl object-cover"
                    />
                  ) : (
                    <span className="text-4xl font-bold text-white">
                      {user?.first_name?.charAt(0) || "U"}
                    </span>
                  )}
                </div>
                <button className="absolute -bottom-2 -right-2 h-10 w-10 rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 text-white flex items-center justify-center hover:from-blue-600 hover:to-blue-700 shadow-lg transition-all">
                  <Camera className="h-5 w-5" />
                </button>
              </div>
              <h3 className="mt-4 text-xl font-bold">
                {user?.last_name} {user?.first_name}
              </h3>
              <p className="text-sm text-muted-foreground">{user?.email}</p>

              <div className="mt-6 w-full space-y-3">
                {currentSystem && (
                  <div className="flex items-center justify-between p-3 rounded-xl bg-gray-50 dark:bg-gray-800/50">
                    <span className="text-sm text-muted-foreground">Систем</span>
                    <span className="text-sm font-semibold">{currentSystem.name}</span>
                  </div>
                )}
                {currentRole && (
                  <div className="flex items-center justify-between p-3 rounded-xl bg-gray-50 dark:bg-gray-800/50">
                    <span className="text-sm text-muted-foreground">Эрх</span>
                    <span className="text-sm font-semibold">{currentRole.name}</span>
                  </div>
                )}
                {currentOrganization && (
                  <div className="flex items-center justify-between p-3 rounded-xl bg-gray-50 dark:bg-gray-800/50">
                    <span className="text-sm text-muted-foreground">Байгууллага</span>
                    <span className="text-sm font-semibold">{currentOrganization.name}</span>
                  </div>
                )}
              </div>

              <Button
                variant="destructive"
                className="mt-6 w-full bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700"
                onClick={handleLogout}
              >
                <LogOut className="mr-2 h-4 w-4" />
                Системээс гарах
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Personal Information */}
          <Card className="border-0 shadow-lg">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="h-10 w-10 rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
                    <User className="h-5 w-5 text-white" />
                  </div>
                  <div>
                    <CardTitle>Хувийн мэдээлэл</CardTitle>
                    <CardDescription>Таны үндсэн мэдээлэл</CardDescription>
                  </div>
                </div>
                <Button
                  variant={isEditing ? "default" : "outline"}
                  size="sm"
                  onClick={() => isEditing ? handleSaveProfile() : setIsEditing(true)}
                  className={isEditing ? "bg-gradient-to-r from-emerald-500 to-emerald-600" : ""}
                >
                  {isEditing ? (
                    <>
                      <Check className="h-4 w-4 mr-2" />
                      Хадгалах
                    </>
                  ) : (
                    <>
                      <Pencil className="h-4 w-4 mr-2" />
                      Засах
                    </>
                  )}
                </Button>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="lastName">Овог</Label>
                  <Input
                    id="lastName"
                    value={profileData.lastName}
                    onChange={(e) => setProfileData({ ...profileData, lastName: e.target.value })}
                    disabled={!isEditing}
                    className="bg-gray-50 dark:bg-gray-800 border-0"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="firstName">Нэр</Label>
                  <Input
                    id="firstName"
                    value={profileData.firstName}
                    onChange={(e) => setProfileData({ ...profileData, firstName: e.target.value })}
                    disabled={!isEditing}
                    className="bg-gray-50 dark:bg-gray-800 border-0"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email">Имэйл</Label>
                  <div className="relative">
                    <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                    <Input
                      id="email"
                      type="email"
                      value={profileData.email}
                      onChange={(e) => setProfileData({ ...profileData, email: e.target.value })}
                      disabled={!isEditing}
                      className="pl-10 bg-gray-50 dark:bg-gray-800 border-0"
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="phone">Утас</Label>
                  <div className="relative">
                    <Phone className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                    <Input
                      id="phone"
                      value={profileData.phone}
                      onChange={(e) => setProfileData({ ...profileData, phone: e.target.value })}
                      disabled={!isEditing}
                      className="pl-10 bg-gray-50 dark:bg-gray-800 border-0"
                      placeholder="99999999"
                    />
                  </div>
                </div>
              </div>
              {isEditing && (
                <div className="flex gap-2 pt-2">
                  <Button onClick={handleSaveProfile} className="bg-gradient-to-r from-emerald-500 to-emerald-600">
                    Хадгалах
                  </Button>
                  <Button variant="outline" onClick={() => setIsEditing(false)}>
                    Цуцлах
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Change Password */}
          <Card className="border-0 shadow-lg">
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-xl bg-gradient-to-br from-amber-500 to-amber-600 flex items-center justify-center">
                  <Key className="h-5 w-5 text-white" />
                </div>
                <div>
                  <CardTitle>Нууц үг солих</CardTitle>
                  <CardDescription>Аюулгүй байдлын үүднээс нууц үгээ тогтмол солиорой</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 md:grid-cols-3">
                <div className="space-y-2">
                  <Label htmlFor="currentPassword">Одоогийн нууц үг</Label>
                  <Input
                    id="currentPassword"
                    type="password"
                    value={passwordData.currentPassword}
                    onChange={(e) => setPasswordData({ ...passwordData, currentPassword: e.target.value })}
                    placeholder="••••••••"
                    className="bg-gray-50 dark:bg-gray-800 border-0"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="newPassword">Шинэ нууц үг</Label>
                  <Input
                    id="newPassword"
                    type="password"
                    value={passwordData.newPassword}
                    onChange={(e) => setPasswordData({ ...passwordData, newPassword: e.target.value })}
                    placeholder="••••••••"
                    className="bg-gray-50 dark:bg-gray-800 border-0"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="confirmPassword">Нууц үг давтах</Label>
                  <Input
                    id="confirmPassword"
                    type="password"
                    value={passwordData.confirmPassword}
                    onChange={(e) => setPasswordData({ ...passwordData, confirmPassword: e.target.value })}
                    placeholder="••••••••"
                    className="bg-gray-50 dark:bg-gray-800 border-0"
                  />
                </div>
              </div>
              <Button onClick={handleChangePassword} className="bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-600 hover:to-amber-700">
                Нууц үг солих
              </Button>
            </CardContent>
          </Card>

          {/* Security Info */}
          <Card className="border-0 shadow-lg">
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-xl bg-gradient-to-br from-violet-500 to-violet-600 flex items-center justify-center">
                  <Shield className="h-5 w-5 text-white" />
                </div>
                <div>
                  <CardTitle>Аюулгүй байдал</CardTitle>
                  <CardDescription>Таны бүртгэлийн аюулгүй байдлын мэдээлэл</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-2">
                {securityItems.map((item, index) => (
                  <div key={index} className="flex items-center gap-4 p-4 rounded-xl bg-gray-50 dark:bg-gray-800/50">
                    <div className={`h-10 w-10 rounded-xl ${item.bgColor} flex items-center justify-center`}>
                      <item.icon className={`h-5 w-5 ${item.color}`} />
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">{item.label}</p>
                      <p className="font-semibold">{item.value}</p>
                    </div>
                  </div>
                ))}
              </div>
              <div className="mt-4">
                <Button variant="outline" className="border-violet-200 text-violet-600 hover:bg-violet-50">
                  <Lock className="h-4 w-4 mr-2" />
                  2 шатлалт баталгаажуулалт идэвхжүүлэх
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
