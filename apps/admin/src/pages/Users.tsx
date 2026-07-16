import { useEffect, useState, useCallback } from "react";
import { apiRequest } from "@/lib/api";
import { useToast } from "@/lib/toast";
import { useConfirm } from "@/lib/confirm";
import { useModuleAccess } from "@/lib/permissions";
import { Button } from "@/components/ui/button";
import { Dialog } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Select } from "@/components/ui/select";
import { Field } from "@/components/ui/field";
import { RequiredLegend } from "@/components/ui/required-legend";
import {
  Table,
  TableHead,
  TableBody,
  TableRow,
  TableHeader,
  TableCell,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Plus, Pencil, Trash2, User, CreditCard } from "lucide-react";
import { useRoles, roleLabel } from "@/lib/roles";

interface OrgUser {
  id: string;
  role_id: string;
  role: string;
  name: string;
  email: string;
  status: string;
  is_principal: boolean;
  managed_by_user_id?: string;
  managed_by_name?: string;
  plan_id?: string;
  plan_name?: string;
}

interface FarmOption {
  id: string;
  name: string;
}

export function Users() {
  const [users, setUsers] = useState<OrgUser[]>([]);
  const [farms, setFarms] = useState<FarmOption[]>([]);
  const [loading, setLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editing, setEditing] = useState<OrgUser | null>(null);
  const [form, setForm] = useState({
    name: "",
    email: "",
    password: "",
    role_id: "",
    status: "active",
    farm_ids: [] as string[],
  });
  const [saving, setSaving] = useState(false);
  const toast = useToast();
  const confirm = useConfirm();
  const { roles } = useRoles();
  const canEdit = useModuleAccess("users") === "write";

  const load = useCallback(async () => {
    try {
      const [userData, farmData] = await Promise.all([
        apiRequest<OrgUser[]>("/users"),
        apiRequest<FarmOption[]>("/farms"),
      ]);
      setUsers(userData);
      setFarms(farmData);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      if (editing) {
        await apiRequest(`/users/${editing.id}`, {
          method: "PUT",
          body: {
            name: form.name,
            email: form.email,
            role_id: form.role_id,
            is_active: form.status === "active",
          },
        });
      } else {
        await apiRequest("/users", {
          method: "POST",
          body: {
            name: form.name,
            email: form.email,
            password: form.password,
            role_id: form.role_id,
            farm_ids: form.farm_ids,
          },
        });
      }
      setDialogOpen(false);
      setEditing(null);
      await load();
      toast.success(editing ? "Usuário atualizado" : "Usuário criado");
    } catch (err) {
      console.error(err);
      toast.error("Erro ao salvar usuário");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!(await confirm({ title: "Remover usuário?", variant: "danger" })))
      return;
    try {
      await apiRequest(`/users/${id}`, { method: "DELETE" });
      await load();
      toast.success("Usuário removido");
    } catch (err) {
      console.error(err);
      toast.error("Erro ao remover usuário");
    }
  };

  if (loading)
    return (
      <div className="flex items-center justify-center h-64 text-muted-foreground">
        Carregando...
      </div>
    );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-display text-2xl font-semibold text-foreground">
            Usuários
          </h1>
          <p className="text-sm text-muted-foreground">
            Gerenciar usuários da organização
          </p>
        </div>
        {canEdit && (
          <Button
            onClick={() => {
              setEditing(null);
              setForm({
                name: "",
                email: "",
                password: "",
                role_id: "",
                status: "active",
                farm_ids: [],
              });
              setDialogOpen(true);
            }}
          >
            <Plus className="h-4 w-4" /> Novo Usuário
          </Button>
        )}
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Nome</TableHeader>
            <TableHeader>Email</TableHeader>
            <TableHeader>Perfil</TableHeader>
            <TableHeader>Plano</TableHeader>
            <TableHeader>Gerenciado Por</TableHeader>
            <TableHeader>Status</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {users.map((u) => (
            <TableRow key={u.id}>
              <TableCell className="font-medium">
                <div className="flex items-center gap-2">
                  <User className="h-4 w-4 text-primary" />
                  {u.name}
                </div>
              </TableCell>
              <TableCell>{u.email}</TableCell>
              <TableCell>
                <Badge>{roleLabel(roles, u.role)}</Badge>
              </TableCell>
              <TableCell>{u.plan_name}</TableCell>
              <TableCell>{u.managed_by_name}</TableCell>
              <TableCell>
                <Badge variant={u.status === "active" ? "success" : "default"}>
                  {u.status === "active" ? "Ativo" : "Inativo"}
                </Badge>
              </TableCell>
              <TableCell className="text-right">
                {canEdit && (
                  <div className="flex justify-end gap-1">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => {
                        setEditing(u);
                        setForm({
                          name: u.name,
                          email: u.email,
                          password: "",
                          role_id: u.role_id,
                          status: u.status,
                          farm_ids: [],
                        });
                        setDialogOpen(true);
                      }}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleDelete(u.id)}
                    >
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                )}
              </TableCell>
            </TableRow>
          ))}
          {users.length === 0 && (
            <TableRow>
              <TableCell
                colSpan={8}
                className="text-center text-muted-foreground py-8"
              >
                Nenhum usuário cadastrado.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>

      <Dialog
        open={dialogOpen}
        onClose={() => {
          setDialogOpen(false);
          setEditing(null);
        }}
        title={editing ? "Editar Usuário" : "Novo Usuário"}
      >
        <form onSubmit={handleSave} className="space-y-4">
          <Field label="Nome" required>
            <Input
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              required
            />
          </Field>
          <Field label="Email" required>
            <Input
              type="email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              required
            />
          </Field>
          {!editing && (
            <Field label="Senha" required>
              <Input
                type="password"
                value={form.password}
                onChange={(e) => setForm({ ...form, password: e.target.value })}
                required
              />
            </Field>
          )}
          <Field label="Perfil">
            <Select
              value={form.role_id}
              onChange={(e) => setForm({ ...form, role_id: e.target.value })}
            >
              <option value="">Selecione...</option>
              {roles.map((role) => (
                <option key={role.id} value={role.id}>
                  {role.name}
                </option>
              ))}
            </Select>
          </Field>
          {editing && (
            <Field label="Status">
              <Select
                value={form.status}
                onChange={(e) => setForm({ ...form, status: e.target.value })}
              >
                <option value="active">Ativo</option>
                <option value="inactive">Inativo</option>
              </Select>
            </Field>
          )}
          {!editing && farms.length > 0 && (
            <fieldset className="space-y-2">
              <legend className="text-sm font-medium text-foreground">
                Vincular Fazendas
              </legend>
              <p className="text-xs text-muted-foreground">
                Selecione as fazendas que este usuário poderá acessar
              </p>
              <div className="space-y-1 max-h-48 overflow-y-auto border rounded-md p-2">
                {farms.map((farm) => (
                  <label
                    key={farm.id}
                    className="flex items-center gap-2 text-sm cursor-pointer hover:bg-muted rounded px-2 py-1"
                  >
                    <input
                      type="checkbox"
                      checked={form.farm_ids.includes(farm.id)}
                      onChange={(e) => {
                        if (e.target.checked) {
                          setForm({
                            ...form,
                            farm_ids: [...form.farm_ids, farm.id],
                          });
                        } else {
                          setForm({
                            ...form,
                            farm_ids: form.farm_ids.filter(
                              (id) => id !== farm.id,
                            ),
                          });
                        }
                      }}
                      className="h-4 w-4 rounded border-gray-300"
                    />
                    {farm.name}
                  </label>
                ))}
              </div>
            </fieldset>
          )}
          <div className="flex items-center justify-between gap-3 pt-2">
            <RequiredLegend />
            <div className="flex gap-3">
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  setDialogOpen(false);
                  setEditing(null);
                }}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={saving}>
                {saving ? "Salvando..." : "Salvar"}
              </Button>
            </div>
          </div>
        </form>
      </Dialog>
    </div>
  );
}
