"use client";

import { Variable, Plus, Search, Globe, Building, User, Key } from "lucide-react";
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

export default function VariablesPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Хувьсагчид</h1>
          <p className="text-muted-foreground">
            Тохиргооны хувьсагчид (Global, Organization, User, Session)
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Шинэ хувьсагч
        </Button>
      </div>

      <div className="flex items-center gap-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input placeholder="Хувьсагч хайх..." className="pl-9" />
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <Badge className="bg-blue-500">Global</Badge>
              <Globe className="h-4 w-4 text-muted-foreground" />
            </div>
            <CardTitle className="text-sm font-mono">APP_NAME</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">Gebase Platform</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <Badge className="bg-green-500">Global</Badge>
              <Globe className="h-4 w-4 text-muted-foreground" />
            </div>
            <CardTitle className="text-sm font-mono">DEFAULT_LANG</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">mn</p>
          </CardContent>
        </Card>

        <Card className="border-dashed">
          <CardHeader className="text-center py-8">
            <Plus className="h-8 w-8 mx-auto text-muted-foreground" />
            <CardDescription>Шинэ хувьсагч нэмэх</CardDescription>
          </CardHeader>
        </Card>
      </div>
    </div>
  );
}
