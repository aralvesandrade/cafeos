import { Outlet, Navigate } from 'react-router-dom'
import { useAuth } from '@/lib/auth'

export function AuthLayout() {
  const { isAuthenticated } = useAuth()

  if (isAuthenticated) {
    return <Navigate to="/" replace />
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-coffee-beige to-white flex items-center justify-center p-4">
      <Outlet />
    </div>
  )
}
