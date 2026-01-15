"use client";

import { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { menusApi, systemsApi } from "@/lib/api";
import { toast } from "sonner";
import {
  Plus,
  Search,
  MoreHorizontal,
  Pencil,
  Trash2,
  Loader2,
  ChevronRight,
  Menu as MenuIcon,
  Eye,
  EyeOff,
  FolderOpen,
  RefreshCw,
  CheckCircle2,
  XCircle,
} from "lucide-react";
import * as LucideIcons from "lucide-react";
import { PermissionGate } from "@/hooks/usePermission";

interface Menu {
  id: number;
  code: string;
  name: string;
  system_id: number;
  parent_id: number | null;
  path: string;
  icon: string;
  component: string;
  sequence: number;
  is_visible: boolean;
  is_active: boolean;
  created_date: string;
  children?: Menu[];
}

interface System {
  id: number;
  code: string;
  name: string;
}

const defaultMenu: Partial<Menu> = {
  code: "",
  name: "",
  path: "",
  icon: "Circle",
  component: "",
  sequence: 0,
  is_visible: true,
  is_active: true,
};

export default function MenusPage() {
  const [menus, setMenus] = useState<Menu[]>([]);
  const [systems, setSystems] = useState<System[]>([]);
  const [selectedSystem, setSelectedSystem] = useState<string>("");
  const [isLoading, setIsLoading] = useState(true);
  const [search, setSearch] = useState("");

  // Dialog state
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingMenu, setEditingMenu] = useState<Partial<Menu> | null>(null);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    loadSystems();
  }, []);

  useEffect(() => {
    if (selectedSystem) {
      loadMenus();
    }
  }, [selectedSystem]);

  const loadSystems = async () => {
    try {
      const response = await systemsApi.list();
      const systemList = response.data.data?.systems || response.data.data || [];
      setSystems(systemList);
      if (systemList.length > 0) {
        setSelectedSystem(systemList[0].id.toString());
      }
    } catch (error) {
      toast.error("Системүүдийг татахад алдаа гарлаа");
    }
  };

  const loadMenus = async () => {
    setIsLoading(true);
    try {
      const response = await menusApi.list({ system_id: selectedSystem });
      setMenus(response.data.data || []);
    } catch (error) {
      toast.error("Цэсүүдийг татахад алдаа гарлаа");
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingMenu({
      ...defaultMenu,
      system_id: parseInt(selectedSystem),
    });
    setDialogOpen(true);
  };

  const handleEdit = (menu: Menu) => {
    setEditingMenu({ ...menu });
    setDialogOpen(true);
  };

  const handleSave = async () => {
    if (!editingMenu?.code || !editingMenu?.name) {
      toast.error("Код болон нэр оруулна уу");
      return;
    }

    setIsSaving(true);
    try {
      if (editingMenu.id) {
        await menusApi.update(editingMenu.id, editingMenu);
        toast.success("Цэс амжилттай шинэчлэгдлээ");
      } else {
        await menusApi.create(editingMenu);
        toast.success("Цэс амжилттай үүсгэгдлээ");
      }
      setDialogOpen(false);
      setEditingMenu(null);
      loadMenus();
    } catch (error) {
      toast.error("Хадгалахад алдаа гарлаа");
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Та энэ цэсийг устгахдаа итгэлтэй байна уу?")) return;

    try {
      await menusApi.delete(id);
      toast.success("Цэс амжилттай устгагдлаа");
      loadMenus();
    } catch (error) {
      toast.error("Устгахад алдаа гарлаа");
    }
  };

  // Get icon component
  const getIcon = (iconName: string) => {
    const icons = LucideIcons as unknown as Record<
      string,
      React.ComponentType<{ className?: string }>
    >;
    return icons[iconName] || LucideIcons.Circle;
  };

  const filteredMenus = menus.filter(
    (menu) =>
      menu.code.toLowerCase().includes(search.toLowerCase()) ||
      menu.name.toLowerCase().includes(search.toLowerCase()) ||
      menu.path.toLowerCase().includes(search.toLowerCase())
  );

  // Stats
  const parentMenus = menus.filter((m) => !m.parent_id);
  const childMenus = menus.filter((m) => m.parent_id);
  const hiddenMenus = menus.filter((m) => m.is_visible === false);


  // Get parent menu options
  const getParentOptions = () => {
    return menus.filter(
      (m) => m.id !== editingMenu?.id && m.parent_id === null
    );
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2">
            <MenuIcon className="h-6 w-6 text-primary" />
            Цэсүүд
          </h1>
          <p className="text-muted-foreground">
            Системийн цэсүүдийн удирдлага
          </p>
        </div>
        <PermissionGate permission="admin.menu.create">
          <Button onClick={handleCreate} disabled={!selectedSystem} className="gap-2">
            <Plus className="h-4 w-4" />
            Цэс нэмэх
          </Button>
        </PermissionGate>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Нийт цэс
            </CardTitle>
            <MenuIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{menus.length}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Үндсэн цэс
            </CardTitle>
            <FolderOpen className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {parentMenus.length}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Дэд цэс
            </CardTitle>
            <ChevronRight className="h-4 w-4 text-purple-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-purple-600">
              {childMenus.length}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Нуугдсан
            </CardTitle>
            <EyeOff className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {hiddenMenus.length}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Menus Table */}
      <Card>
        <CardHeader>
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div>
              <CardTitle>Цэсний жагсаалт</CardTitle>
              <CardDescription>
                Нийт {filteredMenus.length} цэс бүртгэлтэй
              </CardDescription>
            </div>
            <div className="flex items-center gap-2">
              <Select
                value={selectedSystem}
                onValueChange={setSelectedSystem}
              >
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="Систем сонгох" />
                </SelectTrigger>
                <SelectContent>
                  {systems.map((system) => (
                    <SelectItem key={system.id} value={system.id.toString()}>
                      {system.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <div className="relative">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Хайх..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-9 w-[200px]"
                />
              </div>
              <Button variant="outline" size="icon" onClick={loadMenus}>
                <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-4">
              {[...Array(5)].map((_, i) => (
                <div key={i} className="flex items-center gap-4">
                  <Skeleton className="h-10 w-10 rounded-lg" />
                  <div className="space-y-2 flex-1">
                    <Skeleton className="h-4 w-[200px]" />
                    <Skeleton className="h-3 w-[150px]" />
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="rounded-lg border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Цэс</TableHead>
                    <TableHead>Зам</TableHead>
                    <TableHead>Харагдах</TableHead>
                    <TableHead>Төлөв</TableHead>
                    <TableHead className="w-[70px]"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredMenus.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={5} className="text-center py-10">
                        <div className="flex flex-col items-center gap-2 text-muted-foreground">
                          <MenuIcon className="h-10 w-10" />
                          <p>Цэс олдсонгүй</p>
                          <PermissionGate permission="admin.menu.create">
                            <Button
                              variant="outline"
                              className="mt-2"
                              onClick={handleCreate}
                            >
                              <Plus className="mr-2 h-4 w-4" />
                              Эхний цэсээ үүсгэх
                            </Button>
                          </PermissionGate>
                        </div>
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredMenus.map((menu) => {
                      const Icon = getIcon(menu.icon);
                      return (
                        <TableRow key={menu.id}>
                          <TableCell>
                            <div className="flex items-center gap-3">
                              {menu.parent_id && (
                                <ChevronRight className="h-4 w-4 text-muted-foreground ml-4" />
                              )}
                              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                                <Icon className="h-5 w-5" />
                              </div>
                              <div>
                                <p className="font-medium">{menu.name}</p>
                                <code className="text-xs text-muted-foreground">
                                  {menu.code}
                                </code>
                              </div>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center gap-1">
                              <code className="text-xs bg-muted px-2 py-1 rounded">
                                {menu.path || "-"}
                              </code>
                            </div>
                          </TableCell>
                          <TableCell>
                            {menu.is_visible !== false ? (
                              <Badge variant="secondary" className="bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400">
                                <Eye className="h-3 w-3 mr-1" />
                                Тийм
                              </Badge>
                            ) : (
                              <Badge variant="secondary" className="bg-gray-100 text-gray-700 dark:bg-gray-900/30 dark:text-gray-400">
                                <EyeOff className="h-3 w-3 mr-1" />
                                Үгүй
                              </Badge>
                            )}
                          </TableCell>
                          <TableCell>
                            <Badge
                              variant={menu.is_active !== false ? "default" : "secondary"}
                              className={
                                menu.is_active !== false
                                  ? "bg-green-100 text-green-700 hover:bg-green-100 dark:bg-green-900/30 dark:text-green-400"
                                  : "bg-red-100 text-red-700 hover:bg-red-100 dark:bg-red-900/30 dark:text-red-400"
                              }
                            >
                              {menu.is_active !== false ? "Идэвхтэй" : "Идэвхгүй"}
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
                                <PermissionGate permission="admin.menu.update">
                                  <DropdownMenuItem onClick={() => handleEdit(menu)}>
                                    <Pencil className="mr-2 h-4 w-4" />
                                    Засах
                                  </DropdownMenuItem>
                                </PermissionGate>
                                <DropdownMenuSeparator />
                                <PermissionGate permission="admin.menu.delete">
                                  <DropdownMenuItem
                                    onClick={() => handleDelete(menu.id)}
                                    className="text-destructive focus:text-destructive"
                                  >
                                    <Trash2 className="mr-2 h-4 w-4" />
                                    Устгах
                                  </DropdownMenuItem>
                                </PermissionGate>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </TableCell>
                        </TableRow>
                      );
                    })
                  )}
                </TableBody>
              </Table>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Create/Edit Dialog */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>
              {editingMenu?.id ? "Цэс засах" : "Шинэ цэс үүсгэх"}
            </DialogTitle>
            <DialogDescription>
              Цэсийн мэдээллийг оруулна уу
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="code">Код *</Label>
                <Input
                  id="code"
                  value={editingMenu?.code || ""}
                  onChange={(e) =>
                    setEditingMenu((prev) => ({
                      ...prev!,
                      code: e.target.value,
                    }))
                  }
                  placeholder="dashboard"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="name">Нэр *</Label>
                <Input
                  id="name"
                  value={editingMenu?.name || ""}
                  onChange={(e) =>
                    setEditingMenu((prev) => ({
                      ...prev!,
                      name: e.target.value,
                    }))
                  }
                  placeholder="Хянах самбар"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="path">Зам (Path)</Label>
                <Input
                  id="path"
                  value={editingMenu?.path || ""}
                  onChange={(e) =>
                    setEditingMenu((prev) => ({
                      ...prev!,
                      path: e.target.value,
                    }))
                  }
                  placeholder="/dashboard"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="icon">Дүрс (Icon)</Label>
                <Input
                  id="icon"
                  value={editingMenu?.icon || ""}
                  onChange={(e) =>
                    setEditingMenu((prev) => ({
                      ...prev!,
                      icon: e.target.value,
                    }))
                  }
                  placeholder="LayoutDashboard"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="parent">Эцэг цэс</Label>
                <Select
                  value={editingMenu?.parent_id?.toString() || "none"}
                  onValueChange={(value) =>
                    setEditingMenu((prev) => ({
                      ...prev!,
                      parent_id: value === "none" ? null : parseInt(value),
                    }))
                  }
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Сонгоогүй" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">Үндсэн цэс</SelectItem>
                    {getParentOptions().map((menu) => (
                      <SelectItem key={menu.id} value={menu.id.toString()}>
                        {menu.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="sequence">Дараалал</Label>
                <Input
                  id="sequence"
                  type="number"
                  value={editingMenu?.sequence || 0}
                  onChange={(e) =>
                    setEditingMenu((prev) => ({
                      ...prev!,
                      sequence: parseInt(e.target.value) || 0,
                    }))
                  }
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  id="is_visible"
                  checked={editingMenu?.is_visible ?? true}
                  onChange={(e) =>
                    setEditingMenu((prev) => ({
                      ...prev!,
                      is_visible: e.target.checked,
                    }))
                  }
                  className="h-4 w-4 rounded border-gray-300"
                />
                <Label htmlFor="is_visible">Харагдах</Label>
              </div>
              <div className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  id="is_active"
                  checked={editingMenu?.is_active ?? true}
                  onChange={(e) =>
                    setEditingMenu((prev) => ({
                      ...prev!,
                      is_active: e.target.checked,
                    }))
                  }
                  className="h-4 w-4 rounded border-gray-300"
                />
                <Label htmlFor="is_active">Идэвхтэй</Label>
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDialogOpen(false)}
              disabled={isSaving}
            >
              Цуцлах
            </Button>
            <Button onClick={handleSave} disabled={isSaving}>
              {isSaving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Хадгалах
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
