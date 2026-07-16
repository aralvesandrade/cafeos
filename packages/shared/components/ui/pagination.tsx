import { ChevronLeft, ChevronRight } from 'lucide-react'
import { cn } from '../../lib/utils.ts'

interface PaginationProps {
  page: number
  pageCount: number
  totalItems: number
  pageSize: number
  onPageChange: (page: number) => void
  className?: string
}

function Pagination({ page, pageCount, totalItems, pageSize, onPageChange, className }: PaginationProps) {
  const start = (page - 1) * pageSize + 1
  const end = Math.min(page * pageSize, totalItems)

  const pages = Array.from({ length: pageCount }, (_, i) => i + 1)

  return (
    <div
      className={cn(
        'flex items-center justify-between rounded-xl border border-border bg-card px-5 py-3.5',
        className
      )}
    >
      <span className="text-[13px] text-muted-foreground">
        Mostrando {start}–{end} de {totalItems} registros
      </span>
      <div className="flex items-center gap-1.5">
        <button
          type="button"
          disabled={page <= 1}
          onClick={() => onPageChange(page - 1)}
          className="flex h-8 w-8 items-center justify-center rounded-md border border-border hover:bg-muted disabled:pointer-events-none disabled:opacity-50"
        >
          <ChevronLeft className="h-4 w-4" />
        </button>
        {pages.map((p) => (
          <button
            key={p}
            type="button"
            onClick={() => onPageChange(p)}
            className={cn(
              'flex h-8 w-8 items-center justify-center rounded-md text-sm',
              p === page
                ? 'bg-primary text-primary-foreground'
                : 'border border-border bg-transparent hover:bg-muted'
            )}
          >
            {p}
          </button>
        ))}
        <button
          type="button"
          disabled={page >= pageCount}
          onClick={() => onPageChange(page + 1)}
          className="flex h-8 w-8 items-center justify-center rounded-md border border-border hover:bg-muted disabled:pointer-events-none disabled:opacity-50"
        >
          <ChevronRight className="h-4 w-4" />
        </button>
      </div>
    </div>
  )
}

export { Pagination }
