import { cn } from '@/lib/utils'

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
}

export function Button({
  className,
  variant = 'primary',
  size = 'md',
  children,
  ...props
}: ButtonProps) {
  return (
    <button
      className={cn(
        'inline-flex items-center justify-center font-medium rounded-sm transition-all duration-200 cursor-pointer',
        {
          'bg-terreiro text-parchment hover:bg-terreiro-light shadow-sm hover:shadow-md':
            variant === 'primary',
          'bg-gold text-ink hover:brightness-110':
            variant === 'secondary',
          'border border-muted/40 text-parchment hover:border-terreiro hover:text-terreiro':
            variant === 'outline',
          'text-muted hover:text-parchment': variant === 'ghost',
        },
        {
          'px-4 py-2 text-sm': size === 'sm',
          'px-6 py-3 text-base': size === 'md',
          'px-8 py-4 text-lg': size === 'lg',
        },
        className,
      )}
      {...props}
    >
      {children}
    </button>
  )
}
