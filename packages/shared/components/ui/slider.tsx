import { forwardRef, type InputHTMLAttributes } from 'react'
import { cn } from '../../lib/utils.ts'

interface SliderProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'size'> {
  minLabel?: string
  maxLabel?: string
}

const Slider = forwardRef<HTMLInputElement, SliderProps>(
  ({ className, minLabel, maxLabel, ...props }, ref) => {
    return (
      <div className="w-full">
        <input
          ref={ref}
          type="range"
          className={cn('w-full accent-primary', className)}
          {...props}
        />
        {(minLabel || maxLabel) && (
          <div className="mt-1 flex justify-between text-xs text-muted-foreground">
            <span>{minLabel}</span>
            <span>{maxLabel}</span>
          </div>
        )}
      </div>
    )
  }
)
Slider.displayName = 'Slider'

export { Slider }
