import { Button } from '@/components/ui/button'
import { Coffee, Menu, X } from 'lucide-react'
import { useState } from 'react'

const navLinks = [
  { label: 'Funcionalidades', href: '#features' },
  { label: 'Planos', href: '#plans' },
  { label: 'Tecnologia', href: '#tech' },
  { label: 'Roadmap', href: '#roadmap' },
]

export function Header() {
  const [open, setOpen] = useState(false)

  return (
    <header className="fixed top-0 left-0 right-0 z-50 bg-white/90 backdrop-blur-md border-b border-gray-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <a href="#" className="flex items-center gap-2">
            <Coffee className="h-8 w-8 text-coffee-green" />
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
            <Button variant="ghost" size="sm">Entrar</Button>
            <Button variant="primary" size="sm">Começar Grátis</Button>
          </div>

          <button
            className="md:hidden p-2"
            onClick={() => setOpen(!open)}
            aria-label="Menu"
          >
            {open ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
          </button>
        </div>

        {open && (
          <div className="md:hidden pb-4 space-y-3">
            {navLinks.map((link) => (
              <a
                key={link.href}
                href={link.href}
                className="block py-2 text-coffee-text hover:text-coffee-green"
                onClick={() => setOpen(false)}
              >
                {link.label}
              </a>
            ))}
            <div className="flex gap-3 pt-2">
              <Button variant="ghost" size="sm">Entrar</Button>
              <Button variant="primary" size="sm">Começar Grátis</Button>
            </div>
          </div>
        )}
      </div>
    </header>
  )
}
