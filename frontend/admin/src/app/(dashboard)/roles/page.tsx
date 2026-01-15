"use client";

import { useEffect, useState } from "react";
import {
  Shield,
  Plus,
  Search,
  MoreHorizontal,
  Pencil,
  Trash2,
  Eye,
  Key,
  RefreshCw,
  CheckCircle2,
  XCircle,
  Building2,
  Settings,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";
import { rolesApi } from "@/lib/api";
import { PermissionGate } from "@/hooks/usePermission";

interface Role {
  id: number;
  code: string;
  name: string;
  description?: string;
  system_id?: number;
  system?: { id: number; name: string; code: string };
  is_system?: boolean;
  is_default?: boolean;
  organization_id?: number;
  is_active: boolean;
  created_date: string;
}

export default function RolesPage() {
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  const fetchRoles = async () => {
    setLoading(true);
    try {
      const response = await rolesApi.list({ page, page_size: 10 });
      setRoles(response.data.data || []);
      setTotalPages(response.data.meta?.total_pages || 1);
      setTotal(response.data.meta?.total || response.data.data?.length || 0);
    } catch (error) {
      console.error("Failed to fetch roles:", error);
      toast.error("Дүрүүдийг ачаалахад алдаа гарлаа");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRoles();
  }, [page]);

  const filteredRoles = roles.filter((role) => {
    if (!searchQuery) return true;
    const query = searchQuery.toLowerCase();
    return (
      role.name?.toLowerCase().includes(query) ||
      role.code?.toLowerCase().includes(query) ||
      role.description?.toLowerCase().includes(query)
    );
  });

  const handleDelete = async (roleId: number) => {
    if (!confirm("Энэ дүрийг устгах уу?")) return;

    try {
      await rolesApi.delete(roleId);
      toast.success("Дүр амжилттай устгагдлаа");
      fetchRoles();
    } catch (error) {
      toast.error("Дүр устгахад алдаа гарлаа");
    }
  };

  const getSystemColor = (systemId?: number) => {
    switch (systemId) {
      case 1:
        return "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400";
      case 2:
        return "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400";
      case 3:
        return "bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400";
      case 4:
        return "bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400";
      default:
        return "bg-gray-100 text-gray-700 dark:bg-gray-900/30 dark:text-gray-400";
    }
  };

  const getSystemName = (systemId?: number) => {
    switch (systemId) {
      case 1:
        return "Admin";
      case 2:
        return "DSL";
      default:
        return "Unknown";
    }
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2">
            <Shield className="h-6 w-6 text-primary" />
            Дүрүүд
          </h1>
          <p className="text-muted-foreground">
            Системийн дүрүүд болон эрхийн тохиргоо
          </p>
        </div>
        <PermissionGate permission="admin.role.create">
          <Button className="gap-2">
            <Plus className="h-4 w-4" />
            Дүр нэмэх
          </Button>
        </PermissionGate>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Нийт дүр
            </CardTitle>
            <Shield className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{total}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Идэвхтэй
            </CardTitle>
            <CheckCircle2 className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {roles.filter((r) => r.is_active !== false).length}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Системийн дүр
            </CardTitle>
            <Settings className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">
              {roles.filter((r) => r.is_system).length}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Байгууллагын дүр
            </CardTitle>
            <Building2 className="h-4 w-4 text-purple-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-purple-600">
              {roles.filter((r) => r.organization_id).length}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Roles Table */}
      <Card>
        <CardHeader>
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div>
              <CardTitle>Дүрийн жагсаалт</CardTitle>
              <CardDescription>
                Нийт {total} дүр бүртгэлтэй
              </CardDescription>
            </div>
            <div className="flex items-center gap-2">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Хайх..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-9 w-[250px]"
                />
              </div>
              <Button variant="outline" size="icon" onClick={fetchRoles}>
                <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin" : ""}`} />
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="space-y-4">
              {[...Array(5)].map((_, i) => (
                <div key={i} className="flex items-center gap-4">
                  <Skeleton className="h-12 w-12 rounded-lg" />
                  <div className="space-y-2 flex-1">
                    <Skeleton className="h-4 w-[250px]" />
                    <Skeleton className="h-3 w-[180px]" />
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="rounded-lg border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Дүр</TableHead>
                    <TableHead>Код</TableHead>
                    <TableHead>Систем</TableHead>
                    <TableHead>Төрөл</TableHead>
                    <TableHead>Төлөв</TableHead>
                    <TableHead className="w-[70px]"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredRoles.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={6} className="text-center py-10">
                        <div className="flex flex-col items-center gap-2 text-muted-foreground">
                          <Shield className="h-10 w-10" />
                          <p>Дүр олдсонгүй</p>
                        </div>
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredRoles.map((role) => (
                      <TableRow key={role.id}>
                        <TableCell>
                          <div className="flex items-center gap-3">
                            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                              <Shield className="h-5 w-5" />
                            </div>
                            <div>
                              <p className="font-medium">{role.name}</p>
                              {role.description && (
                                <p className="text-xs text-muted-foreground line-clamp-1">
                                  {role.description}
                                </p>
                              )}
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <code className="text-xs bg-muted px-2 py-1 rounded">
                            {role.code}
                          </code>
                        </TableCell>
                        <TableCell>
                          <Badge variant="secondary" className={getSystemColor(role.system_id)}>
                            {role.system?.name || getSystemName(role.system_id)}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            {role.is_system ? (
                              <Badge variant="outline" className="border-blue-500 text-blue-600">
                                <Settings className="h-3 w-3 mr-1" />
                                Системийн
                              </Badge>
                            ) : role.is_default ? (
                              <Badge variant="outline" className="border-green-500 text-green-600">
                                <CheckCircle2 className="h-3 w-3 mr-1" />
                                Үндсэн
                              </Badge>
                            ) : (
                              <Badge variant="outline">
                                Энгийн
                              </Badge>
                            )}
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant={role.is_active !== false ? "default" : "secondary"}
                            className={
                              role.is_active !== false
                                ? "bg-green-100 text-green-700 hover:bg-green-100 dark:bg-green-900/30 dark:text-green-400"
                                : "bg-red-100 text-red-700 hover:bg-red-100 dark:bg-red-900/30 dark:text-red-400"
                            }
                          >
                            {role.is_active !== false ? "Идэвхтэй" : "Идэвхгүй"}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="icon">
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuLabel>Үйлдэл</DropdownMenuLabel>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem>
                                <Eye className="mr-2 h-4 w-4" />
                                Дэлгэрэнгүй
                              </DropdownMenuItem>
                              <PermissionGate permission="admin.role.update">
                                <DropdownMenuItem>
                                  <Pencil className="mr-2 h-4 w-4" />
                                  Засах
                                </DropdownMenuItem>
                              </PermissionGate>
                              <DropdownMenuItem>
                                <Key className="mr-2 h-4 w-4" />
                                Эрхүүд
                              </DropdownMenuItem>
                              <DropdownMenuSeparator />
                              <PermissionGate permission="admin.role.delete">
                                <DropdownMenuItem
                                  className="text-destructive focus:text-destructive"
                                  onClick={() => handleDelete(role.id)}
                                >
                                  <Trash2 className="mr-2 h-4 w-4" />
                                  Устгах
                                </DropdownMenuItem>
                              </PermissionGate>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>
          )}

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between mt-4">
              <p className="text-sm text-muted-foreground">
                Хуудас {page} / {totalPages}
              </p>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page === 1}
                >
                  Өмнөх
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                  disabled={page === totalPages}
                >
                  Дараах
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
