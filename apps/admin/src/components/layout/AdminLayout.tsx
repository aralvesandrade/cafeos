import { useState } from 'react'
import { Outlet, Navigate, useLocation } from 'react-router-dom'
import { useAuth } from '@/lib/auth'
import { moduleForPath, useModuleAccess } from '@/lib/permissions'
import { Forbidden } from '@/pages/Forbidden'
import { Sidebar } from './Sidebar'
import { Header } from './Header'

const platformOnlyRoutes = ['/organizations', '/users']

export function AdminLayout() {
  const { isAuthenticated, user } = useAuth()
  const location = useLocation()
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const isPlatformRoute = platformOnlyRoutes.some((route) => location.pathname.startsWith(route))
  const module = moduleForPath(location.pathname)
  const moduleAccess = useModuleAccess(module ?? 'dashboard')

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (isPlatformRoute && user?.role !== 'platform_owner') {
    return <Navigate to="/" replace />
  }

  const forbidden = !isPlatformRoute && module !== null && moduleAccess === 'none'

  return (
    <div className="flex h-screen bg-background">
      <Sidebar open={sidebarOpen} onClose={() => setSidebarOpen(false)} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <Header onMenuClick={() => setSidebarOpen(true)} />
        <main className="flex-1 overflow-auto p-4 lg:p-6">
          {forbidden ? <Forbidden /> : <Outlet />}
        </main>
      </div>
    </div>
  )
}
