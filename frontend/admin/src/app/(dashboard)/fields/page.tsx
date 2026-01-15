"use client";

import { Columns, Plus, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function FieldsPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Талбарууд</h1>
          <p className="text-muted-foreground">
            Схемийн талбаруудыг удирдах
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Шинэ талбар
        </Button>
      </div>

      <div className="flex items-center gap-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input placeholder="Талбар хайх..." className="pl-9" />
        </div>
      </div>

      <Card>
        <CardHeader className="text-center py-12">
          <Columns className="h-12 w-12 mx-auto text-muted-foreground" />
          <CardTitle>Талбар байхгүй</CardTitle>
          <CardDescription>
            Эхлээд схем сонгоод талбар нэмнэ үү
          </CardDescription>
        </CardHeader>
      </Card>
    </div>
  );
}
