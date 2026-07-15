import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { useConfirm } from '@/lib/confirm'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Select } from '@/components/ui/select'
import { Plus, Pencil, Trash2, Package, ArrowUpDown } from 'lucide-react'

interface StockItem {
  id: string
  product_id: string
  product: { id: string; name: string }
  farm_id: string | null
  quantity: number
  unit: string
  batch: string
  expiry_date: string
  min_stock: number
  location: string
  notes: string
}

interface Product {
  id: string
  name: string
}

interface Farm { id: string; name: string }

export function Stock() {
  const [items, setItems] = useState<StockItem[]>([])
  const [products, setProducts] = useState<Product[]>([])
  const [farms, setFarms] = useState<Farm[]>([])
  const [farmFilter, setFarmFilter] = useState('')
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [movDialogOpen, setMovDialogOpen] = useState(false)
  const [selectedItem, setSelectedItem] = useState<string>('')
  const [editing, setEditing] = useState<StockItem | null>(null)
  const [form, setForm] = useState({ product_id: '', farm_id: '', quantity: '', unit: '', batch: '', expiry_date: '', min_stock: '', location: '', notes: '' })
  const [movForm, setMovForm] = useState({ item_id: '', type: 'in', quantity: '', date: '', reference: '', notes: '' })
  const [saving, setSaving] = useState(false)
  const toast = useToast()
  const confirm = useConfirm()

  const load = useCallback(async () => {
    try {
      const [itemData, prodData, farmsData] = await Promise.all([
        apiRequest<StockItem[]>('/stock/items', { params: farmFilter ? { farm_id: farmFilter } : undefined }),
        apiRequest<Product[]>('/agricultural-products'),
        apiRequest<Farm[]>('/farms'),
      ])
      setItems(itemData); setProducts(prodData); setFarms(farmsData)
    } catch (err) { console.error(err) }
    finally { setLoading(false) }
  }, [farmFilter])

  useEffect(() => { load() }, [load])

  const handleSave = async () => {
    setSaving(true)
    try {
      const body = { ...form, farm_id: form.farm_id || null, quantity: parseFloat(form.quantity), min_stock: parseFloat(form.min_stock) }
      if (editing) await apiRequest(`/stock/items/${editing.id}`, { method: 'PUT', body })
      else await apiRequest('/stock/items', { method: 'POST', body })
      setDialogOpen(false); setEditing(null); await load()
      toast.success(editing ? 'Item atualizado' : 'Item cadastrado')
    } catch (err) { console.error(err); toast.error('Erro ao salvar item') }
    finally { setSaving(false) }
  }

  const handleMovement = async () => {
    setSaving(true)
    try {
      await apiRequest('/stock/movements', { method: 'POST', body: { ...movForm, quantity: parseFloat(movForm.quantity) } })
      setMovDialogOpen(false); await load()
      toast.success('Movimentação registrada')
    } catch (err) { console.error(err); toast.error('Erro ao registrar movimentação') }
    finally { setSaving(false) }
  }

  const handleDeleteItem = async (id: string) => {
    if (!(await confirm({ title: 'Remover item?', variant: 'danger' }))) return
    try {
      await apiRequest(`/stock/items/${id}`, { method: 'DELETE' })
      await load()
      toast.success('Item removido')
    } catch (err) { console.error(err); toast.error('Erro ao remover item') }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (<div className="space-y-6">
    <div className="flex items-center justify-between flex-wrap gap-4">
      <div><h1 className="font-display text-2xl font-semibold text-foreground">Estoque</h1><p className="text-sm text-muted-foreground">Insumos e produtos armazenados</p></div>
      <div className="flex gap-2 flex-wrap">
        <Select value={farmFilter} onChange={(e) => setFarmFilter(e.target.value)} className="w-48">
          <option value="">Todas as fazendas</option>
          {farms.map((f) => (
            <option key={f.id} value={f.id}>{f.name}</option>
          ))}
        </Select>
        <Button variant="outline" onClick={() => { setMovForm({ item_id: '', type: 'in', quantity: '', date: '', reference: '', notes: '' }); setMovDialogOpen(true) }}><ArrowUpDown className="h-4 w-4" /> Movimentar</Button>
        <Button onClick={() => { setEditing(null); setForm({ product_id: '', farm_id: '', quantity: '', unit: '', batch: '', expiry_date: '', min_stock: '', location: '', notes: '' }); setDialogOpen(true) }}><Plus className="h-4 w-4" /> Novo Item</Button>
      </div>
    </div>
    <Table>
      <TableHead><TableRow><TableHeader>Produto</TableHeader><TableHeader>Fazenda</TableHeader><TableHeader>Qtd</TableHeader><TableHeader>Un</TableHeader><TableHeader>Lote</TableHeader><TableHeader>Vencimento</TableHeader><TableHeader>Estoque Mín</TableHeader><TableHeader>Local</TableHeader><TableHeader className="text-right">Ações</TableHeader></TableRow></TableHead>
      <TableBody>{items.map((i) => {
        const prodName = i.product?.name || i.product_id
        const farmName = farms.find(f => f.id === i.farm_id)?.name
        const lowStock = i.quantity <= i.min_stock
        return (<TableRow key={i.id}>
          <TableCell className="font-medium"><div className="flex items-center gap-2"><Package className="h-4 w-4 text-primary" />{prodName}</div></TableCell>
          <TableCell className="text-muted-foreground text-sm">{farmName || '-'}</TableCell>
          <TableCell className={lowStock ? 'text-destructive font-medium' : ''}>{i.quantity}</TableCell>
          <TableCell>{i.unit}</TableCell>
          <TableCell className="text-muted-foreground text-sm">{i.batch || '-'}</TableCell>
          <TableCell className="text-sm">{i.expiry_date ? new Date(i.expiry_date).toLocaleDateString() : '-'}</TableCell>
          <TableCell>{i.min_stock}</TableCell>
          <TableCell className="text-sm">{i.location || '-'}</TableCell>
          <TableCell className="text-right"><div className="flex justify-end gap-1">
            <Button variant="ghost" size="sm" onClick={() => { setEditing(i); setForm({ product_id: i.product_id, farm_id: i.farm_id || '', quantity: String(i.quantity), unit: i.unit, batch: i.batch, expiry_date: i.expiry_date || '', min_stock: String(i.min_stock), location: i.location, notes: i.notes }); setDialogOpen(true) }}><Pencil className="h-4 w-4" /></Button>
            <Button variant="ghost" size="sm" onClick={() => handleDeleteItem(i.id)}><Trash2 className="h-4 w-4 text-destructive" /></Button>
          </div></TableCell>
        </TableRow>)
      })}
      {items.length === 0 && <TableRow><TableCell colSpan={9} className="text-center text-muted-foreground py-8">Nenhum item em estoque.</TableCell></TableRow>}
      </TableBody>
    </Table>
    <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Item' : 'Novo Item'}>
      <div className="space-y-4">
        <div><label className="block text-sm font-medium text-foreground mb-1">Produto</label>
          <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={form.product_id} onChange={(e) => setForm({ ...form, product_id: e.target.value })}>
            <option value="">Selecione</option>{products.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
          </select></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Fazenda</label>
          <Select value={form.farm_id} onChange={(e) => setForm({ ...form, farm_id: e.target.value })}>
            <option value="">Sem fazenda vinculada</option>
            {farms.map((f) => <option key={f.id} value={f.id}>{f.name}</option>)}
          </Select></div>
        <div className="grid grid-cols-2 gap-3">
          <div><label className="block text-sm font-medium text-foreground mb-1">Quantidade</label><Input type="number" step="0.01" value={form.quantity} onChange={(e) => setForm({ ...form, quantity: e.target.value })} /></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Unidade</label><Input value={form.unit} onChange={(e) => setForm({ ...form, unit: e.target.value })} placeholder="kg, l, un" /></div>
        </div>
        <div className="grid grid-cols-2 gap-3">
          <div><label className="block text-sm font-medium text-foreground mb-1">Lote</label><Input value={form.batch} onChange={(e) => setForm({ ...form, batch: e.target.value })} /></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Vencimento</label><Input type="date" value={form.expiry_date} onChange={(e) => setForm({ ...form, expiry_date: e.target.value })} /></div>
        </div>
        <div className="grid grid-cols-2 gap-3">
          <div><label className="block text-sm font-medium text-foreground mb-1">Estoque Mínimo</label><Input type="number" step="0.01" value={form.min_stock} onChange={(e) => setForm({ ...form, min_stock: e.target.value })} /></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Local</label><Input value={form.location} onChange={(e) => setForm({ ...form, location: e.target.value })} /></div>
        </div>
        <div className="flex justify-end gap-3 pt-2">
          <Button variant="outline" onClick={() => { setDialogOpen(false); setEditing(null) }}>Cancelar</Button>
          <Button onClick={handleSave} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
        </div>
      </div>
    </Dialog>
    <Dialog open={movDialogOpen} onClose={() => setMovDialogOpen(false)} title="Movimentar Estoque">
      <div className="space-y-4">
        <div><label className="block text-sm font-medium text-foreground mb-1">Item</label>
          <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={movForm.item_id} onChange={(e) => setMovForm({ ...movForm, item_id: e.target.value })}>
            <option value="">Selecione</option>{items.map((i) => <option key={i.id} value={i.id}>{i.product?.name || i.product_id} ({i.quantity} {i.unit})</option>)}
          </select></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Tipo</label>
          <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={movForm.type} onChange={(e) => setMovForm({ ...movForm, type: e.target.value })}>
            <option value="in">Entrada</option><option value="out">Saída</option>
          </select></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Quantidade</label><Input type="number" step="0.01" value={movForm.quantity} onChange={(e) => setMovForm({ ...movForm, quantity: e.target.value })} /></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Data</label><Input type="date" value={movForm.date} onChange={(e) => setMovForm({ ...movForm, date: e.target.value })} /></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Referência</label><Input value={movForm.reference} onChange={(e) => setMovForm({ ...movForm, reference: e.target.value })} placeholder="Nota fiscal, operação..." /></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Observações</label><Input value={movForm.notes} onChange={(e) => setMovForm({ ...movForm, notes: e.target.value })} /></div>
        <div className="flex justify-end gap-3 pt-2">
          <Button variant="outline" onClick={() => setMovDialogOpen(false)}>Cancelar</Button>
          <Button onClick={handleMovement} disabled={saving}>{saving ? 'Salvando...' : 'Registrar'}</Button>
        </div>
      </div>
    </Dialog>
  </div>)
}
