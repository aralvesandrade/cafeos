import { Sprout, Mail, ArrowUpRight } from 'lucide-react'
import { ADMIN_URL } from '@/lib/api'

const productLinks = [
  { label: 'Funcionalidades', href: '#features' },
  { label: 'Ciclo do Café', href: '#coffee-cycle' },
  { label: 'Indicadores Estratégicos', href: '#indicators' },
  { label: 'Planos', href: '#plans' },
]

const companyLinks = [
  { label: 'Por que o CafeOS?', href: '#about' },
  { label: 'Acessar plataforma', href: ADMIN_URL, external: true },
]

export function Footer() {
  return (
    <footer className="bg-card text-foreground border-t border-border">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div className="col-span-1 md:col-span-2">
            <a href="#" className="flex items-center gap-2 mb-4">
              <Sprout className="h-6 w-6 text-primary" />
              <span className="font-display text-lg font-semibold">CafeOS</span>
            </a>
            <p className="text-muted-foreground text-sm max-w-md mb-6">
              A plataforma especialista em cafeicultura. Gestão operacional,
              produtiva, financeira e analítica para propriedades cafeeiras —
              do talhão à colheita, com indicadores precisos em cada fase da
              safra.
            </p>
            <a
              href="mailto:contato@cafeos.com.br"
              className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-primary transition-colors"
            >
              <Mail className="h-4 w-4" /> contato@cafeos.com.br
            </a>
          </div>

          <div>
            <h4 className="font-display font-semibold mb-4">Produto</h4>
            <ul className="space-y-2 text-sm text-muted-foreground">
              {productLinks.map((link) => (
                <li key={link.href}>
                  <a href={link.href} className="hover:text-foreground transition-colors">
                    {link.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h4 className="font-display font-semibold mb-4">Empresa</h4>
            <ul className="space-y-2 text-sm text-muted-foreground">
              {companyLinks.map((link) => (
                <li key={link.href}>
                  <a
                    href={link.href}
                    target={link.external ? '_blank' : undefined}
                    rel={link.external ? 'noopener noreferrer' : undefined}
                    className="inline-flex items-center gap-1 hover:text-foreground transition-colors"
                  >
                    {link.label}
                    {link.external && <ArrowUpRight className="h-3.5 w-3.5" />}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <div className="mt-8 pt-8 border-t border-border flex flex-col sm:flex-row items-center justify-between gap-3 text-sm text-muted-foreground font-mono">
          <p>&copy; {new Date().getFullYear()} CafeOS. Todos os direitos reservados.</p>
          <p className="text-xs tracking-wide">Feito para quem cultiva café, sacas depois de sacas.</p>
        </div>
      </div>
    </footer>
  )
}
