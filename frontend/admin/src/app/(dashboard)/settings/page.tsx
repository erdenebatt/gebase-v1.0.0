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
import { toast } from "sonner";
import {
  Settings,
  Bell,
  Palette,
  Globe,
  Database,
  Shield,
  Mail,
  Clock,
  Moon,
  Sun,
  Monitor,
  Server,
  Check,
  RefreshCw,
  Send,
} from "lucide-react";
import { useTheme } from "next-themes";

export default function SettingsPage() {
  const { theme, setTheme } = useTheme();

  const [generalSettings, setGeneralSettings] = useState({
    siteName: "Gebase Platform",
    siteDescription: "Enterprise Auth & RBAC Platform",
    adminEmail: "admin@gerege.mn",
    supportEmail: "support@gerege.mn",
  });

  const [securitySettings, setSecuritySettings] = useState({
    sessionTimeout: "24",
    maxLoginAttempts: "5",
    passwordMinLength: "8",
    requireTwoFactor: false,
  });

  const [emailSettings, setEmailSettings] = useState({
    smtpHost: "smtp.gerege.mn",
    smtpPort: "587",
    smtpUser: "",
    smtpPassword: "",
    fromEmail: "noreply@gerege.mn",
    fromName: "Gebase Platform",
  });

  const [notificationSettings, setNotificationSettings] = useState({
    enableEmailNotifications: true,
    enablePushNotifications: true,
    notifyOnLogin: true,
    notifyOnPasswordChange: true,
    notifyOnRoleChange: false,
  });

  const handleSaveGeneral = () => {
    toast.success("Ерөнхий тохиргоо хадгалагдлаа");
  };

  const handleSaveSecurity = () => {
    toast.success("Аюулгүй байдлын тохиргоо хадгалагдлаа");
  };

  const handleSaveEmail = () => {
    toast.success("Имэйл тохиргоо хадгалагдлаа");
  };

  const handleSaveNotifications = () => {
    toast.success("Мэдэгдлийн тохиргоо хадгалагдлаа");
  };

  const stats = [
    {
      title: "Системүүд",
      value: "2",
      icon: Server,
      color: "text-blue-500",
      bgColor: "bg-blue-100 dark:bg-blue-900/30",
    },
    {
      title: "Хэл",
      value: "2",
      icon: Globe,
      color: "text-green-500",
      bgColor: "bg-green-100 dark:bg-green-900/30",
    },
    {
      title: "Горим",
      value: theme === "dark" ? "Бараан" : "Цайвар",
      icon: Palette,
      color: "text-purple-500",
      bgColor: "bg-purple-100 dark:bg-purple-900/30",
    },
    {
      title: "Холболт",
      value: "Идэвхтэй",
      icon: Database,
      color: "text-emerald-500",
      bgColor: "bg-emerald-100 dark:bg-emerald-900/30",
    },
  ];

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2">
            <Settings className="h-6 w-6 text-primary" />
            Системийн тохиргоо
          </h1>
          <p className="text-muted-foreground">
            Платформын ерөнхий тохиргоо болон параметрүүд
          </p>
        </div>
        <Button variant="outline" className="gap-2">
          <RefreshCw className="h-4 w-4" />
          Шинэчлэх
        </Button>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        {stats.map((stat, index) => (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {stat.title}
              </CardTitle>
              <div className={`p-2 rounded-lg ${stat.bgColor}`}>
                <stat.icon className={`h-4 w-4 ${stat.color}`} />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* General Settings */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <Settings className="h-5 w-5" />
              </div>
              <div>
                <CardTitle>Ерөнхий тохиргоо</CardTitle>
                <CardDescription>
                  Платформын үндсэн мэдээлэл
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4">
              <div className="space-y-2">
                <Label htmlFor="siteName">Платформын нэр</Label>
                <Input
                  id="siteName"
                  value={generalSettings.siteName}
                  onChange={(e) =>
                    setGeneralSettings({ ...generalSettings, siteName: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="siteDescription">Тайлбар</Label>
                <Input
                  id="siteDescription"
                  value={generalSettings.siteDescription}
                  onChange={(e) =>
                    setGeneralSettings({ ...generalSettings, siteDescription: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="adminEmail">Админ имэйл</Label>
                <Input
                  id="adminEmail"
                  type="email"
                  value={generalSettings.adminEmail}
                  onChange={(e) =>
                    setGeneralSettings({ ...generalSettings, adminEmail: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="supportEmail">Дэмжлэгийн имэйл</Label>
                <Input
                  id="supportEmail"
                  type="email"
                  value={generalSettings.supportEmail}
                  onChange={(e) =>
                    setGeneralSettings({ ...generalSettings, supportEmail: e.target.value })
                  }
                />
              </div>
            </div>
            <Button onClick={handleSaveGeneral} className="gap-2">
              <Check className="h-4 w-4" />
              Хадгалах
            </Button>
          </CardContent>
        </Card>

        {/* Theme Settings */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-purple-100 dark:bg-purple-900/30 text-purple-500">
                <Palette className="h-5 w-5" />
              </div>
              <div>
                <CardTitle>Харагдах байдал</CardTitle>
                <CardDescription>
                  Системийн өнгө болон загварыг тохируулах
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="space-y-3">
              <Label>Горим сонгох</Label>
              <div className="grid grid-cols-3 gap-3">
                <button
                  onClick={() => setTheme("light")}
                  className={`flex flex-col items-center gap-2 p-4 rounded-lg border-2 transition-all ${
                    theme === "light"
                      ? "border-primary bg-primary/5"
                      : "border-transparent bg-muted hover:border-muted-foreground/20"
                  }`}
                >
                  <Sun className={`h-6 w-6 ${theme === "light" ? "text-primary" : "text-muted-foreground"}`} />
                  <span className={`text-sm font-medium ${theme === "light" ? "text-primary" : ""}`}>
                    Цайвар
                  </span>
                </button>
                <button
                  onClick={() => setTheme("dark")}
                  className={`flex flex-col items-center gap-2 p-4 rounded-lg border-2 transition-all ${
                    theme === "dark"
                      ? "border-primary bg-primary/5"
                      : "border-transparent bg-muted hover:border-muted-foreground/20"
                  }`}
                >
                  <Moon className={`h-6 w-6 ${theme === "dark" ? "text-primary" : "text-muted-foreground"}`} />
                  <span className={`text-sm font-medium ${theme === "dark" ? "text-primary" : ""}`}>
                    Бараан
                  </span>
                </button>
                <button
                  onClick={() => setTheme("system")}
                  className={`flex flex-col items-center gap-2 p-4 rounded-lg border-2 transition-all ${
                    theme === "system"
                      ? "border-primary bg-primary/5"
                      : "border-transparent bg-muted hover:border-muted-foreground/20"
                  }`}
                >
                  <Monitor className={`h-6 w-6 ${theme === "system" ? "text-primary" : "text-muted-foreground"}`} />
                  <span className={`text-sm font-medium ${theme === "system" ? "text-primary" : ""}`}>
                    Систем
                  </span>
                </button>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Security Settings */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-red-100 dark:bg-red-900/30 text-red-500">
                <Shield className="h-5 w-5" />
              </div>
              <div>
                <CardTitle>Аюулгүй байдал</CardTitle>
                <CardDescription>
                  Нэвтрэлт болон аюулгүй байдлын тохиргоо
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4">
              <div className="space-y-2">
                <Label htmlFor="sessionTimeout">Session хугацаа (цаг)</Label>
                <Select
                  value={securitySettings.sessionTimeout}
                  onValueChange={(value) =>
                    setSecuritySettings({ ...securitySettings, sessionTimeout: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="1">1 цаг</SelectItem>
                    <SelectItem value="8">8 цаг</SelectItem>
                    <SelectItem value="24">24 цаг</SelectItem>
                    <SelectItem value="168">7 хоног</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="maxLoginAttempts">Нэвтрэх оролдлогын хязгаар</Label>
                <Select
                  value={securitySettings.maxLoginAttempts}
                  onValueChange={(value) =>
                    setSecuritySettings({ ...securitySettings, maxLoginAttempts: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="3">3 удаа</SelectItem>
                    <SelectItem value="5">5 удаа</SelectItem>
                    <SelectItem value="10">10 удаа</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="passwordMinLength">Нууц үгийн доод урт</Label>
                <Select
                  value={securitySettings.passwordMinLength}
                  onValueChange={(value) =>
                    setSecuritySettings({ ...securitySettings, passwordMinLength: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="6">6 тэмдэгт</SelectItem>
                    <SelectItem value="8">8 тэмдэгт</SelectItem>
                    <SelectItem value="12">12 тэмдэгт</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="flex items-center space-x-3 pt-2">
                <input
                  type="checkbox"
                  id="requireTwoFactor"
                  checked={securitySettings.requireTwoFactor}
                  onChange={(e) =>
                    setSecuritySettings({ ...securitySettings, requireTwoFactor: e.target.checked })
                  }
                  className="h-4 w-4 rounded border-gray-300"
                />
                <Label htmlFor="requireTwoFactor" className="cursor-pointer">
                  2 шатлалт баталгаажуулалт шаардах
                </Label>
              </div>
            </div>
            <Button onClick={handleSaveSecurity} className="gap-2">
              <Check className="h-4 w-4" />
              Хадгалах
            </Button>
          </CardContent>
        </Card>

        {/* Email Settings */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-100 dark:bg-blue-900/30 text-blue-500">
                <Mail className="h-5 w-5" />
              </div>
              <div>
                <CardTitle>Имэйл тохиргоо</CardTitle>
                <CardDescription>
                  SMTP серверийн тохиргоо
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="smtpHost">SMTP сервер</Label>
                <Input
                  id="smtpHost"
                  value={emailSettings.smtpHost}
                  onChange={(e) =>
                    setEmailSettings({ ...emailSettings, smtpHost: e.target.value })
                  }
                  placeholder="smtp.example.com"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="smtpPort">Порт</Label>
                <Input
                  id="smtpPort"
                  value={emailSettings.smtpPort}
                  onChange={(e) =>
                    setEmailSettings({ ...emailSettings, smtpPort: e.target.value })
                  }
                  placeholder="587"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="smtpUser">Хэрэглэгчийн нэр</Label>
                <Input
                  id="smtpUser"
                  value={emailSettings.smtpUser}
                  onChange={(e) =>
                    setEmailSettings({ ...emailSettings, smtpUser: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="smtpPassword">Нууц үг</Label>
                <Input
                  id="smtpPassword"
                  type="password"
                  value={emailSettings.smtpPassword}
                  onChange={(e) =>
                    setEmailSettings({ ...emailSettings, smtpPassword: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="fromEmail">Илгээгчийн имэйл</Label>
                <Input
                  id="fromEmail"
                  type="email"
                  value={emailSettings.fromEmail}
                  onChange={(e) =>
                    setEmailSettings({ ...emailSettings, fromEmail: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="fromName">Илгээгчийн нэр</Label>
                <Input
                  id="fromName"
                  value={emailSettings.fromName}
                  onChange={(e) =>
                    setEmailSettings({ ...emailSettings, fromName: e.target.value })
                  }
                />
              </div>
            </div>
            <div className="flex gap-2">
              <Button onClick={handleSaveEmail} className="gap-2">
                <Check className="h-4 w-4" />
                Хадгалах
              </Button>
              <Button variant="outline" className="gap-2">
                <Send className="h-4 w-4" />
                Тест илгээх
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Notification Settings - Full Width */}
      <Card>
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-amber-100 dark:bg-amber-900/30 text-amber-500">
              <Bell className="h-5 w-5" />
            </div>
            <div>
              <CardTitle>Мэдэгдлийн тохиргоо</CardTitle>
              <CardDescription>
                Системийн мэдэгдлийн тохиргоо
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            <div className="flex items-center justify-between p-4 rounded-lg border">
              <div>
                <p className="font-medium">Имэйл мэдэгдэл</p>
                <p className="text-sm text-muted-foreground">
                  Имэйлээр мэдэгдэл илгээх
                </p>
              </div>
              <input
                type="checkbox"
                checked={notificationSettings.enableEmailNotifications}
                onChange={(e) =>
                  setNotificationSettings({
                    ...notificationSettings,
                    enableEmailNotifications: e.target.checked,
                  })
                }
                className="h-5 w-5 rounded border-gray-300"
              />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg border">
              <div>
                <p className="font-medium">Push мэдэгдэл</p>
                <p className="text-sm text-muted-foreground">
                  Браузер мэдэгдэл илгээх
                </p>
              </div>
              <input
                type="checkbox"
                checked={notificationSettings.enablePushNotifications}
                onChange={(e) =>
                  setNotificationSettings({
                    ...notificationSettings,
                    enablePushNotifications: e.target.checked,
                  })
                }
                className="h-5 w-5 rounded border-gray-300"
              />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg border">
              <div>
                <p className="font-medium">Нэвтрэлтийн мэдэгдэл</p>
                <p className="text-sm text-muted-foreground">
                  Шинэ нэвтрэлт хийгдэхэд
                </p>
              </div>
              <input
                type="checkbox"
                checked={notificationSettings.notifyOnLogin}
                onChange={(e) =>
                  setNotificationSettings({
                    ...notificationSettings,
                    notifyOnLogin: e.target.checked,
                  })
                }
                className="h-5 w-5 rounded border-gray-300"
              />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg border">
              <div>
                <p className="font-medium">Нууц үг солигдсон</p>
                <p className="text-sm text-muted-foreground">
                  Нууц үг солигдоход мэдэгдэх
                </p>
              </div>
              <input
                type="checkbox"
                checked={notificationSettings.notifyOnPasswordChange}
                onChange={(e) =>
                  setNotificationSettings({
                    ...notificationSettings,
                    notifyOnPasswordChange: e.target.checked,
                  })
                }
                className="h-5 w-5 rounded border-gray-300"
              />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg border">
              <div>
                <p className="font-medium">Эрх өөрчлөгдсөн</p>
                <p className="text-sm text-muted-foreground">
                  Хэрэглэгчийн эрх өөрчлөгдөхөд
                </p>
              </div>
              <input
                type="checkbox"
                checked={notificationSettings.notifyOnRoleChange}
                onChange={(e) =>
                  setNotificationSettings({
                    ...notificationSettings,
                    notifyOnRoleChange: e.target.checked,
                  })
                }
                className="h-5 w-5 rounded border-gray-300"
              />
            </div>
          </div>
          <Button onClick={handleSaveNotifications} className="gap-2">
            <Check className="h-4 w-4" />
            Хадгалах
          </Button>
        </CardContent>
      </Card>

      {/* System Info */}
      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-emerald-100 dark:bg-emerald-900/30 text-emerald-500">
                <Database className="h-5 w-5" />
              </div>
              <div>
                <CardTitle>Өгөгдлийн сан</CardTitle>
                <CardDescription>
                  Өгөгдлийн сангийн мэдээлэл
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid gap-3">
              <div className="flex items-center justify-between p-3 rounded-lg bg-muted">
                <span className="text-sm text-muted-foreground">Төрөл</span>
                <span className="font-semibold">PostgreSQL 15</span>
              </div>
              <div className="flex items-center justify-between p-3 rounded-lg bg-muted">
                <span className="text-sm text-muted-foreground">Схем</span>
                <span className="font-semibold">gerege_base</span>
              </div>
              <div className="flex items-center justify-between p-3 rounded-lg bg-muted">
                <span className="text-sm text-muted-foreground">Холболт</span>
                <span className="font-semibold text-green-600">Идэвхтэй</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-cyan-100 dark:bg-cyan-900/30 text-cyan-500">
                <Clock className="h-5 w-5" />
              </div>
              <div>
                <CardTitle>Системийн мэдээлэл</CardTitle>
                <CardDescription>
                  Техникийн мэдээлэл
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid gap-3">
              <div className="flex items-center justify-between p-3 rounded-lg bg-muted">
                <span className="text-sm text-muted-foreground">Хувилбар</span>
                <span className="font-semibold">v1.0.0</span>
              </div>
              <div className="flex items-center justify-between p-3 rounded-lg bg-muted">
                <span className="text-sm text-muted-foreground">Backend</span>
                <span className="font-semibold">Go 1.22</span>
              </div>
              <div className="flex items-center justify-between p-3 rounded-lg bg-muted">
                <span className="text-sm text-muted-foreground">Frontend</span>
                <span className="font-semibold">Next.js 14</span>
              </div>
              <div className="flex items-center justify-between p-3 rounded-lg bg-muted">
                <span className="text-sm text-muted-foreground">Цаг бүс</span>
                <span className="font-semibold">Asia/Ulaanbaatar</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
