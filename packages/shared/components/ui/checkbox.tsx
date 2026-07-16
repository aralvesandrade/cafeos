import { forwardRef, type InputHTMLAttributes } from 'react'
import { Check } from 'lucide-react'
import { cn } from '../../lib/utils.ts'

interface CheckboxProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'size'> {
  label?: string
}

const Checkbox = forwardRef<HTMLInputElement, CheckboxProps>(
  ({ className, checked, disabled, label, ...props }, ref) => {
    return (
      <label
        className={cn(
          'inline-flex items-center gap-2 select-none text-sm',
          disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer',
          className
        )}
      >
        <span className="relative inline-block h-[18px] w-[18px] shrink-0">
          <input
            ref={ref}
            type="checkbox"
            checked={checked}
            disabled={disabled}
            className="peer h-[18px] w-[18px] cursor-pointer appearance-none rounded-[5px] border border-input bg-card checked:border-primary checked:bg-primary disabled:cursor-not-allowed"
            {...props}
          />
          <Check
            strokeWidth={3}
            className="pointer-events-none absolute inset-0 m-auto hidden h-3 w-3 text-primary-foreground peer-checked:block"
          />
        </span>
        {label && <span>{label}</span>}
      </label>
    )
  }
)
Checkbox.displayName = 'Checkbox'

export { Checkbox }
