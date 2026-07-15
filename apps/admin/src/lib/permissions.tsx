import { createContext, useContext, useState, useCallback, useEffect, type ReactNode } from 'react'
import { apiRequest } from './api'
import { useAuth } from './auth'

export type Module = 'dashboard' | 'farms' | 'operations' | 'harvests' | 'resources' | 'financial' | 'users' | 'permissions'
export type AccessLevel = 'none' | 'read' | 'write'

export interface ModuleMeta {
  key: Module
  name: string
  order: number
}

// useModules fetches the module catalog (name/order metadata) from the
// backend — the set of keys is still fixed (routes are wired to them in
// Go), but name/order are now editable through the Permissions screen
// instead of duplicated here as a hardcoded label map.
export function useModules() {
  const [modules, setModules] = useState<ModuleMeta[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let cancelled = false
    apiRequest<ModuleMeta[]>('/modules')
      .then((data) => { if (!cancelled) setModules(data) })
      .catch((err) => console.error(err))
      .finally(() => { if (!cancelled) setLoading(false) })
    return () => { cancelled = true }
  }, [])

  return { modules, loading }
}

export const ACCESS_LABELS: Record<AccessLevel, string> = {
  none: 'Nenhum',
  read: 'Leitura',
  write: 'Escrita',
}

interface PermissionsContextType {
  access: Record<Module, AccessLevel> | null
  loading: boolean
  refresh: () => Promise<void>
}

const PermissionsContext = createContext<PermissionsContextType | null>(null)

export function PermissionsProvider({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth()
  const [access, setAccess] = useState<Record<Module, AccessLevel> | null>(null)
  const [loading, setLoading] = useState(true)

  const refresh = useCallback(async () => {
    if (!isAuthenticated) {
      setAccess(null)
      setLoading(false)
      return
    }
    setLoading(true)
    try {
      const data = await apiRequest<Record<Module, AccessLevel>>('/permissions/me')
      setAccess(data)
    } catch (err) {
      console.error(err)
      setAccess(null)
    } finally {
      setLoading(false)
    }
  }, [isAuthenticated])

  useEffect(() => { refresh() }, [refresh])

  return (
    <PermissionsContext.Provider value={{ access, loading, refresh }}>
      {children}
    </PermissionsContext.Provider>
  )
}

export function usePermissions() {
  const ctx = useContext(PermissionsContext)
  if (!ctx) throw new Error('usePermissions must be used within PermissionsProvider')
  return ctx
}

export function useModuleAccess(module: Module): AccessLevel {
  const { access } = usePermissions()
  return access?.[module] ?? 'none'
}

export const ROUTE_MODULE: Record<string, Module> = {
  '/': 'dashboard',
  '/farms': 'farms',
  '/plots': 'farms',
  '/operation-types': 'farms',
  '/operations': 'operations',
  '/harvests': 'harvests',
  '/financial': 'financial',
  '/cost-centers': 'financial',
  '/budgets': 'financial',
  '/cost-allocations': 'financial',
  '/stock': 'resources',
  '/fleet': 'resources',
  '/labor': 'resources',
  '/team': 'users',
  '/permissions': 'permissions',
  '/roles': 'permissions',
}

export function moduleForPath(pathname: string): Module | null {
  if (ROUTE_MODULE[pathname]) return ROUTE_MODULE[pathname]
  const prefix = Object.keys(ROUTE_MODULE)
    .filter((path) => path !== '/' && pathname.startsWith(path))
    .sort((a, b) => b.length - a.length)[0]
  return prefix ? ROUTE_MODULE[prefix] : null
}
