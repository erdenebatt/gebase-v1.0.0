"use client";

import Link from "next/link";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/store/authStore";
import {
  User,
  Settings,
  Shield,
  Smartphone,
  ArrowRight,
  Activity,
} from "lucide-react";

export default function DashboardPage() {
  const { user, availableSystems } = useAuthStore();

  const quickActions = [
    {
      title: "Профайл",
      description: "Хувийн мэдээллээ удирдах",
      icon: User,
      href: "/profile",
    },
    {
      title: "Тохиргоо",
      description: "Аккаунтын тохиргоо",
      icon: Settings,
      href: "/settings",
    },
    {
      title: "Нууцлал",
      description: "Нууцлалын тохиргоо",
      icon: Shield,
      href: "/settings/security",
    },
    {
      title: "Төхөөрөмжүүд",
      description: "Холбогдсон төхөөрөмжүүд",
      icon: Smartphone,
      href: "/sessions",
    },
  ];

  return (
    <div className="container py-8">
      {/* Welcome */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight mb-2">
          Сайн байна уу, {user?.first_name}!
        </h1>
        <p className="text-muted-foreground">
          Gebase Portal-д тавтай морил. Энд та өөрийн профайл, тохиргоо, системүүдээ
          удирдах боломжтой.
        </p>
      </div>

      {/* Quick Actions */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-8">
        {quickActions.map((action) => (
          <Link key={action.href} href={action.href}>
            <Card className="hover:bg-accent transition-colors cursor-pointer h-full">
              <CardHeader className="pb-2">
                <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center mb-2">
                  <action.icon className="h-5 w-5 text-primary" />
                </div>
                <CardTitle className="text-base">{action.title}</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription>{action.description}</CardDescription>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>

      {/* Available Systems */}
      {availableSystems.length > 0 && (
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>Миний системүүд</CardTitle>
            <CardDescription>
              Та дараах системүүдэд хандах эрхтэй байна
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {availableSystems.map((systemRole) => (
                <Card
                  key={systemRole.system.id}
                  className="hover:bg-accent transition-colors cursor-pointer"
                >
                  <CardHeader className="pb-2">
                    <div className="flex items-center gap-3">
                      <div
                        className="h-12 w-12 rounded-lg flex items-center justify-center"
                        style={{
                          backgroundColor: systemRole.system.color || "#6366f1",
                        }}
                      >
                        <span className="text-xl font-bold text-white">
                          {systemRole.system.name[0]}
                        </span>
                      </div>
                      <div className="flex-1">
                        <CardTitle className="text-lg">
                          {systemRole.system.name}
                        </CardTitle>
                        <CardDescription>
                          {systemRole.roles[0]?.name}
                        </CardDescription>
                      </div>
                      <ArrowRight className="h-5 w-5 text-muted-foreground" />
                    </div>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm text-muted-foreground">
                      {systemRole.system.description || "Системийн тайлбар"}
                    </p>
                  </CardContent>
                </Card>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Recent Activity */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Activity className="h-5 w-5" />
            Сүүлийн үйлдлүүд
          </CardTitle>
          <CardDescription>
            Таны аккаунт дээр хийгдсэн сүүлийн үйлдлүүд
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {[
              {
                action: "Систем нэвтрэлт",
                device: "Chrome, Windows",
                time: "5 минутын өмнө",
                location: "Улаанбаатар",
              },
              {
                action: "Нууц үг солигдлоо",
                device: "Chrome, Windows",
                time: "2 өдрийн өмнө",
                location: "Улаанбаатар",
              },
              {
                action: "Систем нэвтрэлт",
                device: "Safari, iPhone",
                time: "3 өдрийн өмнө",
                location: "Улаанбаатар",
              },
            ].map((activity, index) => (
              <div
                key={index}
                className="flex items-center justify-between py-3 border-b last:border-0"
              >
                <div>
                  <p className="font-medium">{activity.action}</p>
                  <p className="text-sm text-muted-foreground">
                    {activity.device} • {activity.location}
                  </p>
                </div>
                <span className="text-sm text-muted-foreground">
                  {activity.time}
                </span>
              </div>
            ))}
          </div>
          <Button variant="outline" className="w-full mt-4">
            Бүгдийг харах
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
