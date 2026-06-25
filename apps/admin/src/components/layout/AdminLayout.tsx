import { useState } from 'react'
import { Outlet, Navigate, useLocation } from 'react-router-dom'
import { useAuth } from '@/lib/auth'
import { Sidebar } from './Sidebar'
import { Header } from './Header'

const adminRoutes = ['/tenants', '/users']

export function AdminLayout() {
  const { isAuthenticated, user } = useAuth()
  const location = useLocation()
  const [sidebarOpen, setSidebarOpen] = useState(false)

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  const isAdminRoute = adminRoutes.some((route) => location.pathname.startsWith(route))
  if (isAdminRoute && user?.role !== 'platform_owner') {
    return <Navigate to="/" replace />
  }

  return (
    <div className="flex h-screen bg-gray-50">
      <Sidebar open={sidebarOpen} onClose={() => setSidebarOpen(false)} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <Header onMenuClick={() => setSidebarOpen(true)} />
        <main className="flex-1 overflow-auto p-4 lg:p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
