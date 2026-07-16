import { type ReactNode, useEffect, useRef } from 'react'
import { X } from 'lucide-react'

interface DialogProps {
  open: boolean
  onClose: () => void
  title: string
  children: ReactNode
}

export function Dialog({ open, onClose, title, children }: DialogProps) {
  const overlayRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    if (open) document.addEventListener('keydown', handleEsc)
    return () => document.removeEventListener('keydown', handleEsc)
  }, [open, onClose])

  if (!open) return null

  return (
    <div
      ref={overlayRef}
      className="fixed inset-0 z-50 bg-black/50 p-4 overflow-y-auto"
      onClick={(e) => {
        if (e.target === overlayRef.current) onClose()
      }}
    >
      <div className="min-h-full flex items-center justify-center">
        <div className="relative w-full max-w-md rounded-sm bg-card p-6 shadow-lg my-8">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-card-foreground/40 hover:text-card-foreground"
        >
          <X className="h-5 w-5" />
        </button>
        <h2 className="font-display text-lg font-semibold text-card-foreground mb-4">{title}</h2>
        {children}
      </div>
      </div>
    </div>
  )
}
