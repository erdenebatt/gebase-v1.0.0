"use client";

import { useEffect, useState } from "react";
import {
  Building2,
  Plus,
  Search,
  MoreHorizontal,
  Pencil,
  Trash2,
  Eye,
  MapPin,
  Mail,
  Phone,
  RefreshCw,
  Network,
  CheckCircle2,
  XCircle,
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
import { organizationsApi } from "@/lib/api";
import { PermissionGate } from "@/hooks/usePermission";

interface Organization {
  id: number;
  name: string;
  short_name?: string;
  reg_no?: string;
  email?: string;
  phone_no?: string;
  type_id?: number;
  parent_id?: number;
  sso_org_id?: number;
  is_active: boolean;
  created_date: string;
}

export default function OrganizationsPage() {
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  const fetchOrganizations = async () => {
    setLoading(true);
    try {
      const response = await organizationsApi.list({ page, page_size: 10 });
      setOrganizations(response.data.data || []);
      setTotalPages(response.data.meta?.total_pages || 1);
      setTotal(response.data.meta?.total || response.data.data?.length || 0);
    } catch (error) {
      console.error("Failed to fetch organizations:", error);
      toast.error("Байгууллагуудыг ачаалахад алдаа гарлаа");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrganizations();
  }, [page]);

  const filteredOrgs = organizations.filter((org) => {
    if (!searchQuery) return true;
    const query = searchQuery.toLowerCase();
    return (
      org.name?.toLowerCase().includes(query) ||
      org.short_name?.toLowerCase().includes(query) ||
      org.email?.toLowerCase().includes(query) ||
      org.phone_no?.includes(query) ||
      org.reg_no?.toLowerCase().includes(query)
    );
  });

  const handleDelete = async (orgId: number) => {
    if (!confirm("Энэ байгууллагыг устгах уу?")) return;

    try {
      await organizationsApi.delete(orgId);
      toast.success("Байгууллага амжилттай устгагдлаа");
      fetchOrganizations();
    } catch (error) {
      toast.error("Байгууллага устгахад алдаа гарлаа");
    }
  };

  const getTypeColor = (typeId?: number) => {
    switch (typeId) {
      case 1:
        return "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400";
      case 2:
        return "bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400";
      case 3:
        return "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400";
      case 4:
        return "bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400";
      case 5:
        return "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400";
      default:
        return "bg-gray-100 text-gray-700 dark:bg-gray-900/30 dark:text-gray-400";
    }
  };

  const getTypeName = (typeId?: number) => {
    switch (typeId) {
      case 1:
        return "Төрийн байгууллага";
      case 2:
        return "Хувийн хэвшил";
      case 3:
        return "ТББ";
      case 4:
        return "Боловсролын";
      case 5:
        return "Эрүүл мэндийн";
      default:
        return "Бусад";
    }
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2">
            <Building2 className="h-6 w-6 text-primary" />
            Байгууллагууд
          </h1>
          <p className="text-muted-foreground">
            Системд бүртгэлтэй байгууллагуудын жагсаалт
          </p>
        </div>
        <PermissionGate permission="admin.organization.create">
          <Button className="gap-2">
            <Plus className="h-4 w-4" />
            Байгууллага нэмэх
          </Button>
        </PermissionGate>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Нийт байгууллага
            </CardTitle>
            <Building2 className="h-4 w-4 text-muted-foreground" />
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
              {organizations.filter((o) => o.is_active !== false).length}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Идэвхгүй
            </CardTitle>
            <XCircle className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              {organizations.filter((o) => o.is_active === false).length}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Салбартай
            </CardTitle>
            <Network className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary">
              {organizations.filter((o) => o.parent_id).length}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Organizations Table */}
      <Card>
        <CardHeader>
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div>
              <CardTitle>Байгууллагын жагсаалт</CardTitle>
              <CardDescription>
                Нийт {total} байгууллага бүртгэлтэй
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
              <Button variant="outline" size="icon" onClick={fetchOrganizations}>
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
                    <TableHead>Байгууллага</TableHead>
                    <TableHead>Холбоо барих</TableHead>
                    <TableHead>Регистр</TableHead>
                    <TableHead>Төрөл</TableHead>
                    <TableHead>Төлөв</TableHead>
                    <TableHead className="w-[70px]"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredOrgs.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={6} className="text-center py-10">
                        <div className="flex flex-col items-center gap-2 text-muted-foreground">
                          <Building2 className="h-10 w-10" />
                          <p>Байгууллага олдсонгүй</p>
                        </div>
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredOrgs.map((org) => (
                      <TableRow key={org.id}>
                        <TableCell>
                          <div className="flex items-center gap-3">
                            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary font-bold">
                              {org.short_name?.[0] || org.name?.[0] || "O"}
                            </div>
                            <div>
                              <p className="font-medium">{org.name}</p>
                              <p className="text-xs text-muted-foreground">
                                ID: {org.sso_org_id || org.id}
                              </p>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="space-y-1">
                            {org.email && (
                              <div className="flex items-center gap-1 text-sm">
                                <Mail className="h-3 w-3 text-muted-foreground" />
                                {org.email}
                              </div>
                            )}
                            {org.phone_no && (
                              <div className="flex items-center gap-1 text-sm text-muted-foreground">
                                <Phone className="h-3 w-3" />
                                {org.phone_no}
                              </div>
                            )}
                            {!org.email && !org.phone_no && (
                              <span className="text-muted-foreground text-sm">-</span>
                            )}
                          </div>
                        </TableCell>
                        <TableCell>
                          <code className="text-xs bg-muted px-2 py-1 rounded">
                            {org.reg_no || "-"}
                          </code>
                        </TableCell>
                        <TableCell>
                          <Badge variant="secondary" className={getTypeColor(org.type_id)}>
                            {getTypeName(org.type_id)}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant={org.is_active !== false ? "default" : "secondary"}
                            className={
                              org.is_active !== false
                                ? "bg-green-100 text-green-700 hover:bg-green-100 dark:bg-green-900/30 dark:text-green-400"
                                : "bg-red-100 text-red-700 hover:bg-red-100 dark:bg-red-900/30 dark:text-red-400"
                            }
                          >
                            {org.is_active !== false ? "Идэвхтэй" : "Идэвхгүй"}
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
                              <PermissionGate permission="admin.organization.update">
                                <DropdownMenuItem>
                                  <Pencil className="mr-2 h-4 w-4" />
                                  Засах
                                </DropdownMenuItem>
                              </PermissionGate>
                              <DropdownMenuItem>
                                <Network className="mr-2 h-4 w-4" />
                                Салбарууд
                              </DropdownMenuItem>
                              <DropdownMenuSeparator />
                              <PermissionGate permission="admin.organization.delete">
                                <DropdownMenuItem
                                  className="text-destructive focus:text-destructive"
                                  onClick={() => handleDelete(org.id)}
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
