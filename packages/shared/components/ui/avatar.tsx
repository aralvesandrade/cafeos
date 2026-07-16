import type { HTMLAttributes } from 'react'
import { cn } from '../../lib/utils.ts'

interface AvatarProps extends HTMLAttributes<HTMLDivElement> {
  initials: string
  size?: 32 | 40
  tone?: 'primary' | 'gold' | 'secondary' | 'muted'
}

const toneStyles: Record<NonNullable<AvatarProps['tone']>, string> = {
  primary: 'bg-primary text-primary-foreground',
  gold: 'bg-gold text-gold-foreground',
  secondary: 'bg-secondary text-secondary-foreground',
  muted: 'bg-muted text-muted-foreground',
}

function Avatar({ className, initials, size = 40, tone = 'primary', ...props }: AvatarProps) {
  return (
    <div
      className={cn(
        'flex shrink-0 items-center justify-center rounded-full font-bold',
        size === 40 ? 'h-10 w-10 text-sm' : 'h-8 w-8 text-xs',
        toneStyles[tone],
        className
      )}
      {...props}
    >
      {initials}
    </div>
  )
}

interface AvatarGroupProps extends HTMLAttributes<HTMLDivElement> {}

function AvatarGroup({ className, children, ...props }: AvatarGroupProps) {
  return (
    <div className={cn('flex', className)} {...props}>
      {Array.isArray(children)
        ? children.map((child, i) => (
            <div key={i} className={cn(i > 0 && '-ml-2.5', 'rounded-full border-2 border-card')}>
              {child}
            </div>
          ))
        : children}
    </div>
  )
}

export { Avatar, AvatarGroup }
