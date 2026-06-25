import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Coffee, LogIn } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useAuth } from '@/lib/auth'
import { apiRequest } from '@/lib/api'

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
        tenant_id: string
        user: { id: string; email: string; name: string; role: string }
      }>('/auth/login', {
        method: 'POST',
        body: { email, password },
      })
      login(data.token, data.tenant_id, data.user)
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
        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-coffee-green/10 mb-4">
          <Coffee className="h-8 w-8 text-coffee-green" />
        </div>
        <h1 className="text-2xl font-bold text-coffee-green-dark">CafeOS</h1>
        <p className="text-sm text-coffee-text-light mt-1">
          Plataforma especialista em cafeicultura
        </p>
      </div>

      <form onSubmit={handleSubmit} className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 space-y-4">
        <h2 className="text-lg font-semibold text-coffee-green-dark">Entrar</h2>

        {error && (
          <div className="bg-red-50 text-red-700 text-sm rounded-lg px-4 py-2">
            {error}
          </div>
        )}

        <div>
          <label className="block text-sm font-medium text-coffee-text mb-1">Email</label>
          <Input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="seu@email.com"
            required
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-coffee-text mb-1">Senha</label>
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
    </div>
  )
}
