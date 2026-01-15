"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuthStore } from "@/store/authStore";
import { Button } from "@/components/ui/button";
import { ArrowRight, Shield, Zap, Users, Building2 } from "lucide-react";

export default function HomePage() {
  const router = useRouter();
  const { isAuthenticated, isLoading } = useAuthStore();

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.push("/dashboard");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      {/* Header */}
      <header className="border-b">
        <div className="container flex h-16 items-center justify-between">
          <Link href="/" className="flex items-center gap-2">
            <div className="h-8 w-8 rounded-lg bg-primary flex items-center justify-center">
              <span className="text-lg font-bold text-white">G</span>
            </div>
            <span className="font-semibold text-lg">Gebase</span>
          </Link>
          <div className="flex items-center gap-4">
            <Link href="/login">
              <Button variant="ghost">Нэвтрэх</Button>
            </Link>
            <Link href="/register">
              <Button>Бүртгүүлэх</Button>
            </Link>
          </div>
        </div>
      </header>

      {/* Hero */}
      <section className="py-20 md:py-32">
        <div className="container text-center">
          <h1 className="text-4xl md:text-6xl font-bold tracking-tight mb-6">
            Gebase Platform
          </h1>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto mb-8">
            Байгууллагын удирдлага, хэрэглэгчийн эрх, төхөөрөмжийн бүртгэл зэрэг
            бүх зүйлийг нэг дороос удирдах боломжтой платформ.
          </p>
          <div className="flex items-center justify-center gap-4">
            <Link href="/register">
              <Button size="lg" className="gap-2">
                Эхлэх
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
            <Link href="/about">
              <Button size="lg" variant="outline">
                Дэлгэрэнгүй
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="py-20 bg-gray-50 dark:bg-gray-900">
        <div className="container">
          <h2 className="text-3xl font-bold text-center mb-12">
            Платформын боломжууд
          </h2>
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            <div className="p-6 bg-white dark:bg-gray-800 rounded-xl shadow-sm">
              <div className="h-12 w-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <Users className="h-6 w-6 text-primary" />
              </div>
              <h3 className="font-semibold mb-2">Хэрэглэгчийн удирдлага</h3>
              <p className="text-sm text-muted-foreground">
                Хэрэглэгчдийг бүртгэх, эрх олгох, хяналт тавих
              </p>
            </div>
            <div className="p-6 bg-white dark:bg-gray-800 rounded-xl shadow-sm">
              <div className="h-12 w-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <Building2 className="h-6 w-6 text-primary" />
              </div>
              <h3 className="font-semibold mb-2">Байгууллагын бүтэц</h3>
              <p className="text-sm text-muted-foreground">
                Байгууллагын шатлал, хэлтэс нэгжүүдийг удирдах
              </p>
            </div>
            <div className="p-6 bg-white dark:bg-gray-800 rounded-xl shadow-sm">
              <div className="h-12 w-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <Shield className="h-6 w-6 text-primary" />
              </div>
              <h3 className="font-semibold mb-2">Эрхийн удирдлага</h3>
              <p className="text-sm text-muted-foreground">
                Нарийвчилсан эрхийн тохиргоо, хандалтын хяналт
              </p>
            </div>
            <div className="p-6 bg-white dark:bg-gray-800 rounded-xl shadow-sm">
              <div className="h-12 w-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <Zap className="h-6 w-6 text-primary" />
              </div>
              <h3 className="font-semibold mb-2">Олон платформ</h3>
              <p className="text-sm text-muted-foreground">
                Вэб, гар утас, таблет, киоск зэрэг бүх төхөөрөмжид
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="py-20">
        <div className="container text-center">
          <h2 className="text-3xl font-bold mb-4">Бидэнтэй нэгдээрэй</h2>
          <p className="text-muted-foreground mb-8 max-w-xl mx-auto">
            Gebase платформыг ашиглан байгууллагынхаа удирдлагыг дараагийн
            түвшинд гаргаарай.
          </p>
          <Link href="/register">
            <Button size="lg">Үнэгүй эхлэх</Button>
          </Link>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t py-8">
        <div className="container text-center text-sm text-muted-foreground">
          <p>&copy; 2024 Gebase Platform. Бүх эрх хуулиар хамгаалагдсан.</p>
        </div>
      </footer>
    </div>
  );
}
