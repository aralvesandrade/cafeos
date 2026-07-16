import { useEffect, useRef, useState, useCallback } from 'react'
import { Bell, Check, X } from 'lucide-react'
import { apiRequest } from '@/lib/api'
import { Badge } from '@/components/ui/badge'

interface Alert {
  id: string
  harvest_id: string
  rule_id: string
  message: string
  severity: string
  status: string
  created_at: string
}

const POLL_INTERVAL_MS = 60000

export function NotificationBell() {
  const [alerts, setAlerts] = useState<Alert[]>([])
  const [open, setOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<Alert[]>('/alerts')
      setAlerts(data)
    } catch (err) { console.error(err) }
  }, [])

  useEffect(() => {
    load()
    const interval = setInterval(load, POLL_INTERVAL_MS)
    return () => clearInterval(interval)
  }, [load])

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const updateStatus = async (id: string, status: string) => {
    try {
      await apiRequest(`/alerts/${id}`, { method: 'PUT', body: { status } })
      await load()
    } catch (err) { console.error(err) }
  }

  const openAlerts = alerts.filter((a) => a.status === 'aberto')

  return (
    <div className="relative" ref={containerRef}>
      <button
        onClick={() => setOpen((o) => !o)}
        className="relative inline-flex items-center justify-center h-[34px] w-[34px] rounded-lg border border-border text-foreground hover:bg-muted"
      >
        <Bell className="h-4 w-4" />
        {openAlerts.length > 0 && (
          <span className="absolute -top-1 -right-1 h-4 w-4 rounded-full bg-danger-bg text-danger-foreground border border-danger-border text-[10px] flex items-center justify-center font-medium">
            {openAlerts.length > 9 ? '9+' : openAlerts.length}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 mt-2 w-80 max-h-96 overflow-y-auto rounded-lg border border-border bg-card shadow-lg z-50">
          <div className="p-3 border-b border-border">
            <h3 className="text-sm font-semibold text-foreground">Alertas</h3>
          </div>
          {openAlerts.length === 0 ? (
            <p className="p-4 text-sm text-muted-foreground text-center">Nenhum alerta em aberto.</p>
          ) : (
            <div className="divide-y divide-border">
              {openAlerts.map((a) => (
                <div key={a.id} className="p-3 space-y-2">
                  <div className="flex items-start justify-between gap-2">
                    <Badge variant={a.severity === 'warning' ? 'warning' : 'danger'}>{a.severity}</Badge>
                    <span className="text-xs text-muted-foreground whitespace-nowrap">{new Date(a.created_at).toLocaleDateString('pt-BR')}</span>
                  </div>
                  <p className="text-sm text-foreground">{a.message}</p>
                  <div className="flex justify-end gap-2">
                    <button
                      onClick={() => updateStatus(a.id, 'resolvido')}
                      className="inline-flex items-center gap-1 text-xs text-success-foreground hover:underline"
                    >
                      <Check className="h-3 w-3" /> Resolver
                    </button>
                    <button
                      onClick={() => updateStatus(a.id, 'descartado')}
                      className="inline-flex items-center gap-1 text-xs text-muted-foreground hover:underline"
                    >
                      <X className="h-3 w-3" /> Descartar
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  )
}
