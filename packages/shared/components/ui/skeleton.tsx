import type { HTMLAttributes } from 'react'
import { cn } from '../../lib/utils.ts'

function Skeleton({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn('animate-ds-shimmer rounded-md bg-muted', className)}
      {...props}
    />
  )
}

export { Skeleton }
