"use client";

import { Database, Plus, Search, MoreHorizontal } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function SchemasPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Схемүүд</h1>
          <p className="text-muted-foreground">
            Динамик өгөгдлийн бүтэц тодорхойлох
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Шинэ схем
        </Button>
      </div>

      <div className="flex items-center gap-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input placeholder="Схем хайх..." className="pl-9" />
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card className="border-dashed">
          <CardHeader className="text-center">
            <Database className="h-12 w-12 mx-auto text-muted-foreground" />
            <CardTitle className="text-lg">Схем үүсгэх</CardTitle>
            <CardDescription>
              Шинэ динамик хүснэгт үүсгэхийн тулд схем тодорхойлно уу
            </CardDescription>
          </CardHeader>
          <CardContent className="text-center">
            <Button variant="outline">
              <Plus className="mr-2 h-4 w-4" />
              Үүсгэх
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
