import { Sprout, Mail, Globe, ExternalLink } from 'lucide-react'

export function Footer() {
  return (
    <footer className="bg-ink-raised text-parchment border-t border-rule">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div className="col-span-1 md:col-span-2">
            <div className="flex items-center gap-2 mb-4">
              <Sprout className="h-6 w-6 text-terreiro" />
              <span className="font-display text-lg font-semibold">CafeOS</span>
            </div>
            <p className="text-muted text-sm max-w-md">
              A plataforma especialista em cafeicultura. Gestão operacional,
              produtiva, financeira e analítica para propriedades cafeeiras.
            </p>
          </div>

          <div>
            <h4 className="font-display font-semibold mb-4">Produto</h4>
            <ul className="space-y-2 text-sm text-muted">
              <li><a href="#features" className="hover:text-parchment transition-colors">Funcionalidades</a></li>
              <li><a href="#plans" className="hover:text-parchment transition-colors">Planos</a></li>
            </ul>
          </div>

          <div>
            <h4 className="font-display font-semibold mb-4">Contato</h4>
            <ul className="space-y-2 text-sm text-muted">
              <li className="flex items-center gap-2">
                <Mail className="h-4 w-4" /> contato@cafeos.com.br
              </li>
              <li className="flex items-center gap-2">
                <Globe className="h-4 w-4" /> github.com/cafeos
              </li>
              <li className="flex items-center gap-2">
                <ExternalLink className="h-4 w-4" /> linkedin.com/company/cafeos
              </li>
            </ul>
          </div>
        </div>

        <div className="mt-8 pt-8 border-t border-rule text-center text-sm text-muted font-mono">
          <p>&copy; {new Date().getFullYear()} CafeOS. Todos os direitos reservados.</p>
        </div>
      </div>
    </footer>
  )
}
