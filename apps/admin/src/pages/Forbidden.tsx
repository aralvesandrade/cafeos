import { Button } from '@/components/ui/button'
import { useNavigate } from 'react-router-dom'
import { ShieldOff } from 'lucide-react'

export function Forbidden() {
  const navigate = useNavigate()

  return (
    <div className="flex flex-col items-center justify-center h-full text-center">
      <ShieldOff className="h-24 w-24 text-primary/30 mb-4" />
      <h1 className="font-display text-3xl font-semibold text-foreground mb-2">Sem acesso a esta tela</h1>
      <p className="text-muted-foreground mb-6">Seu perfil não tem permissão para essa área. Fale com quem administra sua organização se precisar de acesso.</p>
      <Button onClick={() => navigate('/')}>Voltar ao Dashboard</Button>
    </div>
  )
}
