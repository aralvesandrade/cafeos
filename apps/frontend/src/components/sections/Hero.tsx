import { Button } from '@/components/ui/button'
import { ArrowRight, Sprout, Leaf, BarChart3, Wallet, Bell, Tractor, SplitSquareHorizontal } from 'lucide-react'

const badges = [
  { icon: Leaf, label: 'Gestão completa' },
  { icon: BarChart3, label: 'Indicadores em tempo real' },
  { icon: Wallet, label: 'Controle financeiro' },
  { icon: Bell, label: 'Alertas inteligentes' },
  { icon: Tractor, label: 'Operações agrícolas' },
  { icon: SplitSquareHorizontal, label: 'Rateio de custo' },
]

const manifest = [
  { label: 'TALHÕES', value: '128' },
  { label: 'SACAS/HA', value: '32.4' },
  { label: 'CUSTO/SACA', value: 'R$ 312' },
  { label: 'SAFRAS ATIVAS', value: '6' },
]

export function Hero() {
  return (
    <section className="relative min-h-screen flex items-center bg-background pt-16 overflow-hidden">
      <div className="absolute inset-0 pointer-events-none opacity-[0.07]" style={{ backgroundImage: 'radial-gradient(circle, var(--color-foreground) 1px, transparent 1px)', backgroundSize: '28px 28px' }} />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 relative w-full">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div>
            <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full border border-primary/40 bg-primary/10 text-primary text-xs font-mono uppercase tracking-wider mb-6">
              <Sprout className="h-3.5 w-3.5" />
              Plataforma especialista em cafeicultura
            </div>

            <h1 className="font-display text-4xl sm:text-5xl lg:text-6xl font-semibold text-foreground leading-[1.05] mb-6">
              Café não se gerencia{' '}
              <span className="italic text-primary">de cabeça</span>.
            </h1>

            <p className="text-lg text-muted-foreground mb-8 max-w-lg">
              Gerencie fazendas, talhões, operações, safras e custos em um só
              lugar. Com indicadores precisos e alertas inteligentes para
              máxima produtividade.
            </p>

            <div className="flex flex-col sm:flex-row gap-4 mb-10">
              <a href="#plans">
                <Button variant="primary" size="lg">
                  Assinar agora
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
              </a>
              <a href="#features">
                <Button variant="outline" size="lg">
                  Ver Funcionalidades
                </Button>
              </a>
            </div>

            <div className="flex flex-wrap items-center gap-x-6 gap-y-3 mb-10 text-sm text-muted-foreground">
              {badges.map((b) => (
                <div key={b.label} className="flex items-center gap-2">
                  <b.icon className="h-4 w-4 text-primary" />
                  {b.label}
                </div>
              ))}
            </div>

            {/* Manifesto — assinatura visual, tira de dados estilo ficha de pesagem */}
            <div className="border-t border-border pt-4">
              <div className="flex flex-wrap gap-x-8 gap-y-3 font-mono">
                {manifest.map((m, i) => (
                  <div key={m.label} className={i > 0 ? 'pl-8 border-l border-border' : ''}>
                    <div className="text-[10px] tracking-widest text-muted-foreground mb-1">{m.label}</div>
                    <div className="text-xl text-gold font-medium">{m.value}</div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className="hidden lg:flex items-center justify-center">
            <div className="relative">
              <div className="w-96 h-96 rounded-full border border-border flex items-center justify-center">
                <Sprout className="h-40 w-40 text-primary/20" />
              </div>
              <div className="absolute -top-3 -right-3 bg-card border border-border rounded-sm shadow-lg p-4 rotate-3">
                <BarChart3 className="h-7 w-7 text-primary" />
              </div>
              <div className="absolute -bottom-3 -left-3 bg-card border border-border rounded-sm shadow-lg p-4 -rotate-3">
                <Leaf className="h-7 w-7 text-leaf" />
              </div>
              <div className="absolute top-1/2 -translate-y-1/2 -right-6 bg-card border border-border rounded-sm shadow-lg p-4 rotate-2">
                <Wallet className="h-7 w-7 text-gold" />
              </div>
              <div className="absolute top-1/2 -translate-y-1/2 -left-6 bg-card border border-border rounded-sm shadow-lg p-4 -rotate-2">
                <Bell className="h-7 w-7 text-primary" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
