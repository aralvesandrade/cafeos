import { Button } from '@/components/ui/button'
import { useNavigate } from 'react-router-dom'
import { FileQuestion } from 'lucide-react'

export function NotFound() {
  const navigate = useNavigate()

  return (
    <div className="flex flex-col items-center justify-center h-full text-center">
      <FileQuestion className="h-24 w-24 text-primary/30 mb-4" />
      <h1 className="text-3xl font-bold text-primary mb-2">Página não encontrada</h1>
      <p className="text-muted-foreground mb-6">A página que você procura não existe.</p>
      <Button onClick={() => navigate('/')}>Voltar ao Dashboard</Button>
    </div>
  )
}
