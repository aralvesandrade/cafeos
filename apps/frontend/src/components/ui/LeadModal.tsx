import { useState } from 'react'
import { Dialog } from './dialog'
import { Button } from './button'
import { Mail, Send, CheckCircle } from 'lucide-react'

interface LeadModalProps {
  open: boolean
  onClose: () => void
  mode: 'signup' | 'contact'
}

export function LeadModal({ open, onClose, mode }: LeadModalProps) {
  const [form, setForm] = useState({ name: '', email: '', phone: '', message: '' })
  const [sent, setSent] = useState(false)
  const [sending, setSending] = useState(false)

  const isSignup = mode === 'signup'

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.name || !form.email) return

    setSending(true)
    // Simulate sending — no backend yet, just show success
    await new Promise((r) => setTimeout(r, 800))
    setSending(false)
    setSent(true)
  }

  if (sent) {
    return (
      <Dialog open={open} onClose={() => { setSent(false); onClose() }} title="">
        <div className="text-center py-6">
          <CheckCircle className="h-12 w-12 text-coffee-green mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-coffee-green-dark mb-2">
            {isSignup ? 'Cadastro realizado!' : 'Mensagem enviada!'}
          </h3>
          <p className="text-sm text-coffee-text-light mb-6">
            {isSignup
              ? 'Em breve você receberá um email com instruções de acesso.'
              : 'Nosso time comercial entrará em contato em até 24 horas.'}
          </p>
          <Button variant="primary" onClick={() => { setSent(false); onClose() }}>
            Fechar
          </Button>
        </div>
      </Dialog>
    )
  }

  return (
    <Dialog open={open} onClose={onClose} title={isSignup ? 'Criar Conta Grátis' : 'Falar com Vendas'}>
      <form onSubmit={handleSubmit} className="space-y-4">
        {isSignup && (
          <p className="text-sm text-coffee-text-light mb-2">
            Preencha seus dados para começar a usar o CafeOS gratuitamente.
          </p>
        )}

        <div>
          <label className="block text-sm font-medium text-coffee-text mb-1">Nome</label>
          <input
            className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-coffee-green/50"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            required
            placeholder="Seu nome"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-coffee-text mb-1">Email</label>
          <input
            type="email"
            className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-coffee-green/50"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
            required
            placeholder="seu@email.com"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-coffee-text mb-1">Telefone</label>
          <input
            className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-coffee-green/50"
            value={form.phone}
            onChange={(e) => setForm({ ...form, phone: e.target.value })}
            placeholder="(DDD) 99999-9999"
          />
        </div>

        {!isSignup && (
          <div>
            <label className="block text-sm font-medium text-coffee-text mb-1">Mensagem</label>
            <textarea
              className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-coffee-green/50 min-h-[80px]"
              value={form.message}
              onChange={(e) => setForm({ ...form, message: e.target.value })}
              placeholder="Conte-nos sobre sua necessidade..."
              rows={3}
            />
          </div>
        )}

        <Button type="submit" variant="primary" className="w-full gap-2" disabled={sending}>
          {sending ? 'Enviando...' : isSignup ? 'Assinar Grátis' : 'Enviar Mensagem'}
          <Send className="h-4 w-4" />
        </Button>
      </form>
    </Dialog>
  )
}
