import { useEffect, useState, useCallback } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { apiRequest } from "@/lib/api";
import { useToast } from "@/lib/toast";
import { useConfirm } from "@/lib/confirm";
import {
  Table,
  TableHead,
  TableBody,
  TableRow,
  TableHeader,
  TableCell,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Select } from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Dialog } from "@/components/ui/dialog";
import {
  Search,
  Tractor,
  Plus,
  Pencil,
  Trash2,
  ExternalLink,
} from "lucide-react";
import { useModuleAccess } from "@/lib/permissions";

interface Operation {
  id: string;
  type_id: string;
  type_name: string;
  type_color: "info" | "warning" | "success" | "default" | "danger";
  date: string;
  plot_id: string;
  plot_name: string;
  responsible: string;
  product_used: string;
  quantity: number;
  cost: number;
  notes: string;
}

interface Farm {
  id: string;
  name: string;
}
interface Plot {
  id: string;
  name: string;
  farm_id: string;
}
interface OperationType {
  id: string;
  name: string;
  code: string;
  color: "info" | "warning" | "success" | "default" | "danger";
}

const emptyForm = {
  plot_id: "",
  type_id: "",
  date: "",
  responsible: "",
  product_used: "",
  quantity: "",
  cost: "",
  notes: "",
};

export function Operations() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [operations, setOperations] = useState<Operation[]>([]);
  const [farms, setFarms] = useState<Farm[]>([]);
  const [plots, setPlots] = useState<Plot[]>([]);
  const [operationTypes, setOperationTypes] = useState<OperationType[]>([]);
  const [farmFilter, setFarmFilter] = useState("");
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editing, setEditing] = useState<Operation | null>(null);
  const [form, setForm] = useState(emptyForm);
  const [saving, setSaving] = useState(false);
  const toast = useToast();
  const confirm = useConfirm();
  const canEdit = useModuleAccess("operations") === "write";

  const load = useCallback(async () => {
    try {
      const [opsData, farmsData, plotsData, typesData] = await Promise.all([
        apiRequest<Operation[]>("/operations", {
          params: farmFilter ? { farm_id: farmFilter } : undefined,
        }),
        apiRequest<Farm[]>("/farms"),
        apiRequest<Plot[]>("/plots"),
        apiRequest<OperationType[]>("/operation-types"),
      ]);
      setOperations(opsData);
      setFarms(farmsData);
      setPlots(plotsData);
      setOperationTypes(typesData);
      return opsData;
    } catch (err) {
      console.error(err);
      return [];
    } finally {
      setLoading(false);
    }
  }, [farmFilter]);

  const openEdit = useCallback((op: Operation) => {
    setEditing(op);
    setForm({
      plot_id: op.plot_id,
      type_id: op.type_id,
      date: op.date ? op.date.slice(0, 10) : "",
      responsible: op.responsible,
      product_used: op.product_used,
      quantity: String(op.quantity),
      cost: String(op.cost),
      notes: op.notes,
    });
    setDialogOpen(true);
  }, []);

  useEffect(() => {
    load().then((opsData) => {
      const editId = searchParams.get("edit");
      if (editId) {
        const op = opsData.find((o) => o.id === editId);
        if (op) openEdit(op);
        setSearchParams(
          (params) => {
            params.delete("edit");
            return params;
          },
          { replace: true },
        );
      }
    });
  }, [load, searchParams, setSearchParams, openEdit]);

  const handleSave = async () => {
    setSaving(true);
    try {
      const body = {
        ...form,
        quantity: parseFloat(form.quantity) || 0,
        cost: parseFloat(form.cost) || 0,
      };
      if (editing)
        await apiRequest(`/operations/${editing.id}`, { method: "PUT", body });
      else await apiRequest("/operations", { method: "POST", body });
      setDialogOpen(false);
      setEditing(null);
      await load();
      toast.success(editing ? "Operação atualizada" : "Operação registrada");
    } catch (err) {
      console.error(err);
      toast.error("Erro ao salvar operação");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!(await confirm({ title: "Remover operação?", variant: "danger" }))) return;
    try {
      await apiRequest(`/operations/${id}`, { method: "DELETE" });
      await load();
      toast.success("Operação removida");
    } catch (err) {
      console.error(err);
      toast.error("Erro ao remover operação");
    }
  };

  const filtered = operations.filter(
    (op) =>
      (op.plot_name || op.plot_id || "")
        .toLowerCase()
        .includes(search.toLowerCase()) ||
      (op.type_name || "").toLowerCase().includes(search.toLowerCase()) ||
      (op.responsible || "").toLowerCase().includes(search.toLowerCase()),
  );

  if (loading)
    return (
      <div className="flex items-center justify-center h-64 text-muted-foreground">
        Carregando...
      </div>
    );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="font-display text-2xl font-semibold text-foreground">Operações</h1>
          <p className="text-sm text-muted-foreground">
            Histórico de operações agrícolas
          </p>
        </div>
        {canEdit && (
          <Button
            onClick={() => {
              setEditing(null);
              setForm(emptyForm);
              setDialogOpen(true);
            }}
          >
            <Plus className="h-4 w-4" /> Nova Operação
          </Button>
        )}
      </div>

      <div className="flex flex-wrap gap-3">
        <div className="relative w-full max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            className="pl-9"
            placeholder="Buscar por talhão, tipo, responsável..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
        <Select
          value={farmFilter}
          onChange={(e) => setFarmFilter(e.target.value)}
          className="w-48"
        >
          <option value="">Todas as fazendas</option>
          {farms.map((f) => (
            <option key={f.id} value={f.id}>
              {f.name}
            </option>
          ))}
        </Select>
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Data</TableHeader>
            <TableHeader>Tipo</TableHeader>
            <TableHeader>Talhão</TableHeader>
            <TableHeader>Responsável</TableHeader>
            <TableHeader>Produto</TableHeader>
            <TableHeader>Quantidade</TableHeader>
            <TableHeader className="text-right">Custo</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {filtered.map((op) => (
            <TableRow key={op.id}>
              <TableCell>{new Date(op.date).toLocaleDateString()}</TableCell>
              <TableCell>
                <div className="flex items-center gap-2">
                  <Tractor className="h-4 w-4 text-primary" />
                  <Badge variant={op.type_color || "default"}>
                    {op.type_name || "—"}
                  </Badge>
                </div>
              </TableCell>
              <TableCell>{op.plot_name || op.plot_id.slice(0, 8)}</TableCell>
              <TableCell>{op.responsible}</TableCell>
              <TableCell>{op.product_used}</TableCell>
              <TableCell>{op.quantity}</TableCell>
              <TableCell className="text-right">
                R$ {op.cost.toFixed(2)}
              </TableCell>
              <TableCell className="text-right">
                <div className="flex justify-end gap-1">
                  <Button variant="ghost" size="sm" asChild>
                    <Link to={`/operations/${op.id}`}>
                      <ExternalLink className="h-4 w-4" />
                    </Link>
                  </Button>
                  {canEdit && (
                    <>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => openEdit(op)}
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDelete(op.id)}
                      >
                        <Trash2 className="h-4 w-4 text-destructive" />
                      </Button>
                    </>
                  )}
                </div>
              </TableCell>
            </TableRow>
          ))}
          {filtered.length === 0 && (
            <TableRow>
              <TableCell
                colSpan={8}
                className="text-center text-muted-foreground py-8"
              >
                Nenhuma operação encontrada.
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
        title={editing ? "Editar Operação" : "Nova Operação"}
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">
              Talhão
            </label>
            <select
              className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm"
              value={form.plot_id}
              onChange={(e) => setForm({ ...form, plot_id: e.target.value })}
            >
              <option value="">Selecione</option>
              {plots.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-sm font-medium text-foreground mb-1">
                Tipo
              </label>
              <select
                className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm"
                value={form.type_id}
                onChange={(e) => setForm({ ...form, type_id: e.target.value })}
              >
                <option value="">Selecione</option>
                {operationTypes.map((ot) => (
                  <option key={ot.id} value={ot.id}>
                    {ot.name}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-foreground mb-1">
                Data
              </label>
              <Input
                type="date"
                value={form.date}
                onChange={(e) => setForm({ ...form, date: e.target.value })}
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">
              Responsável
            </label>
            <Input
              value={form.responsible}
              onChange={(e) =>
                setForm({ ...form, responsible: e.target.value })
              }
            />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-sm font-medium text-foreground mb-1">
                Produto Usado
              </label>
              <Input
                value={form.product_used}
                onChange={(e) =>
                  setForm({ ...form, product_used: e.target.value })
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-foreground mb-1">
                Quantidade
              </label>
              <Input
                type="number"
                step="0.01"
                value={form.quantity}
                onChange={(e) => setForm({ ...form, quantity: e.target.value })}
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">
              Custo (R$)
            </label>
            <Input
              type="number"
              step="0.01"
              value={form.cost}
              onChange={(e) => setForm({ ...form, cost: e.target.value })}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">
              Observações
            </label>
            <Input
              value={form.notes}
              onChange={(e) => setForm({ ...form, notes: e.target.value })}
            />
          </div>
          <div className="flex justify-end gap-3 pt-2">
            <Button
              variant="outline"
              onClick={() => {
                setDialogOpen(false);
                setEditing(null);
              }}
            >
              Cancelar
            </Button>
            <Button onClick={handleSave} disabled={saving}>
              {saving ? "Salvando..." : "Salvar"}
            </Button>
          </div>
        </div>
      </Dialog>
    </div>
  );
}
