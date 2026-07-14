import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Sprout, LogIn, UserCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useAuth } from '@/lib/auth'
import { apiRequest } from '@/lib/api'

const profiles = [
  { label: 'Admin Plataforma', email: 'admin@cafeos.com.br', password: 'admin123', role: 'platform_owner' },
  { label: 'Admin Organização', email: 'fernanda@cafeos.com.br', password: '123456', role: 'organization_admin' },
  { label: 'Consultor', email: 'rodrigo@cafeos.com.br', password: '123456', role: 'consultor_externo' },
  { label: 'Proprietário (2 fazendas)', email: 'joao@cafeos.com.br', password: '123456', role: 'proprietario' },
  { label: 'Proprietário (1 fazenda)', email: 'maria@cafeos.com.br', password: '123456', role: 'proprietario' },
]

export function Login() {
  const [email, setEmail] = useState('admin@cafeos.com.br')
  const [password, setPassword] = useState('admin123')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const data = await apiRequest<{
        token: string
        organization_id: string
        user: { id: string; email: string; name: string; role: string }
      }>('/auth/login', {
        method: 'POST',
        body: { email, password },
      })
      login(data.token, data.organization_id, data.user)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao fazer login')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="w-full max-w-sm">
      <div className="text-center mb-8">
        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary/10 mb-4">
          <Sprout className="h-8 w-8 text-primary" />
        </div>
        <h1 className="text-2xl font-bold text-primary">CafeOS</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Plataforma especialista em cafeicultura
        </p>
      </div>

      <form onSubmit={handleSubmit} className="bg-card rounded-xl shadow-sm border border-border p-6 space-y-4">
        <h2 className="text-lg font-semibold text-primary">Entrar</h2>

        {error && (
          <div className="bg-danger-bg text-destructive text-sm rounded-lg px-4 py-2">
            {error}
          </div>
        )}

        <div>
          <label className="block text-sm font-medium text-foreground mb-1">Email</label>
          <Input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="seu@email.com"
            required
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-foreground mb-1">Senha</label>
          <Input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
            required
          />
        </div>

        <Button type="submit" variant="primary" className="w-full" disabled={loading}>
          {loading ? 'Entrando...' : 'Entrar'}
          <LogIn className="h-4 w-4" />
        </Button>
      </form>

      <div className="bg-card rounded-xl shadow-sm border border-border p-4 mt-4">
        <p className="text-xs font-medium text-muted-foreground mb-3 flex items-center gap-1">
          <UserCircle className="h-3 w-3" />
          ACESSO RÁPIDO — SELECIONE UM PERFIL
        </p>
        <div className="grid grid-cols-2 gap-2">
          {profiles.map((p) => (
            <button
              key={p.email}
              type="button"
              onClick={() => {
                setEmail(p.email)
                setPassword(p.password)
                setError('')
              }}
              className={`text-left text-xs px-3 py-2 rounded-lg border transition-colors ${
                email === p.email
                  ? 'border-primary bg-primary/10 text-primary font-medium'
                  : 'border-border text-foreground hover:border-primary/50 hover:bg-muted'
              }`}
            >
              <div className="truncate font-medium">{p.label}</div>
              <div className="truncate text-[10px] opacity-60">{p.role}</div>
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}
