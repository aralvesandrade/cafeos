import { useState, type ReactNode } from 'react'
import { cn } from '../../lib/utils.ts'

interface TooltipProps {
  content: ReactNode
  children: ReactNode
  className?: string
}

function Tooltip({ content, children, className }: TooltipProps) {
  const [open, setOpen] = useState(false)

  return (
    <span
      className="relative inline-flex"
      onMouseEnter={() => setOpen(true)}
      onMouseLeave={() => setOpen(false)}
      onFocus={() => setOpen(true)}
      onBlur={() => setOpen(false)}
    >
      {children}
      <span
        role="tooltip"
        className={cn(
          'pointer-events-none absolute bottom-[calc(100%+9px)] left-1/2 z-30 -translate-x-1/2 whitespace-nowrap rounded-md bg-foreground px-2.5 py-[7px] text-xs font-semibold text-background shadow-ds-md transition-all duration-150',
          open ? 'opacity-100' : 'opacity-0',
          className
        )}
      >
        {content}
        <span className="absolute left-1/2 top-full h-0 w-0 -translate-x-1/2 border-[5px] border-transparent border-t-foreground" />
      </span>
    </span>
  )
}

export { Tooltip }
