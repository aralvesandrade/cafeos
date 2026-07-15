import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { usePermissions } from '@/lib/permissions'
import { useModules, ACCESS_LABELS, type Module, type AccessLevel } from '@/lib/permissions'
import { useRoles } from '@/lib/roles'
import { Button } from '@/components/ui/button'
import { Select } from '@/components/ui/select'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Save } from 'lucide-react'
import { Link } from 'react-router-dom'

interface RolePermission {
  role_id: string
  module: Module
  access: AccessLevel
}

type Matrix = Record<string, Record<string, AccessLevel>> // module -> role_id -> access

export function Permissions() {
  const { modules, loading: modulesLoading } = useModules()
  const { roles, loading: rolesLoading } = useRoles()
  const [matrix, setMatrix] = useState<Matrix>({})
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const toast = useToast()
  const { refresh } = usePermissions()

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const rows = await apiRequest<RolePermission[]>('/permissions')
      const next: Matrix = {}
      for (const module of modules) next[module.key] = {}
      for (const row of rows) {
        if (!next[row.module]) next[row.module] = {}
        next[row.module][row.role_id] = row.access
      }
      setMatrix(next)
    } catch (err) {
      console.error(err)
      toast.error('Erro ao carregar permissões')
    } finally {
      setLoading(false)
    }
  }, [toast, modules])

  useEffect(() => {
    if (modules.length > 0) load()
  }, [load, modules])

  const cellAccess = (module: Module, roleId: string): AccessLevel => matrix[module]?.[roleId] ?? 'none'

  const setCell = (module: Module, roleId: string, access: AccessLevel) => {
    setMatrix((prev) => ({ ...prev, [module]: { ...prev[module], [roleId]: access } }))
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      const entries: RolePermission[] = modules.flatMap((module) =>
        roles.map((role) => ({ role_id: role.id, module: module.key, access: cellAccess(module.key, role.id) }))
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

  if (loading || modulesLoading || rolesLoading) {
    return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="font-display text-2xl font-semibold text-foreground">Permissões</h1>
          <p className="text-sm text-muted-foreground">
            Defina o que cada perfil pode ver e editar nesta organização.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Link to="/roles">
            <Button variant="outline">Gerenciar Papéis</Button>
          </Link>
          <Button onClick={handleSave} disabled={saving}>
            <Save className="h-4 w-4" /> {saving ? 'Salvando...' : 'Salvar alterações'}
          </Button>
        </div>
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Módulo</TableHeader>
            {roles.map((role) => (
              <TableHeader key={role.id}>{role.name}</TableHeader>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {modules.map((module) => (
            <TableRow key={module.key}>
              <TableCell className="font-medium whitespace-nowrap">{module.name}</TableCell>
              {roles.map((role) => (
                <TableCell key={role.id}>
                  <Select
                    value={cellAccess(module.key, role.id)}
                    onChange={(e) => setCell(module.key, role.id, e.target.value as AccessLevel)}
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
