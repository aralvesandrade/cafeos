import { createContext, useContext, useState, useCallback, useEffect, type ReactNode } from 'react'
import { apiRequest } from './api'
import { useAuth } from './auth'

// Role is a global catalog entry — like Module, it's shared by every
// organization. What varies per organization is access (see
// lib/permissions.tsx), not the set of available roles.
export interface Role {
  id: string
  key: string
  name: string
  is_system: boolean
}

interface RolesContextType {
  roles: Role[]
  loading: boolean
  refresh: () => Promise<void>
}

const RolesContext = createContext<RolesContextType | null>(null)

export function RolesProvider({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth()
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(true)

  const refresh = useCallback(async () => {
    if (!isAuthenticated) {
      setRoles([])
      setLoading(false)
      return
    }
    setLoading(true)
    try {
      const data = await apiRequest<Role[]>('/roles')
      setRoles(data)
    } catch (err) {
      console.error(err)
      setRoles([])
    } finally {
      setLoading(false)
    }
  }, [isAuthenticated])

  useEffect(() => { refresh() }, [refresh])

  return (
    <RolesContext.Provider value={{ roles, loading, refresh }}>
      {children}
    </RolesContext.Provider>
  )
}

export function useRoles() {
  const ctx = useContext(RolesContext)
  if (!ctx) throw new Error('useRoles must be used within RolesProvider')
  return ctx
}

// roleLabel looks up a role's display name by key, falling back to the key
// itself for roles that were deleted after being assigned (or not yet
// loaded).
export function roleLabel(roles: Role[], key: string): string {
  return roles.find((r) => r.key === key)?.name ?? key
}
