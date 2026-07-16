import { createContext, useCallback, useContext, useState, type ReactNode } from 'react'
import { CircleAlert, CircleCheck, X } from 'lucide-react'
import { cn } from '../../lib/utils.ts'

type ToastVariant = 'success' | 'error'

interface ToastItem {
  id: string
  title: string
  description?: string
  variant: ToastVariant
}

interface ToastContextValue {
  toasts: ToastItem[]
  addToast: (toast: Omit<ToastItem, 'id'>) => void
  dismissToast: (id: string) => void
}

const ToastContext = createContext<ToastContextValue | null>(null)

function useToast() {
  const ctx = useContext(ToastContext)
  if (!ctx) throw new Error('useToast must be used within a ToastProvider')
  return ctx
}

const TOAST_DURATION_MS = 4000

function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([])

  const dismissToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id))
  }, [])

  const addToast = useCallback(
    (toast: Omit<ToastItem, 'id'>) => {
      const id = crypto.randomUUID()
      setToasts((prev) => [...prev, { ...toast, id }])
      setTimeout(() => dismissToast(id), TOAST_DURATION_MS)
    },
    [dismissToast]
  )

  return (
    <ToastContext.Provider value={{ toasts, addToast, dismissToast }}>
      {children}
      <Toaster />
    </ToastContext.Provider>
  )
}

function Toaster() {
  const { toasts, dismissToast } = useToast()

  return (
    <div className="pointer-events-none fixed bottom-5 right-5 z-[90] flex flex-col gap-2.5">
      {toasts.map((toast) => {
        const Icon = toast.variant === 'success' ? CircleCheck : CircleAlert
        return (
          <div
            key={toast.id}
            className="animate-toast-in pointer-events-auto w-80 rounded-xl border border-border bg-card p-[15px] shadow-ds-lg"
          >
            <div className="flex items-start gap-2.5">
              <Icon
                className={cn(
                  'h-[18px] w-[18px] shrink-0',
                  toast.variant === 'success' ? 'text-success-foreground' : 'text-danger-foreground'
                )}
              />
              <div className="flex-1">
                <div className="text-sm font-bold">{toast.title}</div>
                {toast.description && (
                  <div className="text-[13px] text-muted-foreground">{toast.description}</div>
                )}
              </div>
              <button
                type="button"
                onClick={() => dismissToast(toast.id)}
                className="flex h-6 w-6 shrink-0 items-center justify-center rounded-md hover:bg-muted"
              >
                <X className="h-3.5 w-3.5" />
              </button>
            </div>
          </div>
        )
      })}
    </div>
  )
}

export { ToastProvider, useToast }
