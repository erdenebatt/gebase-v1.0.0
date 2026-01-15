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
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAuthStore } from "@/store/authStore";
import { toast } from "sonner";
import {
  Bell,
  Globe,
  Palette,
  Monitor,
  Smartphone,
  Laptop,
  Moon,
  Sun,
  Eye,
  EyeOff,
  Trash2,
} from "lucide-react";
import { useTheme } from "next-themes";

const mockDevices = [
  {
    id: 1,
    name: "Chrome - Windows",
    type: "desktop",
    lastActive: "Идэвхтэй",
    current: true,
  },
  {
    id: 2,
    name: "Safari - iPhone",
    type: "mobile",
    lastActive: "2 цагийн өмнө",
    current: false,
  },
  {
    id: 3,
    name: "Firefox - MacOS",
    type: "desktop",
    lastActive: "1 өдрийн өмнө",
    current: false,
  },
];

export default function UserSettingsPage() {
  const { theme, setTheme } = useTheme();
  const { user } = useAuthStore();

  const [language, setLanguage] = useState(user?.language_code || "mn");
  const [timezone, setTimezone] = useState("Asia/Ulaanbaatar");

  const [notifications, setNotifications] = useState({
    emailNotifications: true,
    pushNotifications: true,
    loginAlerts: true,
    weeklyReport: false,
    marketingEmails: false,
  });

  const [privacy, setPrivacy] = useState({
    showOnlineStatus: true,
    showLastSeen: true,
    showEmail: false,
  });

  const handleSaveLanguage = () => {
    toast.success("Хэлний тохиргоо хадгалагдлаа");
  };

  const handleSaveNotifications = () => {
    toast.success("Мэдэгдлийн тохиргоо хадгалагдлаа");
  };

  const handleSavePrivacy = () => {
    toast.success("Нууцлалын тохиргоо хадгалагдлаа");
  };

  const handleRemoveDevice = (deviceId: number) => {
    toast.success("Төхөөрөмж амжилттай салгагдлаа");
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Хэрэглэгчийн тохиргоо</h1>
        <p className="text-muted-foreground">
          Таны хувийн тохиргоо болон сонголтууд
        </p>
      </div>

      <div className="grid gap-6">
        {/* Theme Settings */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Palette className="h-5 w-5" />
              Харагдах байдал
            </CardTitle>
            <CardDescription>
              Интерфэйсийн өнгө, загварыг өөрийн хүсэлтээр тохируулах
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label>Өнгөний горим</Label>
              <div className="flex gap-2">
                <Button
                  variant={theme === "light" ? "default" : "outline"}
                  size="sm"
                  onClick={() => setTheme("light")}
                  className="flex items-center gap-2"
                >
                  <Sun className="h-4 w-4" />
                  Цайвар
                </Button>
                <Button
                  variant={theme === "dark" ? "default" : "outline"}
                  size="sm"
                  onClick={() => setTheme("dark")}
                  className="flex items-center gap-2"
                >
                  <Moon className="h-4 w-4" />
                  Бараан
                </Button>
                <Button
                  variant={theme === "system" ? "default" : "outline"}
                  size="sm"
                  onClick={() => setTheme("system")}
                  className="flex items-center gap-2"
                >
                  <Monitor className="h-4 w-4" />
                  Автомат
                </Button>
              </div>
              <p className="text-sm text-muted-foreground">
                Автомат горим нь таны төхөөрөмжийн тохиргоог дагана
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Language & Region */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Globe className="h-5 w-5" />
              Хэл & Бүс
            </CardTitle>
            <CardDescription>
              Хэл болон цагийн бүсийн тохиргоо
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="language">Хэл</Label>
                <Select value={language} onValueChange={setLanguage}>
                  <SelectTrigger>
                    <SelectValue placeholder="Хэл сонгох" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="mn">🇲🇳 Монгол</SelectItem>
                    <SelectItem value="en">🇺🇸 English</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="timezone">Цагийн бүс</Label>
                <Select value={timezone} onValueChange={setTimezone}>
                  <SelectTrigger>
                    <SelectValue placeholder="Цагийн бүс сонгох" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="Asia/Ulaanbaatar">
                      Улаанбаатар (UTC+8)
                    </SelectItem>
                    <SelectItem value="Asia/Hovd">Ховд (UTC+7)</SelectItem>
                    <SelectItem value="Asia/Choibalsan">
                      Чойбалсан (UTC+8)
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <Button onClick={handleSaveLanguage}>Хадгалах</Button>
          </CardContent>
        </Card>

        {/* Notification Settings */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Bell className="h-5 w-5" />
              Мэдэгдэл
            </CardTitle>
            <CardDescription>
              Мэдэгдэл хүлээн авах төрлүүдийг тохируулах
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Имэйл мэдэгдэл</p>
                  <p className="text-sm text-muted-foreground">
                    Чухал мэдэгдлүүдийг имэйлээр авах
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={notifications.emailNotifications}
                  onChange={(e) =>
                    setNotifications({
                      ...notifications,
                      emailNotifications: e.target.checked,
                    })
                  }
                  className="h-4 w-4 rounded border-gray-300"
                />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Push мэдэгдэл</p>
                  <p className="text-sm text-muted-foreground">
                    Браузер дээр шууд мэдэгдэл авах
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={notifications.pushNotifications}
                  onChange={(e) =>
                    setNotifications({
                      ...notifications,
                      pushNotifications: e.target.checked,
                    })
                  }
                  className="h-4 w-4 rounded border-gray-300"
                />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Нэвтрэлтийн сэрэмжлүүлэг</p>
                  <p className="text-sm text-muted-foreground">
                    Шинэ төхөөрөмжөөс нэвтрэхэд мэдэгдэл авах
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={notifications.loginAlerts}
                  onChange={(e) =>
                    setNotifications({
                      ...notifications,
                      loginAlerts: e.target.checked,
                    })
                  }
                  className="h-4 w-4 rounded border-gray-300"
                />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Долоо хоногийн тайлан</p>
                  <p className="text-sm text-muted-foreground">
                    7 хоног тутам үйл ажиллагааны тайлан авах
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={notifications.weeklyReport}
                  onChange={(e) =>
                    setNotifications({
                      ...notifications,
                      weeklyReport: e.target.checked,
                    })
                  }
                  className="h-4 w-4 rounded border-gray-300"
                />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Сурталчилгааны имэйл</p>
                  <p className="text-sm text-muted-foreground">
                    Шинэ боломжууд, санал урилгын мэдээлэл авах
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={notifications.marketingEmails}
                  onChange={(e) =>
                    setNotifications({
                      ...notifications,
                      marketingEmails: e.target.checked,
                    })
                  }
                  className="h-4 w-4 rounded border-gray-300"
                />
              </div>
            </div>
            <Button onClick={handleSaveNotifications}>Хадгалах</Button>
          </CardContent>
        </Card>

        {/* Privacy Settings */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Eye className="h-5 w-5" />
              Нууцлал
            </CardTitle>
            <CardDescription>
              Хувийн мэдээллийн харагдах байдлыг тохируулах
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Онлайн төлөв харуулах</p>
                  <p className="text-sm text-muted-foreground">
                    Бусад хэрэглэгчид таны онлайн төлөвийг харах боломжтой
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={privacy.showOnlineStatus}
                  onChange={(e) =>
                    setPrivacy({
                      ...privacy,
                      showOnlineStatus: e.target.checked,
                    })
                  }
                  className="h-4 w-4 rounded border-gray-300"
                />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Сүүлд идэвхтэй байсан цаг</p>
                  <p className="text-sm text-muted-foreground">
                    Сүүлд хэзээ идэвхтэй байсныг харуулах
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={privacy.showLastSeen}
                  onChange={(e) =>
                    setPrivacy({ ...privacy, showLastSeen: e.target.checked })
                  }
                  className="h-4 w-4 rounded border-gray-300"
                />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Имэйл хаяг харуулах</p>
                  <p className="text-sm text-muted-foreground">
                    Бусад хэрэглэгчид таны имэйл хаягийг харах боломжтой
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={privacy.showEmail}
                  onChange={(e) =>
                    setPrivacy({ ...privacy, showEmail: e.target.checked })
                  }
                  className="h-4 w-4 rounded border-gray-300"
                />
              </div>
            </div>
            <Button onClick={handleSavePrivacy}>Хадгалах</Button>
          </CardContent>
        </Card>

        {/* Connected Devices */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Smartphone className="h-5 w-5" />
              Холбогдсон төхөөрөмжүүд
            </CardTitle>
            <CardDescription>
              Таны бүртгэлд нэвтэрсэн төхөөрөмжүүдийн жагсаалт
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {mockDevices.map((device) => (
                <div
                  key={device.id}
                  className="flex items-center justify-between p-4 rounded-lg border"
                >
                  <div className="flex items-center gap-4">
                    {device.type === "desktop" ? (
                      <Laptop className="h-8 w-8 text-muted-foreground" />
                    ) : (
                      <Smartphone className="h-8 w-8 text-muted-foreground" />
                    )}
                    <div>
                      <p className="font-medium">
                        {device.name}
                        {device.current && (
                          <span className="ml-2 text-xs bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300 px-2 py-0.5 rounded">
                            Энэ төхөөрөмж
                          </span>
                        )}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {device.lastActive}
                      </p>
                    </div>
                  </div>
                  {!device.current && (
                    <Button
                      variant="ghost"
                      size="sm"
                      className="text-destructive hover:text-destructive"
                      onClick={() => handleRemoveDevice(device.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              ))}
            </div>
            <div className="mt-4">
              <Button variant="outline" className="text-destructive">
                Бусад бүх төхөөрөмжөөс гарах
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
