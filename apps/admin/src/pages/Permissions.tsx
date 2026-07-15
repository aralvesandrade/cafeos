import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { usePermissions } from '@/lib/permissions'
import { ALL_MODULES, MODULE_LABELS, ACCESS_LABELS, type Module, type AccessLevel } from '@/lib/permissions'
import { ALL_ROLES, ROLE_LABELS, type UserRole } from '@/lib/roles'
import { Button } from '@/components/ui/button'
import { Select } from '@/components/ui/select'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Save } from 'lucide-react'

interface RolePermission {
  role: UserRole
  module: Module
  access: AccessLevel
}

type Matrix = Record<Module, Record<UserRole, AccessLevel>>

function emptyMatrix(): Matrix {
  const matrix = {} as Matrix
  for (const module of ALL_MODULES) {
    matrix[module] = {} as Record<UserRole, AccessLevel>
    for (const role of ALL_ROLES) {
      matrix[module][role] = 'none'
    }
  }
  return matrix
}

export function Permissions() {
  const [matrix, setMatrix] = useState<Matrix>(emptyMatrix)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const toast = useToast()
  const { refresh } = usePermissions()

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const rows = await apiRequest<RolePermission[]>('/permissions')
      const next = emptyMatrix()
      for (const row of rows) {
        if (next[row.module]) next[row.module][row.role] = row.access
      }
      setMatrix(next)
    } catch (err) {
      console.error(err)
      toast.error('Erro ao carregar permissões')
    } finally {
      setLoading(false)
    }
  }, [toast])

  useEffect(() => { load() }, [load])

  const setCell = (module: Module, role: UserRole, access: AccessLevel) => {
    setMatrix((prev) => ({ ...prev, [module]: { ...prev[module], [role]: access } }))
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      const entries: RolePermission[] = ALL_MODULES.flatMap((module) =>
        ALL_ROLES.map((role) => ({ role, module, access: matrix[module][role] }))
      )
      await apiRequest('/permissions', { method: 'PUT', body: entries })
      await refresh()
      toast.success('Permissões salvas')
    } catch (err) {
      console.error(err)
      toast.error('Erro ao salvar permissões')
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="font-display text-2xl font-semibold text-foreground">Permissões</h1>
          <p className="text-sm text-muted-foreground">
            Defina o que cada perfil pode ver e editar nesta organização.
          </p>
        </div>
        <Button onClick={handleSave} disabled={saving}>
          <Save className="h-4 w-4" /> {saving ? 'Salvando...' : 'Salvar alterações'}
        </Button>
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Módulo</TableHeader>
            {ALL_ROLES.map((role) => (
              <TableHeader key={role}>{ROLE_LABELS[role]}</TableHeader>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {ALL_MODULES.map((module) => (
            <TableRow key={module}>
              <TableCell className="font-medium whitespace-nowrap">{MODULE_LABELS[module]}</TableCell>
              {ALL_ROLES.map((role) => (
                <TableCell key={role}>
                  <Select
                    value={matrix[module][role]}
                    onChange={(e) => setCell(module, role, e.target.value as AccessLevel)}
                    className="w-32"
                  >
                    {(['none', 'read', 'write'] as AccessLevel[]).map((level) => (
                      <option key={level} value={level}>{ACCESS_LABELS[level]}</option>
                    ))}
                  </Select>
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
