import {
  createContext,
  useContext,
  useState,
  type ButtonHTMLAttributes,
  type HTMLAttributes,
  type ReactNode,
} from 'react'
import { cn } from '../../lib/utils.ts'

interface DropdownMenuContextValue {
  open: boolean
  setOpen: (open: boolean) => void
}

const DropdownMenuContext = createContext<DropdownMenuContextValue | null>(null)

function useDropdownMenu() {
  const ctx = useContext(DropdownMenuContext)
  if (!ctx) throw new Error('DropdownMenu components must be used within a DropdownMenu')
  return ctx
}

interface DropdownMenuProps {
  children: ReactNode
  className?: string
}

function DropdownMenu({ children, className }: DropdownMenuProps) {
  const [open, setOpen] = useState(false)
  return (
    <DropdownMenuContext.Provider value={{ open, setOpen }}>
      <div className={cn('relative inline-block', className)}>{children}</div>
    </DropdownMenuContext.Provider>
  )
}

function DropdownMenuTrigger({
  className,
  onClick,
  children,
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement>) {
  const { open, setOpen } = useDropdownMenu()
  return (
    <button
      type="button"
      className={className}
      onClick={(e) => {
        onClick?.(e)
        setOpen(!open)
      }}
      {...props}
    >
      {children}
    </button>
  )
}

function DropdownMenuContent({ className, children }: HTMLAttributes<HTMLDivElement>) {
  const { open, setOpen } = useDropdownMenu()
  if (!open) return null

  return (
    <>
      <div className="fixed inset-0 z-20" onClick={() => setOpen(false)} />
      <div
        className={cn(
          'animate-ds-pop absolute left-0 top-[calc(100%+6px)] z-30 min-w-[184px] rounded-lg border border-border bg-card p-[5px] shadow-ds-md',
          className
        )}
      >
        {children}
      </div>
    </>
  )
}

interface DropdownMenuItemProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  destructive?: boolean
}

function DropdownMenuItem({ className, destructive, onClick, children, ...props }: DropdownMenuItemProps) {
  const { setOpen } = useDropdownMenu()
  return (
    <button
      type="button"
      className={cn(
        'flex w-full items-center gap-2.5 rounded-md px-[9px] py-2 text-left text-sm hover:bg-muted',
        destructive ? 'text-destructive hover:bg-danger-bg' : 'text-foreground',
        className
      )}
      onClick={(e) => {
        onClick?.(e)
        setOpen(false)
      }}
      {...props}
    >
      {children}
    </button>
  )
}

function DropdownMenuLabel({ className, children }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div className={cn('px-[9px] py-1.5 text-xs font-bold text-muted-foreground', className)}>
      {children}
    </div>
  )
}

function DropdownMenuSeparator({ className }: HTMLAttributes<HTMLDivElement>) {
  return <div className={cn('my-[5px] h-px bg-border', className)} />
}

export {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
}
