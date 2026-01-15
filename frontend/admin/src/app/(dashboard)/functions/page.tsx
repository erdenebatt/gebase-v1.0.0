"use client";

import { Code, Plus, Search, Play } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

export default function FunctionsPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Функцүүд</h1>
          <p className="text-muted-foreground">
            Built-in болон custom функцүүд
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Шинэ функц
        </Button>
      </div>

      <div className="flex items-center gap-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input placeholder="Функц хайх..." className="pl-9" />
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <Badge variant="secondary">builtin</Badge>
              <Code className="h-4 w-4 text-muted-foreground" />
            </div>
            <CardTitle className="text-lg">now()</CardTitle>
            <CardDescription>Одоогийн огноо цаг буцаана</CardDescription>
          </CardHeader>
          <CardContent>
            <code className="text-sm bg-muted px-2 py-1 rounded">
              now() → datetime
            </code>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <Badge variant="secondary">builtin</Badge>
              <Code className="h-4 w-4 text-muted-foreground" />
            </div>
            <CardTitle className="text-lg">concat()</CardTitle>
            <CardDescription>Текстүүдийг нэгтгэнэ</CardDescription>
          </CardHeader>
          <CardContent>
            <code className="text-sm bg-muted px-2 py-1 rounded">
              concat(a, b, ...) → string
            </code>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <Badge variant="secondary">builtin</Badge>
              <Code className="h-4 w-4 text-muted-foreground" />
            </div>
            <CardTitle className="text-lg">sum()</CardTitle>
            <CardDescription>Тоонуудын нийлбэр</CardDescription>
          </CardHeader>
          <CardContent>
            <code className="text-sm bg-muted px-2 py-1 rounded">
              sum(1, 2, 3) → number
            </code>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
