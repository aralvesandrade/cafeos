import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

interface FarmData {
  name: string
  owner: string
  location: string
  total_area: string
  planted_area: string
}

interface FarmFormProps {
  initial?: FarmData
  onSave: (data: FarmData) => void
  onCancel: () => void
  loading?: boolean
}

export function FarmForm({ initial, onSave, onCancel, loading }: FarmFormProps) {
  const [form, setForm] = useState<FarmData>({
    name: '',
    owner: '',
    location: '',
    total_area: '',
    planted_area: '',
  })

  useEffect(() => {
    if (initial) setForm(initial)
  }, [initial])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSave(form)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-coffee-text mb-1">Nome</label>
        <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
      </div>
      <div>
        <label className="block text-sm font-medium text-coffee-text mb-1">Proprietário</label>
        <Input value={form.owner} onChange={(e) => setForm({ ...form, owner: e.target.value })} required />
      </div>
      <div>
        <label className="block text-sm font-medium text-coffee-text mb-1">Localização</label>
        <Input value={form.location} onChange={(e) => setForm({ ...form, location: e.target.value })} required />
      </div>
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-coffee-text mb-1">Área Total (ha)</label>
          <Input type="number" step="0.01" value={form.total_area} onChange={(e) => setForm({ ...form, total_area: e.target.value })} required />
        </div>
        <div>
          <label className="block text-sm font-medium text-coffee-text mb-1">Área Plantada (ha)</label>
          <Input type="number" step="0.01" value={form.planted_area} onChange={(e) => setForm({ ...form, planted_area: e.target.value })} required />
        </div>
      </div>
      <div className="flex justify-end gap-3 pt-2">
        <Button type="button" variant="outline" onClick={onCancel}>Cancelar</Button>
        <Button type="submit" disabled={loading}>
          {loading ? 'Salvando...' : 'Salvar'}
        </Button>
      </div>
    </form>
  )
}
