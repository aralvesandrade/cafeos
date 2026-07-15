import { cn } from '@/lib/utils'

interface BadgeProps {
  children: React.ReactNode
  variant?: 'default' | 'success' | 'warning'
  className?: string
}

export function Badge({ children, variant = 'default', className }: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center px-3 py-1 rounded-full text-xs font-medium font-mono uppercase tracking-wide',
        {
          'bg-terreiro/15 text-terreiro-light border border-terreiro/30': variant === 'default',
          'bg-gold text-ink': variant === 'success',
          'bg-leaf/20 text-leaf border border-leaf/30': variant === 'warning',
        },
        className,
      )}
    >
      {children}
    </span>
  )
}
