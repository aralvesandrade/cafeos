import { useState } from 'react'
import { Dialog } from './dialog'
import { Button } from './button'
import { Mail, MessageCircle, Send, CheckCircle, ChevronLeft } from 'lucide-react'

interface LeadModalProps {
  open: boolean
  onClose: () => void
  mode: 'signup' | 'contact'
  plan?: string
}

type ContactStep = 'channel' | 'form'

export function LeadModal({ open, onClose, mode, plan }: LeadModalProps) {
  const [step, setStep] = useState<ContactStep>('channel')
  const [channel, setChannel] = useState<'email' | 'whatsapp'>('email')
  const [form, setForm] = useState({ name: '', email: '', phone: '', message: '' })
  const [sent, setSent] = useState(false)
  const [sending, setSending] = useState(false)

  const isSignup = mode === 'signup'

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.name || !form.email) return

    setSending(true)
    await new Promise((r) => setTimeout(r, 800))
    setSending(false)
    setSent(true)
  }

  const handleClose = () => {
    setSent(false)
    setStep('channel')
    setForm({ name: '', email: '', phone: '', message: '' })
    onClose()
  }

  if (sent) {
    return (
      <Dialog open={open} onClose={handleClose} title="">
        <div className="text-center py-6">
          <CheckCircle className="h-12 w-12 text-primary mx-auto mb-4" />
          <h3 className="font-display text-lg font-semibold text-card-foreground mb-2">
            {isSignup ? 'Cadastro realizado!' : 'Mensagem enviada!'}
          </h3>
          <p className="text-sm text-card-foreground/60 mb-6">
            {isSignup
              ? 'Em breve você receberá um email com instruções de acesso.'
              : 'Nosso time comercial entrará em contato em até 24 horas pelo canal selecionado.'}
          </p>
          <Button variant="primary" onClick={handleClose}>
            Fechar
          </Button>
        </div>
      </Dialog>
    )
  }

  if (isSignup) {
    return (
      <Dialog open={open} onClose={handleClose} title="Criar Conta">
        <form onSubmit={handleSubmit} className="space-y-4">
          {plan && (
            <div className="bg-card-foreground/5 rounded-lg px-4 py-3 text-sm">
              <span className="text-card-foreground/60">Plano selecionado: </span>
              <span className="font-semibold text-card-foreground">{plan}</span>
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-card-foreground mb-1">Nome</label>
            <input
              className="w-full rounded-lg border border-card-foreground/15 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/40"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              required
              placeholder="Seu nome"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-card-foreground mb-1">Email</label>
            <input
              type="email"
              className="w-full rounded-lg border border-card-foreground/15 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/40"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              required
              placeholder="seu@email.com"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-card-foreground mb-1">Telefone</label>
            <input
              className="w-full rounded-lg border border-card-foreground/15 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/40"
              value={form.phone}
              onChange={(e) => setForm({ ...form, phone: e.target.value })}
              placeholder="(DDD) 99999-9999"
            />
          </div>

          <Button type="submit" variant="primary" className="w-full gap-2" disabled={sending}>
            {sending ? 'Enviando...' : 'Assinar'}
            <Send className="h-4 w-4" />
          </Button>
        </form>
      </Dialog>
    )
  }

  // Contact mode
  return (
    <Dialog open={open} onClose={handleClose} title="Fale Conosco">
      {step === 'channel' ? (
        <div className="space-y-4">
          <p className="text-sm text-card-foreground/60 mb-2">
            Escolha o canal de atendimento:
          </p>

          <button
            onClick={() => { setChannel('whatsapp'); setStep('form') }}
            className="w-full flex items-center gap-4 p-4 rounded-lg border border-card-foreground/10 hover:border-primary/50 hover:bg-card-foreground/5 transition-colors text-left"
          >
            <div className="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center">
              <MessageCircle className="h-5 w-5 text-green-600" />
            </div>
            <div>
              <p className="font-medium text-card-foreground">WhatsApp</p>
              <p className="text-xs text-card-foreground/60">Atendimento rápido via mensagem</p>
            </div>
          </button>

          <button
            onClick={() => { setChannel('email'); setStep('form') }}
            className="w-full flex items-center gap-4 p-4 rounded-lg border border-card-foreground/10 hover:border-primary/50 hover:bg-card-foreground/5 transition-colors text-left"
          >
            <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center">
              <Mail className="h-5 w-5 text-blue-600" />
            </div>
            <div>
              <p className="font-medium text-card-foreground">Email</p>
              <p className="text-xs text-card-foreground/60">Respondemos em até 24 horas</p>
            </div>
          </button>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="space-y-4">
          <button
            type="button"
            onClick={() => setStep('channel')}
            className="flex items-center gap-1 text-sm text-card-foreground/60 hover:text-primary transition-colors mb-2"
          >
            <ChevronLeft className="h-4 w-4" />
            Voltar
          </button>

          <div className="flex items-center gap-2 text-sm text-card-foreground mb-2">
            {channel === 'whatsapp' ? (
              <MessageCircle className="h-4 w-4 text-green-600" />
            ) : (
              <Mail className="h-4 w-4 text-blue-600" />
            )}
            <span className="font-medium">
              {channel === 'whatsapp' ? 'WhatsApp' : 'Email'}
            </span>
          </div>

          <div>
            <label className="block text-sm font-medium text-card-foreground mb-1">Nome</label>
            <input
              className="w-full rounded-lg border border-card-foreground/15 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/40"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              required
              placeholder="Seu nome"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-card-foreground mb-1">
              {channel === 'whatsapp' ? 'WhatsApp' : 'Email'}
            </label>
            {channel === 'whatsapp' ? (
              <input
                className="w-full rounded-lg border border-card-foreground/15 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/40"
                value={form.phone}
                onChange={(e) => setForm({ ...form, phone: e.target.value })}
                placeholder="(DDD) 99999-9999"
              />
            ) : (
              <input
                type="email"
                className="w-full rounded-lg border border-card-foreground/15 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/40"
                value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })}
                placeholder="seu@email.com"
              />
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-card-foreground mb-1">Mensagem</label>
            <textarea
              className="w-full rounded-lg border border-card-foreground/15 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/40 min-h-[80px]"
              value={form.message}
              onChange={(e) => setForm({ ...form, message: e.target.value })}
              placeholder="Conte-nos sobre sua necessidade..."
              rows={3}
            />
          </div>

          <Button type="submit" variant="primary" className="w-full gap-2" disabled={sending}>
            {sending ? 'Enviando...' : 'Enviar'}
            <Send className="h-4 w-4" />
          </Button>
        </form>
      )}
    </Dialog>
  )
}
