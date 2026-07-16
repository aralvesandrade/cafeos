import { createContext, useContext, useState, type ReactNode } from 'react'
import { ChevronDown } from 'lucide-react'
import { cn } from '../../lib/utils.ts'

interface AccordionContextValue {
  openItems: Set<string>
  toggle: (value: string) => void
}

const AccordionContext = createContext<AccordionContextValue | null>(null)

interface AccordionProps {
  className?: string
  children: ReactNode
  type?: 'single' | 'multiple'
  defaultValue?: string[]
}

function Accordion({ className, children, type = 'single', defaultValue = [] }: AccordionProps) {
  const [openItems, setOpenItems] = useState<Set<string>>(new Set(defaultValue))

  const toggle = (value: string) => {
    setOpenItems((prev) => {
      const next = type === 'single' ? new Set<string>() : new Set(prev)
      if (prev.has(value)) {
        if (type === 'single') return new Set()
        next.delete(value)
      } else {
        next.add(value)
      }
      return next
    })
  }

  return (
    <AccordionContext.Provider value={{ openItems, toggle }}>
      <div className={cn('rounded-xl border border-border bg-card', className)}>{children}</div>
    </AccordionContext.Provider>
  )
}

interface AccordionItemProps {
  value: string
  title: ReactNode
  children: ReactNode
  className?: string
  last?: boolean
}

function AccordionItem({ value, title, children, className, last }: AccordionItemProps) {
  const ctx = useContext(AccordionContext)
  if (!ctx) throw new Error('AccordionItem must be used within an Accordion')
  const open = ctx.openItems.has(value)

  return (
    <div className={cn(!last && 'border-b border-border', className)}>
      <button
        type="button"
        className="flex w-full items-center justify-between px-5 py-4 text-left text-sm font-semibold"
        onClick={() => ctx.toggle(value)}
        aria-expanded={open}
      >
        {title}
        <ChevronDown
          className={cn('h-4 w-4 shrink-0 transition-transform duration-150', open && 'rotate-180')}
        />
      </button>
      {open && (
        <div className="px-5 pb-4 text-[13px] leading-relaxed text-muted-foreground">{children}</div>
      )}
    </div>
  )
}

export { Accordion, AccordionItem }
