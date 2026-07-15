import { Menu, LogOut, User, Sun, Moon } from 'lucide-react'
import { useAuth } from '@/lib/auth'
import { useTheme } from '@/lib/theme'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { NotificationBell } from '@/components/layout/NotificationBell'

interface HeaderProps {
  onMenuClick: () => void
}

export function Header({ onMenuClick }: HeaderProps) {
  const { user, logout } = useAuth()
  const { theme, toggleTheme } = useTheme()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <header className="h-16 bg-background border-b border-border flex items-center justify-between px-4 lg:px-6">
      <button
        onClick={onMenuClick}
        className="lg:hidden text-foreground hover:text-primary"
      >
        <Menu className="h-6 w-6" />
      </button>

      <div className="hidden lg:block" />

      <div className="flex items-center gap-4">
        <NotificationBell />
        <button
          onClick={toggleTheme}
          className="inline-flex items-center gap-2 h-[34px] px-3 rounded-lg border border-border text-foreground text-sm font-medium hover:bg-muted"
        >
          {theme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
          <span className="hidden sm:inline">{theme === 'dark' ? 'Claro' : 'Escuro'}</span>
        </button>
        <div className="flex items-center gap-2 text-sm text-foreground">
          <User className="h-5 w-5 text-primary" />
          <span className="hidden sm:inline">{user?.name || 'Usuário'}</span>
        </div>
        <Button variant="ghost" size="sm" onClick={handleLogout}>
          <LogOut className="h-4 w-4" />
          <span className="hidden sm:inline">Sair</span>
        </Button>
      </div>
    </header>
  )
}
