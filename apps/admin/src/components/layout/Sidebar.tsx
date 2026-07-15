import { useState } from 'react'
import { NavLink } from 'react-router-dom'
import { cn } from '@/lib/utils'
import {
  LayoutDashboard,
  Tractor,
  Grid3X3,
  MapPin,
  Calendar,
  Building2,
  Users,
  Sprout,
  X,
  DollarSign,
  Package,
  Truck,
  UserCog,
  CircleDollarSign,
  Tags,
  Wallet,
  SplitSquareHorizontal,
  Sun,
  Moon,
  PanelLeftClose,
  PanelLeftOpen,
  ShieldCheck,
} from 'lucide-react'
import { useAuth } from '@/lib/auth'
import { useTheme } from '@/lib/theme'
import { usePermissions, type Module } from '@/lib/permissions'

const navItems: { to: string; icon: typeof LayoutDashboard; label: string; module: Module }[] = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard', module: 'dashboard' },
  { to: '/farms', icon: MapPin, label: 'Fazendas', module: 'farms' },
  { to: '/plots', icon: Grid3X3, label: 'Talhões', module: 'farms' },
  { to: '/operations', icon: Tractor, label: 'Operações', module: 'operations' },
  { to: '/operation-types', icon: Tags, label: 'Tipos de Operação', module: 'farms' },
  { to: '/harvests', icon: Calendar, label: 'Safras', module: 'harvests' },
]

const phase2Items: { to: string; icon: typeof LayoutDashboard; label: string; module: Module }[] = [
  { to: '/financial', icon: DollarSign, label: 'Financeiro', module: 'financial' },
  { to: '/cost-centers', icon: CircleDollarSign, label: 'Centros de Custo', module: 'financial' },
  { to: '/budgets', icon: Wallet, label: 'Orçamento', module: 'financial' },
  { to: '/cost-allocations', icon: SplitSquareHorizontal, label: 'Rateio de Custo', module: 'financial' },
  { to: '/stock', icon: Package, label: 'Estoque', module: 'resources' },
  { to: '/fleet', icon: Truck, label: 'Frota', module: 'resources' },
  { to: '/labor', icon: UserCog, label: 'Equipes', module: 'resources' },
  { to: '/team', icon: Users, label: 'Minha Equipe', module: 'users' },
]

const adminItems = [
  { to: '/organizations', icon: Building2, label: 'Organizações' },
  { to: '/users', icon: Users, label: 'Usuários' },
]

interface SidebarProps {
  open: boolean
  onClose: () => void
}

const COLLAPSED_KEY = 'cafeos-sidebar-collapsed'

export function Sidebar({ open, onClose }: SidebarProps) {
  const { user } = useAuth()
  const { theme, toggleTheme } = useTheme()
  const isAdmin = user?.role === 'platform_owner'
  const { access } = usePermissions()
  const [collapsed, setCollapsed] = useState(() => localStorage.getItem(COLLAPSED_KEY) === 'true')
  const visibleNavItems = navItems.filter((item) => (access?.[item.module] ?? 'none') !== 'none')
  const visiblePhase2Items = phase2Items.filter((item) => (access?.[item.module] ?? 'none') !== 'none')
  const canConfigurePermissions = (access?.permissions ?? 'none') !== 'none'

  const toggleCollapsed = () => {
    setCollapsed((prev) => {
      const next = !prev
      localStorage.setItem(COLLAPSED_KEY, String(next))
      return next
    })
  }

  const linkClass = ({ isActive }: { isActive: boolean }) =>
    cn(
      'flex items-center gap-3 rounded-md text-sm transition-colors border-l-2',
      collapsed ? 'justify-center px-2 py-2' : 'px-3 py-2',
      isActive
        ? 'bg-sidebar-active-bg text-sidebar-active-foreground font-medium border-sidebar-accent'
        : 'border-transparent text-sidebar-foreground/80 hover:text-sidebar-foreground hover:bg-sidebar-active-bg/50'
    )

  return (
    <>
      {open && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={onClose}
        />
      )}

      <aside
        className={cn(
          'fixed top-0 left-0 z-50 h-full w-64 bg-sidebar text-sidebar-foreground border-r border-sidebar-border flex flex-col transition-all lg:translate-x-0 lg:static lg:z-auto',
          open ? 'translate-x-0' : '-translate-x-full',
          collapsed ? 'lg:w-16' : 'lg:w-64'
        )}
      >
        <div className={cn('flex items-center h-16 border-b border-sidebar-border', collapsed ? 'justify-center px-2' : 'justify-between px-4')}>
          <NavLink to="/" className="flex items-center gap-2 font-display font-semibold text-lg text-sidebar-active-foreground">
            <span className="w-7 h-7 rounded-md bg-sidebar-accent text-sidebar-active-foreground flex items-center justify-center shrink-0">
              <Sprout className="h-4 w-4" />
            </span>
            {!collapsed && 'CafeOS'}
          </NavLink>
          <button onClick={onClose} className="lg:hidden text-sidebar-foreground/60 hover:text-sidebar-foreground">
            <X className="h-5 w-5" />
          </button>
        </div>

        <nav className="flex-1 py-4 space-y-1 px-3 overflow-y-auto overflow-x-hidden">
          {visibleNavItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.to === '/'}
              onClick={onClose}
              title={collapsed ? item.label : undefined}
              className={linkClass}
            >
              <item.icon className="h-5 w-5 shrink-0" />
              {!collapsed && item.label}
            </NavLink>
          ))}

          {visiblePhase2Items.length > 0 && (
            <div className="pt-4 mt-4 border-t border-sidebar-border">
              {!collapsed && (
                <p className="px-3 text-xs font-medium text-sidebar-foreground/50 uppercase tracking-wider mb-2">
                  Gestão
                </p>
              )}
              {visiblePhase2Items.map((item) => (
                <NavLink
                  key={item.to}
                  to={item.to}
                  onClick={onClose}
                  title={collapsed ? item.label : undefined}
                  className={linkClass}
                >
                  <item.icon className="h-5 w-5 shrink-0" />
                  {!collapsed && item.label}
                </NavLink>
              ))}
            </div>
          )}

          {(isAdmin || canConfigurePermissions) && (
            <div className="pt-4 mt-4 border-t border-sidebar-border">
              {!collapsed && (
                <p className="px-3 text-xs font-medium text-sidebar-foreground/50 uppercase tracking-wider mb-2">
                  Administração
                </p>
              )}
              {isAdmin && adminItems.map((item) => (
                <NavLink
                  key={item.to}
                  to={item.to}
                  onClick={onClose}
                  title={collapsed ? item.label : undefined}
                  className={linkClass}
                >
                  <item.icon className="h-5 w-5 shrink-0" />
                  {!collapsed && item.label}
                </NavLink>
              ))}
              {canConfigurePermissions && (
                <NavLink
                  to="/permissions"
                  onClick={onClose}
                  title={collapsed ? 'Permissões' : undefined}
                  className={linkClass}
                >
                  <ShieldCheck className="h-5 w-5 shrink-0" />
                  {!collapsed && 'Permissões'}
                </NavLink>
              )}
            </div>
          )}
        </nav>

        <div className="px-3 py-4 border-t border-sidebar-border space-y-1">
          <button
            onClick={toggleTheme}
            title={collapsed ? (theme === 'dark' ? 'Modo Claro' : 'Modo Escuro') : undefined}
            className={cn(
              'flex items-center gap-3 rounded-lg text-sm w-full text-sidebar-foreground/80 hover:text-sidebar-foreground hover:bg-sidebar-active-bg/50 transition-colors',
              collapsed ? 'justify-center px-2 py-2' : 'px-3 py-2'
            )}
          >
            {theme === 'dark' ? <Sun className="h-5 w-5 shrink-0" /> : <Moon className="h-5 w-5 shrink-0" />}
            {!collapsed && (theme === 'dark' ? 'Modo Claro' : 'Modo Escuro')}
          </button>
          <button
            onClick={toggleCollapsed}
            title={collapsed ? 'Expandir menu' : undefined}
            className={cn(
              'hidden lg:flex items-center gap-3 rounded-lg text-sm w-full text-sidebar-foreground/80 hover:text-sidebar-foreground hover:bg-sidebar-active-bg/50 transition-colors',
              collapsed ? 'justify-center px-2 py-2' : 'px-3 py-2'
            )}
          >
            {collapsed ? <PanelLeftOpen className="h-5 w-5 shrink-0" /> : <PanelLeftClose className="h-5 w-5 shrink-0" />}
            {!collapsed && 'Recolher menu'}
          </button>
        </div>
      </aside>
    </>
  )
}
