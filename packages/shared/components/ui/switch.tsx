import { forwardRef, type InputHTMLAttributes } from 'react'
import { cn } from '../../lib/utils.ts'

interface SwitchProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'size'> {}

const Switch = forwardRef<HTMLInputElement, SwitchProps>(
  ({ className, checked, disabled, ...props }, ref) => {
    return (
      <label
        className={cn(
          'relative inline-flex h-6 w-[42px] shrink-0 select-none items-center rounded-full p-[3px] transition-colors',
          checked ? 'bg-primary' : 'bg-input',
          disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer',
          className
        )}
      >
        <input
          ref={ref}
          type="checkbox"
          role="switch"
          checked={checked}
          disabled={disabled}
          className="peer absolute inset-0 h-full w-full cursor-pointer appearance-none disabled:cursor-not-allowed"
          {...props}
        />
        <span
          className={cn(
            'pointer-events-none h-[18px] w-[18px] rounded-full bg-card shadow-ds-sm transition-transform',
            checked ? 'translate-x-[18px]' : 'translate-x-0'
          )}
        />
      </label>
    )
  }
)
Switch.displayName = 'Switch'

export { Switch }
