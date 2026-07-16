import type { HTMLAttributes, ReactNode } from 'react'
import { CircleAlert, CircleCheck, Info, TriangleAlert } from 'lucide-react'
import { cn } from '../../lib/utils.ts'

type AlertVariant = 'info' | 'success' | 'warning' | 'error'

const variantStyles: Record<AlertVariant, string> = {
  info: 'border-border bg-card text-foreground',
  success: 'border-success-border bg-success-bg text-success-foreground',
  warning: 'border-warning-border bg-warning-bg text-warning-foreground',
  error: 'border-danger-border bg-danger-bg text-danger-foreground',
}

const variantIcons: Record<AlertVariant, typeof Info> = {
  info: Info,
  success: CircleCheck,
  warning: TriangleAlert,
  error: CircleAlert,
}

interface AlertProps extends Omit<HTMLAttributes<HTMLDivElement>, 'title'> {
  variant?: AlertVariant
  title: ReactNode
  description?: ReactNode
}

function Alert({ className, variant = 'info', title, description, ...props }: AlertProps) {
  const Icon = variantIcons[variant]
  return (
    <div
      role="alert"
      className={cn(
        'flex gap-3 rounded-lg border px-4 py-3.5',
        variantStyles[variant],
        className
      )}
      {...props}
    >
      <Icon className={cn('h-[18px] w-[18px] shrink-0', variant === 'info' && 'text-muted-foreground')} />
      <div>
        <div className="text-sm font-bold">{title}</div>
        {description && <div className="text-[13px] opacity-85">{description}</div>}
      </div>
    </div>
  )
}

export { Alert }
