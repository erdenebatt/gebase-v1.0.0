"use client";

import { Terminal, Play, Square, RotateCcw } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export default function ExecutorPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Executor</h1>
          <p className="text-muted-foreground">
            DSL expression болон workflow гүйцэтгэх
          </p>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Terminal className="h-5 w-5" />
              Expression Editor
            </CardTitle>
            <CardDescription>
              DSL expression бичиж шууд тест хийх
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Select>
              <SelectTrigger>
                <SelectValue placeholder="Төрөл сонгох" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="expression">Expression</SelectItem>
                <SelectItem value="rule">Rule Test</SelectItem>
                <SelectItem value="workflow">Workflow</SelectItem>
              </SelectContent>
            </Select>

            <Textarea
              placeholder="// DSL expression бичих...
sum(1, 2, 3)
concat('Hello', ' ', 'World')
if(age >= 18, 'Adult', 'Minor')"
              className="font-mono min-h-[200px]"
            />

            <div className="flex gap-2">
              <Button className="flex-1">
                <Play className="mr-2 h-4 w-4" />
                Гүйцэтгэх
              </Button>
              <Button variant="outline">
                <RotateCcw className="h-4 w-4" />
              </Button>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Үр дүн</CardTitle>
            <CardDescription>Гүйцэтгэлийн үр дүн энд харагдана</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="bg-muted rounded-lg p-4 min-h-[280px] font-mono text-sm">
              <span className="text-muted-foreground">
                // Гүйцэтгэх товч дарна уу...
              </span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
