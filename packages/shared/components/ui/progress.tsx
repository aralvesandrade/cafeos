import type { HTMLAttributes } from 'react'
import { cn } from '../../lib/utils.ts'

interface ProgressProps extends HTMLAttributes<HTMLDivElement> {
  value: number
  max?: number
  label?: string
  tone?: 'primary' | 'gold'
}

function Progress({ className, value, max = 100, label, tone = 'primary', ...props }: ProgressProps) {
  const pct = Math.min(100, Math.max(0, (value / max) * 100))
  return (
    <div className={className} {...props}>
      {label && (
        <div className="mb-1 flex justify-between text-[13px]">
          <span>{label}</span>
          <span className="text-muted-foreground">{Math.round(pct)}%</span>
        </div>
      )}
      <div
        role="progressbar"
        aria-valuenow={value}
        aria-valuemin={0}
        aria-valuemax={max}
        className="h-2 overflow-hidden rounded-full bg-muted"
      >
        <div
          className={cn('h-full rounded-full', tone === 'primary' ? 'bg-primary' : 'bg-gold')}
          style={{ width: `${pct}%` }}
        />
      </div>
    </div>
  )
}

export { Progress }
