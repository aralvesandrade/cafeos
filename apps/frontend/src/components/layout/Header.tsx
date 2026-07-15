import { Button } from '@/components/ui/button'
import { Sprout, Menu, X } from 'lucide-react'
import { useState } from 'react'

const navLinks = [
  { label: 'Por que o CafeOS?', href: '#about' },
  { label: 'Funcionalidades', href: '#features' },
  { label: 'Ciclo do Café', href: '#coffee-cycle' },
  { label: 'Indicadores Estratégicos', href: '#indicators' },
  { label: 'Planos', href: '#plans' },
]

export function Header() {
  const [menuOpen, setMenuOpen] = useState(false)

  return (
    <header className="fixed top-0 left-0 right-0 z-50 bg-white/90 backdrop-blur-md border-b border-gray-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <a href="#" className="flex items-center gap-2">
            <Sprout className="h-8 w-8 text-coffee-green" />
            <span className="text-xl font-bold text-coffee-green-dark">CafeOS</span>
          </a>

          <nav className="hidden md:flex items-center gap-8">
            {navLinks.map((link) => (
              <a
                key={link.href}
                href={link.href}
                className="text-sm text-coffee-text hover:text-coffee-green transition-colors"
              >
                {link.label}
              </a>
            ))}
          </nav>

          <div className="hidden md:flex items-center gap-4">
            <a href="http://localhost:5174" target="_blank" rel="noopener noreferrer">
              <Button variant="ghost" size="sm">Entrar</Button>
            </a>
          </div>

          <button
            className="md:hidden p-2"
            onClick={() => setMenuOpen(!menuOpen)}
            aria-label="Menu"
          >
            {menuOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
          </button>
        </div>

        {menuOpen && (
          <div className="md:hidden pb-4 space-y-3">
            {navLinks.map((link) => (
              <a
                key={link.href}
                href={link.href}
                className="block py-2 text-coffee-text hover:text-coffee-green"
                onClick={() => setMenuOpen(false)}
              >
                {link.label}
              </a>
            ))}
            <div className="flex gap-3 pt-2">
              <a href="http://localhost:5174" target="_blank" rel="noopener noreferrer">
                <Button variant="ghost" size="sm">Entrar</Button>
              </a>
            </div>
          </div>
        )}
      </div>
    </header>
  )
}
