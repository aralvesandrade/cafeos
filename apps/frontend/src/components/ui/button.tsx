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
        'inline-flex items-center justify-center font-medium rounded-lg transition-all duration-200 cursor-pointer',
        {
          'bg-coffee-green text-white hover:bg-coffee-green-dark shadow-sm hover:shadow-md':
            variant === 'primary',
          'bg-coffee-brown text-white hover:bg-coffee-brown-light':
            variant === 'secondary',
          'border-2 border-coffee-green text-coffee-green hover:bg-coffee-green hover:text-white':
            variant === 'outline',
          'text-coffee-text hover:text-coffee-green': variant === 'ghost',
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
