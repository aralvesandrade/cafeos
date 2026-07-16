import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { ArrowRight, Mail } from 'lucide-react'
import { LeadModal } from '@/components/ui/LeadModal'

export function CtaSection() {
  const [showContact, setShowContact] = useState(false)

  return (
    <section className="py-20 bg-gradient-to-br from-primary to-primary/60 text-foreground">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <h2 className="font-display text-3xl sm:text-4xl font-semibold mb-4">
          Pronto para transformar sua cafeicultura?
        </h2>
        <p className="text-lg text-foreground/85 max-w-2xl mx-auto mb-8">
          Comece hoje mesmo e descubra como o CafeOS pode ajudar você a
          aumentar a produtividade, reduzir custos e tomar decisões melhores.
        </p>

        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <a href="#plans">
            <Button variant="secondary" size="lg" className="gap-2">
              Assinar agora
              <ArrowRight className="h-5 w-5" />
            </Button>
          </a>
          <Button
            variant="outline"
            size="lg"
            className="border-foreground/40 text-foreground hover:bg-foreground/10 gap-2"
            onClick={() => setShowContact(true)}
          >
            <Mail className="h-5 w-5" />
            Fale Conosco
          </Button>
        </div>
      </div>

      <LeadModal open={showContact} onClose={() => setShowContact(false)} mode="contact" />
    </section>
  )
}
