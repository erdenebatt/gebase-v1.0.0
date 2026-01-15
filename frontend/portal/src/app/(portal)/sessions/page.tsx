"use client";

import { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { portalApi } from "@/lib/api";
import { toast } from "sonner";
import {
  Smartphone,
  Monitor,
  Tablet,
  Loader2,
  LogOut,
  MapPin,
  Clock,
} from "lucide-react";

interface Session {
  id: number;
  device_name: string;
  platform: string;
  ip_address: string;
  location: string;
  last_activity: string;
  is_current: boolean;
}

export default function SessionsPage() {
  const [sessions, setSessions] = useState<Session[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadSessions();
  }, []);

  const loadSessions = async () => {
    try {
      const response = await portalApi.sessions.list();
      setSessions(response.data.data || []);
    } catch (error) {
      // Mock data for demo
      setSessions([
        {
          id: 1,
          device_name: "Chrome - Windows",
          platform: "web",
          ip_address: "202.131.1.xxx",
          location: "Улаанбаатар, Монгол",
          last_activity: new Date().toISOString(),
          is_current: true,
        },
        {
          id: 2,
          device_name: "Safari - iPhone",
          platform: "ios",
          ip_address: "202.131.2.xxx",
          location: "Улаанбаатар, Монгол",
          last_activity: new Date(
            Date.now() - 3600000 * 24
          ).toISOString(),
          is_current: false,
        },
        {
          id: 3,
          device_name: "Chrome - MacBook",
          platform: "mac_desktop",
          ip_address: "202.131.3.xxx",
          location: "Дархан, Монгол",
          last_activity: new Date(
            Date.now() - 3600000 * 48
          ).toISOString(),
          is_current: false,
        },
      ]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRevoke = async (id: number) => {
    if (!confirm("Та энэ сешнийг хүчингүй болгохдоо итгэлтэй байна уу?")) return;

    try {
      await portalApi.sessions.revoke(id);
      toast.success("Сешн хүчингүй болгогдлоо");
      loadSessions();
    } catch (error) {
      toast.error("Алдаа гарлаа");
    }
  };

  const handleRevokeAll = async () => {
    if (
      !confirm(
        "Та бүх сешнийг хүчингүй болгохдоо итгэлтэй байна уу? Энэ үйлдлийг хийсний дараа дахин нэвтрэх шаардлагатай."
      )
    )
      return;

    try {
      await portalApi.sessions.revokeAll();
      toast.success("Бүх сешн хүчингүй болгогдлоо");
      window.location.href = "/login";
    } catch (error) {
      toast.error("Алдаа гарлаа");
    }
  };

  const getDeviceIcon = (platform: string) => {
    switch (platform) {
      case "ios":
      case "android":
        return Smartphone;
      case "tablet_ios":
      case "tablet_android":
        return Tablet;
      default:
        return Monitor;
    }
  };

  const formatLastActivity = (dateStr: string) => {
    const date = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return "Одоо идэвхтэй";
    if (diffMins < 60) return `${diffMins} минутын өмнө`;
    if (diffHours < 24) return `${diffHours} цагийн өмнө`;
    return `${diffDays} өдрийн өмнө`;
  };

  return (
    <div className="container py-8 max-w-3xl">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">
            Идэвхтэй сешнүүд
          </h1>
          <p className="text-muted-foreground">
            Таны аккаунтад нэвтэрсэн төхөөрөмжүүд
          </p>
        </div>
        <Button variant="destructive" onClick={handleRevokeAll}>
          <LogOut className="mr-2 h-4 w-4" />
          Бүгдийг гаргах
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Төхөөрөмжүүд</CardTitle>
          <CardDescription>
            Нийт {sessions.length} төхөөрөмж холбогдсон байна
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : (
            <div className="space-y-4">
              {sessions.map((session) => {
                const DeviceIcon = getDeviceIcon(session.platform);
                return (
                  <div
                    key={session.id}
                    className="flex items-start gap-4 p-4 rounded-lg border"
                  >
                    <div className="h-12 w-12 rounded-lg bg-muted flex items-center justify-center">
                      <DeviceIcon className="h-6 w-6 text-muted-foreground" />
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <p className="font-medium">{session.device_name}</p>
                        {session.is_current && (
                          <span className="inline-flex items-center rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700">
                            Энэ төхөөрөмж
                          </span>
                        )}
                      </div>
                      <div className="mt-1 flex items-center gap-4 text-sm text-muted-foreground">
                        <span className="flex items-center gap-1">
                          <MapPin className="h-3 w-3" />
                          {session.location}
                        </span>
                        <span className="flex items-center gap-1">
                          <Clock className="h-3 w-3" />
                          {formatLastActivity(session.last_activity)}
                        </span>
                      </div>
                      <p className="mt-1 text-xs text-muted-foreground">
                        IP: {session.ip_address}
                      </p>
                    </div>
                    {!session.is_current && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleRevoke(session.id)}
                      >
                        Гаргах
                      </Button>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
