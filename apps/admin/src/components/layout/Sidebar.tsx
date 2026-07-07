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
} from 'lucide-react'
import { useAuth } from '@/lib/auth'

const navItems = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/farms', icon: MapPin, label: 'Fazendas' },
  { to: '/plots', icon: Grid3X3, label: 'Talhões' },
  { to: '/operations', icon: Tractor, label: 'Operações' },
  { to: '/harvests', icon: Calendar, label: 'Safras' },
]

const phase2Items = [
  { to: '/financial', icon: DollarSign, label: 'Financeiro' },
  { to: '/cost-centers', icon: CircleDollarSign, label: 'Centros de Custo' },
  { to: '/stock', icon: Package, label: 'Estoque' },
  { to: '/fleet', icon: Truck, label: 'Frota' },
  { to: '/labor', icon: UserCog, label: 'Equipes' },
]

const adminItems = [
  { to: '/tenants', icon: Building2, label: 'Tenants' },
  { to: '/users', icon: Users, label: 'Usuários' },
]

interface SidebarProps {
  open: boolean
  onClose: () => void
}

export function Sidebar({ open, onClose }: SidebarProps) {
  const { user } = useAuth()
  const isAdmin = user?.role === 'platform_owner'
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
          'fixed top-0 left-0 z-50 h-full w-64 bg-coffee-green-dark text-white flex flex-col transition-transform lg:translate-x-0 lg:static lg:z-auto',
          open ? 'translate-x-0' : '-translate-x-full'
        )}
      >
        <div className="flex items-center justify-between px-4 h-16 border-b border-white/10">
          <NavLink to="/" className="flex items-center gap-2 font-bold text-lg">
            <Sprout className="h-6 w-6" />
            CafeOS
          </NavLink>
          <button onClick={onClose} className="lg:hidden text-white/60 hover:text-white">
            <X className="h-5 w-5" />
          </button>
        </div>

        <nav className="flex-1 py-4 space-y-1 px-3">
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.to === '/'}
              onClick={onClose}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors',
                  isActive
                    ? 'bg-white/20 text-white font-medium'
                    : 'text-white/70 hover:text-white hover:bg-white/10'
                )
              }
            >
              <item.icon className="h-5 w-5" />
              {item.label}
            </NavLink>
          ))}

          <div className="pt-4 mt-4 border-t border-white/10">
            <p className="px-3 text-xs font-medium text-white/40 uppercase tracking-wider mb-2">
              Gestão
            </p>
            {phase2Items.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                onClick={onClose}
                className={({ isActive }) =>
                  cn(
                    'flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors',
                    isActive
                      ? 'bg-white/20 text-white font-medium'
                      : 'text-white/70 hover:text-white hover:bg-white/10'
                  )
                }
              >
                <item.icon className="h-5 w-5" />
                {item.label}
              </NavLink>
            ))}
          </div>

          {isAdmin && (
            <div className="pt-4 mt-4 border-t border-white/10">
              <p className="px-3 text-xs font-medium text-white/40 uppercase tracking-wider mb-2">
                Administração
              </p>
              {adminItems.map((item) => (
                <NavLink
                  key={item.to}
                  to={item.to}
                  onClick={onClose}
                  className={({ isActive }) =>
                    cn(
                      'flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors',
                      isActive
                        ? 'bg-white/20 text-white font-medium'
                        : 'text-white/70 hover:text-white hover:bg-white/10'
                    )
                  }
                >
                  <item.icon className="h-5 w-5" />
                  {item.label}
                </NavLink>
              ))}
            </div>
          )}
        </nav>
      </aside>
    </>
  )
}
