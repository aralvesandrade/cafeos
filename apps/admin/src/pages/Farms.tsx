import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { FarmList } from '@/components/farms/FarmList'
import { FarmForm } from '@/components/farms/FarmForm'
import { Plus } from 'lucide-react'

interface Farm {
  id: string
  name: string
  owner: string
  total_area: number
  planted_area: number
  location: string
}

export function Farms() {
  const [farms, setFarms] = useState<Farm[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingFarm, setEditingFarm] = useState<Farm | null>(null)
  const [saving, setSaving] = useState(false)

  const loadFarms = useCallback(async () => {
    try {
      const data = await apiRequest<Farm[]>('/farms')
      setFarms(data)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadFarms()
  }, [loadFarms])

  const handleSave = async (formData: { name: string; owner: string; location: string; total_area: string; planted_area: string }) => {
    setSaving(true)
    try {
      const payload = {
        name: formData.name,
        owner: formData.owner,
        location: formData.location,
        total_area: parseFloat(formData.total_area),
        planted_area: parseFloat(formData.planted_area),
      }

      if (editingFarm) {
        await apiRequest(`/farms/${editingFarm.id}`, { method: 'PUT', body: payload })
      } else {
        await apiRequest('/farms', { method: 'POST', body: payload })
      }

      setDialogOpen(false)
      setEditingFarm(null)
      await loadFarms()
    } catch (err) {
      console.error(err)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Remover fazenda?')) return
    try {
      await apiRequest(`/farms/${id}`, { method: 'DELETE' })
      await loadFarms()
    } catch (err) {
      console.error(err)
    }
  }

  const openEdit = (farm: Farm) => {
    setEditingFarm(farm)
    setDialogOpen(true)
  }

  const openCreate = () => {
    setEditingFarm(null)
    setDialogOpen(true)
  }

  if (loading) {
    return <div className="flex items-center justify-center h-64 text-coffee-text-light">Carregando...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-coffee-green-dark">Fazendas</h1>
          <p className="text-sm text-coffee-text-light">Gerencie suas propriedades rurais</p>
        </div>
        <Button onClick={openCreate}>
          <Plus className="h-4 w-4" />
          Nova Fazenda
        </Button>
      </div>

      <FarmList farms={farms} onEdit={openEdit} onDelete={handleDelete} />

      <Dialog
        open={dialogOpen}
        onClose={() => { setDialogOpen(false); setEditingFarm(null) }}
        title={editingFarm ? 'Editar Fazenda' : 'Nova Fazenda'}
      >
        <FarmForm
          initial={editingFarm ? {
            name: editingFarm.name,
            owner: editingFarm.owner,
            location: editingFarm.location,
            total_area: String(editingFarm.total_area),
            planted_area: String(editingFarm.planted_area),
          } : undefined}
          onSave={handleSave}
          onCancel={() => { setDialogOpen(false); setEditingFarm(null) }}
          loading={saving}
        />
      </Dialog>
    </div>
  )
}
