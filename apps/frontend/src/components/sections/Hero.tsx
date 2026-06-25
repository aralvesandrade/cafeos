import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { ArrowRight, Sprout, Leaf, BarChart3 } from 'lucide-react'
import { LeadModal } from '@/components/ui/LeadModal'

export function Hero() {
  const [showSignup, setShowSignup] = useState(false)

  return (
    <section className="relative min-h-screen flex items-center bg-gradient-to-b from-coffee-beige to-white pt-16">
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 -right-40 w-96 h-96 rounded-full bg-coffee-green/5 blur-3xl" />
        <div className="absolute -bottom-40 -left-40 w-96 h-96 rounded-full bg-coffee-brown/5 blur-3xl" />
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 relative">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div>
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-coffee-green/10 text-coffee-green text-sm font-medium mb-6">
              <Sprout className="h-4 w-4" />
              Plataforma especialista em cafeicultura
            </div>

            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-coffee-green-dark leading-tight mb-6">
              Gestão inteligente para{' '}
              <span className="text-coffee-green">cafeicultura</span>
            </h1>

            <p className="text-lg text-coffee-text-light mb-8 max-w-lg">
              Gerencie fazendas, talhões, operações, safras e custos em um só
              lugar. Com indicadores precisos e alertas inteligentes para
              máxima produtividade.
            </p>

            <div className="flex flex-col sm:flex-row gap-4">
              <Button variant="primary" size="lg" onClick={() => setShowSignup(true)}>
                Começar Grátis
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
              <a href="#features">
                <Button variant="outline" size="lg">
                  Ver Funcionalidades
                </Button>
              </a>
            </div>

            <div className="flex items-center gap-6 mt-10 text-sm text-coffee-text-light">
              <div className="flex items-center gap-2">
                <Leaf className="h-5 w-5 text-coffee-green" />
                Gestão completa
              </div>
              <div className="flex items-center gap-2">
                <BarChart3 className="h-5 w-5 text-coffee-green" />
                Indicadores em tempo real
              </div>
            </div>
          </div>

          <div className="hidden lg:flex items-center justify-center">
            <div className="relative">
              <div className="w-96 h-96 rounded-full bg-gradient-to-br from-coffee-green/20 to-coffee-brown/20 flex items-center justify-center">
                <Sprout className="h-48 w-48 text-coffee-green/30" />
              </div>
              <div className="absolute -top-4 -right-4 bg-white rounded-xl shadow-lg p-4">
                <BarChart3 className="h-8 w-8 text-coffee-green" />
              </div>
              <div className="absolute -bottom-4 -left-4 bg-white rounded-xl shadow-lg p-4">
                <Leaf className="h-8 w-8 text-coffee-green" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <LeadModal open={showSignup} onClose={() => setShowSignup(false)} mode="signup" />
    </section>
  )
}
