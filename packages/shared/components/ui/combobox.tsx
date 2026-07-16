import { useMemo, useState } from 'react'
import { Check, ChevronDown } from 'lucide-react'
import { cn } from '../../lib/utils.ts'

interface ComboboxOption {
  value: string
  label: string
}

interface ComboboxProps {
  options: ComboboxOption[]
  value?: string
  onValueChange: (value: string) => void
  placeholder?: string
  searchPlaceholder?: string
  emptyMessage?: string
  className?: string
  disabled?: boolean
}

function Combobox({
  options,
  value,
  onValueChange,
  placeholder = 'Selecione…',
  searchPlaceholder = 'Buscar…',
  emptyMessage = 'Nenhum resultado.',
  className,
  disabled,
}: ComboboxProps) {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState('')

  const selected = options.find((o) => o.value === value)
  const filtered = useMemo(
    () => options.filter((o) => o.label.toLowerCase().includes(query.toLowerCase())),
    [options, query]
  )

  return (
    <div className={cn('relative', className)}>
      <button
        type="button"
        disabled={disabled}
        onClick={() => setOpen(!open)}
        className={cn(
          'flex h-[38px] w-full items-center justify-between rounded-lg border border-input bg-card px-3 text-left text-sm',
          disabled && 'cursor-not-allowed bg-muted text-muted-foreground'
        )}
      >
        <span className={cn(!selected && 'text-muted-foreground')}>
          {selected ? selected.label : placeholder}
        </span>
        <ChevronDown className="h-4 w-4 shrink-0 text-muted-foreground" />
      </button>

      {open && (
        <>
          <div className="fixed inset-0 z-20" onClick={() => setOpen(false)} />
          <div className="animate-ds-pop absolute left-0 top-[calc(100%+6px)] z-30 max-h-[200px] w-full overflow-auto rounded-lg border border-border bg-card p-[5px] shadow-ds-md">
            <input
              autoFocus
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder={searchPlaceholder}
              className="mb-1 w-full rounded-md border border-input bg-background px-2.5 py-1.5 text-sm outline-none focus:border-ring"
            />
            {filtered.length === 0 && (
              <div className="px-2.5 py-2 text-sm text-muted-foreground">{emptyMessage}</div>
            )}
            {filtered.map((option) => (
              <button
                key={option.value}
                type="button"
                onClick={() => {
                  onValueChange(option.value)
                  setOpen(false)
                  setQuery('')
                }}
                className={cn(
                  'flex w-full items-center justify-between rounded-md px-2.5 py-2 text-left text-sm hover:bg-muted',
                  option.value === value && 'bg-muted font-semibold'
                )}
              >
                {option.label}
                {option.value === value && <Check className="h-3.5 w-3.5" />}
              </button>
            ))}
          </div>
        </>
      )}
    </div>
  )
}

export { Combobox, type ComboboxOption }
