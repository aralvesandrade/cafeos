import type { ReactNode } from 'react'

interface FieldProps {
  label: string
  required?: boolean
  className?: string
  children: ReactNode
}

export function Field({ label, required, className, children }: FieldProps) {
  return (
    <div className={className}>
      <label className="block text-sm font-medium text-foreground mb-1">
        {label}
        {required && <span className="text-destructive"> *</span>}
      </label>
      {children}
    </div>
  )
}
