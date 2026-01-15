"use client";

import { GitBranch, Plus, Search, Play, Pause } from "lucide-react";
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

export default function WorkflowsPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Workflow</h1>
          <p className="text-muted-foreground">
            Ажлын урсгал ба автоматжуулалт
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Шинэ workflow
        </Button>
      </div>

      <div className="flex items-center gap-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input placeholder="Workflow хайх..." className="pl-9" />
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card className="border-dashed">
          <CardHeader className="text-center">
            <GitBranch className="h-12 w-12 mx-auto text-muted-foreground" />
            <CardTitle className="text-lg">Workflow үүсгэх</CardTitle>
            <CardDescription>
              Олон алхамт процесс автоматжуулах
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
