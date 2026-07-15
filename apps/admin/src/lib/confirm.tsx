import { createContext, useCallback, useContext, useRef, useState, type ReactNode } from 'react'
import { Dialog } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'

interface ConfirmOptions {
  title: string
  description?: string
  confirmLabel?: string
  variant?: 'default' | 'danger'
}

type ConfirmFn = (opts: ConfirmOptions) => Promise<boolean>

const ConfirmContext = createContext<ConfirmFn | null>(null)

export function ConfirmProvider({ children }: { children: ReactNode }) {
  const [opts, setOpts] = useState<ConfirmOptions | null>(null)
  const resolveRef = useRef<(value: boolean) => void>(null)

  const confirm = useCallback<ConfirmFn>((options) => {
    setOpts(options)
    return new Promise((resolve) => {
      resolveRef.current = resolve
    })
  }, [])

  const handle = (result: boolean) => {
    setOpts(null)
    resolveRef.current?.(result)
  }

  return (
    <ConfirmContext.Provider value={confirm}>
      {children}
      <Dialog open={!!opts} onClose={() => handle(false)} title={opts?.title || ''}>
        {opts?.description && <p className="text-sm text-muted-foreground mb-4">{opts.description}</p>}
        <div className="flex justify-end gap-3 pt-2">
          <Button variant="outline" onClick={() => handle(false)}>Cancelar</Button>
          <Button variant={opts?.variant === 'danger' ? 'danger' : 'primary'} onClick={() => handle(true)}>
            {opts?.confirmLabel || 'Confirmar'}
          </Button>
        </div>
      </Dialog>
    </ConfirmContext.Provider>
  )
}

export function useConfirm() {
  const ctx = useContext(ConfirmContext)
  if (!ctx) throw new Error('useConfirm must be used within ConfirmProvider')
  return ctx
}
