"use client";

import { FileCode, Plus, Search, Mail, MessageSquare, FileText } from "lucide-react";
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

export default function TemplatesPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Загварууд</h1>
          <p className="text-muted-foreground">
            Email, SMS, Document загварууд
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Шинэ загвар
        </Button>
      </div>

      <div className="flex items-center gap-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input placeholder="Загвар хайх..." className="pl-9" />
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card className="border-dashed">
          <CardHeader className="text-center">
            <Mail className="h-12 w-12 mx-auto text-muted-foreground" />
            <CardTitle className="text-lg">Email загвар</CardTitle>
            <CardDescription>
              И-мэйл мэдэгдлийн загвар
            </CardDescription>
          </CardHeader>
          <CardContent className="text-center">
            <Button variant="outline">
              <Plus className="mr-2 h-4 w-4" />
              Үүсгэх
            </Button>
          </CardContent>
        </Card>

        <Card className="border-dashed">
          <CardHeader className="text-center">
            <MessageSquare className="h-12 w-12 mx-auto text-muted-foreground" />
            <CardTitle className="text-lg">SMS загвар</CardTitle>
            <CardDescription>
              SMS мэдэгдлийн загвар
            </CardDescription>
          </CardHeader>
          <CardContent className="text-center">
            <Button variant="outline">
              <Plus className="mr-2 h-4 w-4" />
              Үүсгэх
            </Button>
          </CardContent>
        </Card>

        <Card className="border-dashed">
          <CardHeader className="text-center">
            <FileText className="h-12 w-12 mx-auto text-muted-foreground" />
            <CardTitle className="text-lg">Document загвар</CardTitle>
            <CardDescription>
              Баримт бичгийн загвар
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
