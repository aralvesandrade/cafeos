import { Outlet, Navigate } from 'react-router-dom'
import { Sun, Moon } from 'lucide-react'
import { useAuth } from '@/lib/auth'
import { useTheme } from '@/lib/theme'

export function AuthLayout() {
  const { isAuthenticated } = useAuth()
  const { theme, toggleTheme } = useTheme()

  if (isAuthenticated) {
    return <Navigate to="/" replace />
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-muted to-background flex items-center justify-center p-4 relative">
      <button
        onClick={toggleTheme}
        className="absolute top-4 right-4 inline-flex items-center gap-2 h-[34px] px-3 rounded-lg border border-border text-foreground text-sm font-medium hover:bg-muted bg-background"
      >
        {theme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
        <span className="hidden sm:inline">{theme === 'dark' ? 'Claro' : 'Escuro'}</span>
      </button>
      <Outlet />
    </div>
  )
}
