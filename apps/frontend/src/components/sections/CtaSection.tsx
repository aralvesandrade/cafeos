import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { ArrowRight, Mail } from 'lucide-react'
import { LeadModal } from '@/components/ui/LeadModal'

export function CtaSection() {
  const [showSignup, setShowSignup] = useState(false)
  const [showContact, setShowContact] = useState(false)

  return (
    <section className="py-20 bg-gradient-to-r from-coffee-green-dark to-coffee-green text-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <h2 className="text-3xl sm:text-4xl font-bold mb-4">
          Pronto para transformar sua cafeicultura?
        </h2>
        <p className="text-lg text-coffee-beige/80 max-w-2xl mx-auto mb-8">
          Comece grátis hoje e descubra como o CafeOS pode ajudar você a
          aumentar a produtividade, reduzir custos e tomar decisões melhores.
        </p>

        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button
            variant="secondary"
            size="lg"
            className="gap-2"
            onClick={() => setShowSignup(true)}
          >
            Assinar Grátis
            <ArrowRight className="h-5 w-5" />
          </Button>
          <Button
            variant="outline"
            size="lg"
            className="border-white/30 text-white hover:bg-white/10 gap-2"
            onClick={() => setShowContact(true)}
          >
            <Mail className="h-5 w-5" />
            Falar com Vendas
          </Button>
        </div>
      </div>

      <LeadModal open={showSignup} onClose={() => setShowSignup(false)} mode="signup" />
      <LeadModal open={showContact} onClose={() => setShowContact(false)} mode="contact" />
    </section>
  )
}
