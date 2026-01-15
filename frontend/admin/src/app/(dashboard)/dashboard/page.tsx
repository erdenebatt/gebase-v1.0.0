"use client";

import { useEffect, useState } from "react";
import {
  Users,
  Building2,
  Smartphone,
  Shield,
  Activity,
  UserCheck,
  Monitor,
  Settings,
  Key,
  ArrowUpRight,
  Layers,
} from "lucide-react";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { useAuthStore } from "@/store/authStore";
import { useSystemStore } from "@/store/systemStore";

interface StatCardProps {
  title: string;
  value: string | number;
  description?: string;
  icon: React.ReactNode;
  gradient?: string;
  accentColor?: string;
}

function StatCard({ title, value, description, icon, gradient, accentColor }: StatCardProps) {
  return (
    <Card className="relative overflow-hidden">
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle>
        <div className={`p-2.5 rounded-xl bg-gradient-to-br ${gradient || "from-primary/20 to-primary/10"}`}>
          {icon}
        </div>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        {description && (
          <p className="text-xs text-muted-foreground mt-1">{description}</p>
        )}
      </CardContent>
      <div className={`absolute bottom-0 left-0 right-0 h-1 ${accentColor || "bg-primary"}`} />
    </Card>
  );
}

interface RecentActivityItem {
  id: number;
  action: string;
  user: string;
  time: string;
  type: "user" | "system" | "device" | "org";
}

const recentActivities: RecentActivityItem[] = [
  { id: 1, action: "Шинэ хэрэглэгч бүртгэгдлээ", user: "Болд", time: "2 минутын өмнө", type: "user" },
  { id: 2, action: "Систем шинэчлэгдлээ", user: "System", time: "15 минутын өмнө", type: "system" },
  { id: 3, action: "Төхөөрөмж холбогдлоо", user: "POS Terminal", time: "1 цагийн өмнө", type: "device" },
  { id: 4, action: "Байгууллага нэмэгдлээ", user: "Админ", time: "2 цагийн өмнө", type: "org" },
  { id: 5, action: "Эрх шинэчлэгдлээ", user: "System", time: "3 цагийн өмнө", type: "system" },
];

const getActivityIcon = (type: string) => {
  switch (type) {
    case "user":
      return <UserCheck className="h-4 w-4 text-blue-500" />;
    case "system":
      return <Activity className="h-4 w-4 text-purple-500" />;
    case "device":
      return <Smartphone className="h-4 w-4 text-green-500" />;
    case "org":
      return <Building2 className="h-4 w-4 text-orange-500" />;
    default:
      return <Activity className="h-4 w-4" />;
  }
};

export default function DashboardPage() {
  const { user, availableSystems } = useAuthStore();
  const { currentSystem, currentRole, permissions } = useSystemStore();

  const [stats, setStats] = useState({
    userCount: 0,
    orgCount: 0,
    deviceCount: 0,
    sessionCount: 0,
  });

  useEffect(() => {
    // Simulate stats - replace with actual API calls
    setStats({
      userCount: 156,
      orgCount: 24,
      deviceCount: 89,
      sessionCount: 3,
    });
  }, []);

  return (
    <div className="space-y-6">
      {/* Welcome Section */}
      <div className="flex flex-col gap-1">
        <h1 className="text-2xl font-bold tracking-tight">
          Сайн байна уу, {user?.first_name || "Хэрэглэгч"}! 👋
        </h1>
        <p className="text-muted-foreground">
          Админ системийн удирдлагын хяналтын самбар
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
        <StatCard
          title="Нийт хэрэглэгч"
          value={stats.userCount}
          description="бүртгэлтэй"
          icon={<Users className="h-5 w-5 text-blue-600" />}
          gradient="from-blue-100 to-blue-50 dark:from-blue-900/40 dark:to-blue-900/20"
          accentColor="bg-blue-500"
        />
        <StatCard
          title="Байгууллага"
          value={stats.orgCount}
          description="бүртгэлтэй"
          icon={<Building2 className="h-5 w-5 text-green-600" />}
          gradient="from-green-100 to-green-50 dark:from-green-900/40 dark:to-green-900/20"
          accentColor="bg-green-500"
        />
        <StatCard
          title="Төхөөрөмж"
          value={stats.deviceCount}
          description="бүртгэлтэй"
          icon={<Smartphone className="h-5 w-5 text-orange-600" />}
          gradient="from-orange-100 to-orange-50 dark:from-orange-900/40 dark:to-orange-900/20"
          accentColor="bg-orange-500"
        />
        <StatCard
          title="Идэвхтэй сессүүд"
          value={stats.sessionCount}
          description="миний төхөөрөмжүүд"
          icon={<Monitor className="h-5 w-5 text-cyan-600" />}
          gradient="from-cyan-100 to-cyan-50 dark:from-cyan-900/40 dark:to-cyan-900/20"
          accentColor="bg-cyan-500"
        />
        <StatCard
          title="Эрхүүд"
          value={permissions.length}
          description="таньд олгогдсон"
          icon={<Shield className="h-5 w-5 text-purple-600" />}
          gradient="from-purple-100 to-purple-50 dark:from-purple-900/40 dark:to-purple-900/20"
          accentColor="bg-purple-500"
        />
      </div>

      {/* Main Content Grid */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {/* System Overview */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Layers className="h-5 w-5 text-primary" />
              Системийн тойм
            </CardTitle>
            <CardDescription>Таны хандах боломжтой системүүд</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 sm:grid-cols-2">
              {availableSystems.map((systemRole) => (
                <div
                  key={systemRole.system.id}
                  className={`flex items-center gap-4 p-4 rounded-xl border transition-all cursor-pointer ${
                    currentSystem?.id === systemRole.system.id
                      ? "border-primary bg-primary/5"
                      : "hover:border-primary/50"
                  }`}
                >
                  <div
                    className="flex h-12 w-12 items-center justify-center rounded-xl text-white font-bold"
                    style={{ backgroundColor: systemRole.system.color || "#6366f1" }}
                  >
                    <Layers className="h-6 w-6" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <p className="font-medium truncate">{systemRole.system.name}</p>
                      {currentSystem?.id === systemRole.system.id && (
                        <Badge variant="secondary" className="text-xs">Идэвхтэй</Badge>
                      )}
                    </div>
                    <p className="text-sm text-muted-foreground truncate">{systemRole.system.code}</p>
                  </div>
                  <ArrowUpRight className="h-4 w-4 text-muted-foreground" />
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Recent Activity */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Activity className="h-5 w-5 text-primary" />
              Сүүлийн үйлдлүүд
            </CardTitle>
            <CardDescription>Системд болсон үйлдлүүд</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentActivities.map((activity) => (
                <div key={activity.id} className="flex items-start gap-3">
                  <div className="flex h-8 w-8 items-center justify-center rounded-full bg-muted">
                    {getActivityIcon(activity.type)}
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">{activity.action}</p>
                    <p className="text-xs text-muted-foreground">
                      {activity.user} • {activity.time}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card className="bg-gradient-to-br from-blue-500 to-blue-600 text-white transition-colors cursor-pointer hover:from-blue-600 hover:to-blue-700">
          <CardContent className="flex items-center gap-4 p-6">
            <div className="h-12 w-12 rounded-full bg-white/20 flex items-center justify-center">
              <Users className="h-6 w-6" />
            </div>
            <div>
              <p className="font-semibold">Хэрэглэгч нэмэх</p>
              <p className="text-sm opacity-80">Шинэ хэрэглэгч бүртгэх</p>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-green-500 to-green-600 text-white transition-colors cursor-pointer hover:from-green-600 hover:to-green-700">
          <CardContent className="flex items-center gap-4 p-6">
            <div className="h-12 w-12 rounded-full bg-white/20 flex items-center justify-center">
              <Building2 className="h-6 w-6" />
            </div>
            <div>
              <p className="font-semibold">Байгууллага нэмэх</p>
              <p className="text-sm opacity-80">Шинэ байгууллага бүртгэх</p>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-purple-500 to-purple-600 text-white transition-colors cursor-pointer hover:from-purple-600 hover:to-purple-700">
          <CardContent className="flex items-center gap-4 p-6">
            <div className="h-12 w-12 rounded-full bg-white/20 flex items-center justify-center">
              <Shield className="h-6 w-6" />
            </div>
            <div>
              <p className="font-semibold">Дүр удирдах</p>
              <p className="text-sm opacity-80">Эрх тохируулах</p>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-orange-500 to-orange-600 text-white transition-colors cursor-pointer hover:from-orange-600 hover:to-orange-700">
          <CardContent className="flex items-center gap-4 p-6">
            <div className="h-12 w-12 rounded-full bg-white/20 flex items-center justify-center">
              <Settings className="h-6 w-6" />
            </div>
            <div>
              <p className="font-semibold">Тохиргоо</p>
              <p className="text-sm opacity-80">Системийн тохиргоо</p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Permissions List */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Key className="h-5 w-5 text-primary" />
            Таны эрхүүд
          </CardTitle>
          <CardDescription>Админ системд таньд олгогдсон эрхүүд</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-2">
            {permissions.length > 0 ? (
              permissions.map((perm) => (
                <Badge key={perm} variant="secondary" className="text-xs">
                  {perm}
                </Badge>
              ))
            ) : (
              <p className="text-sm text-muted-foreground">Эрх олдсонгүй</p>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
