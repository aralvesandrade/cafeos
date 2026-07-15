import { Header } from '@/components/layout/Header'
import { Footer } from '@/components/layout/Footer'
import { Hero } from '@/components/sections/Hero'
import { About } from '@/components/sections/About'
import { Features } from '@/components/sections/Features'
import { CoffeeCycle } from '@/components/sections/CoffeeCycle'
import { Indicators } from '@/components/sections/Indicators'
import { Plans } from '@/components/sections/Plans'
import { CtaSection } from '@/components/sections/CtaSection'

function App() {
  return (
    <div className="min-h-screen bg-ink text-parchment">
      <Header />

      <main>
        <Hero />
        <About />
        <Features />
        <CoffeeCycle />
        <Indicators />
        <Plans />
        <CtaSection />
      </main>

      <Footer />
    </div>
  )
}

export default App
