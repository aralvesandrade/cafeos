import { createContext, forwardRef, useContext, type InputHTMLAttributes, type ReactNode } from 'react'
import { cn } from '../../lib/utils.ts'

interface RadioGroupContextValue {
  name: string
  value?: string
  onValueChange?: (value: string) => void
}

const RadioGroupContext = createContext<RadioGroupContextValue | null>(null)

interface RadioGroupProps {
  name: string
  value?: string
  onValueChange?: (value: string) => void
  className?: string
  children: ReactNode
}

function RadioGroup({ name, value, onValueChange, className, children }: RadioGroupProps) {
  return (
    <RadioGroupContext.Provider value={{ name, value, onValueChange }}>
      <div role="radiogroup" className={cn('flex flex-col gap-2', className)}>
        {children}
      </div>
    </RadioGroupContext.Provider>
  )
}

interface RadioGroupItemProps
  extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'size' | 'name' | 'onChange'> {
  value: string
  label?: string
}

const RadioGroupItem = forwardRef<HTMLInputElement, RadioGroupItemProps>(
  ({ className, value, label, disabled, ...props }, ref) => {
    const ctx = useContext(RadioGroupContext)
    if (!ctx) throw new Error('RadioGroupItem must be used within a RadioGroup')

    return (
      <label
        className={cn(
          'inline-flex items-center gap-2 select-none',
          disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer',
          className
        )}
      >
        <span className="relative inline-block h-[18px] w-[18px] shrink-0">
          <input
            ref={ref}
            type="radio"
            name={ctx.name}
            value={value}
            checked={ctx.value === value}
            disabled={disabled}
            onChange={() => ctx.onValueChange?.(value)}
            className="peer h-[18px] w-[18px] cursor-pointer appearance-none rounded-full border border-input bg-card checked:border-primary disabled:cursor-not-allowed"
            {...props}
          />
          <span className="pointer-events-none absolute inset-0 m-auto hidden h-[9px] w-[9px] rounded-full bg-primary peer-checked:block" />
        </span>
        {label && <span className="text-sm">{label}</span>}
      </label>
    )
  }
)
RadioGroupItem.displayName = 'RadioGroupItem'

export { RadioGroup, RadioGroupItem }
