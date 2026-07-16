import { createContext, useCallback, useContext, useState, type ReactNode } from 'react'
import { CircleCheck, CircleAlert, X } from 'lucide-react'

type ToastKind = 'success' | 'error'

interface Toast {
  id: number
  kind: ToastKind
  title: string
  desc?: string
}

interface ToastContextType {
  success: (title: string, desc?: string) => void
  error: (title: string, desc?: string) => void
}

const ToastContext = createContext<ToastContextType | null>(null)

const AUTO_DISMISS_MS = 4000

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([])

  const dismiss = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id))
  }, [])

  const addToast = useCallback((kind: ToastKind, title: string, desc?: string) => {
    const id = Date.now() + Math.random()
    setToasts((prev) => [...prev, { id, kind, title, desc }])
    setTimeout(() => dismiss(id), AUTO_DISMISS_MS)
  }, [dismiss])

  const success = useCallback((title: string, desc?: string) => addToast('success', title, desc), [addToast])
  const error = useCallback((title: string, desc?: string) => addToast('error', title, desc), [addToast])

  return (
    <ToastContext.Provider value={{ success, error }}>
      {children}
      <div className="fixed bottom-5 right-5 z-[90] flex flex-col gap-2.5 pointer-events-none">
        {toasts.map((t) => (
          <div
            key={t.id}
            className="w-80 bg-card text-card-foreground border border-border rounded-xl shadow-lg p-4 flex gap-3 items-start pointer-events-auto animate-toast-in"
          >
            {t.kind === 'success' ? (
              <CircleCheck className="h-[18px] w-[18px] mt-0.5 shrink-0 text-success-foreground" />
            ) : (
              <CircleAlert className="h-[18px] w-[18px] mt-0.5 shrink-0 text-danger-foreground" />
            )}
            <div className="flex-1 min-w-0">
              <div className="text-sm font-semibold">{t.title}</div>
              {t.desc && <div className="text-xs text-muted-foreground mt-0.5">{t.desc}</div>}
            </div>
            <button
              onClick={() => dismiss(t.id)}
              className="inline-flex items-center justify-center w-6 h-6 rounded-md text-muted-foreground hover:bg-muted hover:text-foreground shrink-0"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  )
}

export function useToast() {
  const ctx = useContext(ToastContext)
  if (!ctx) throw new Error('useToast must be used within ToastProvider')
  return ctx
}
