import { useState, type ReactNode } from 'react'
import { cn } from '@/lib/utils'

interface Tab {
  key: string
  label: string
}

interface TabsProps {
  tabs: Tab[]
  activeTab?: string
  defaultTab?: string
  children: (activeTab: string) => ReactNode
}

export function Tabs({ tabs, activeTab, defaultTab, children }: TabsProps) {
  const [internalActive, setInternalActive] = useState(defaultTab ?? tabs[0]?.key)
  const active = activeTab ?? internalActive

  return (
    <div>
      <div className="inline-flex gap-1 p-1 mb-6 rounded-lg bg-muted">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            type="button"
            onClick={() => setInternalActive(tab.key)}
            className={cn(
              'px-3 py-1.5 rounded-md text-sm font-medium transition-colors',
              active === tab.key
                ? 'bg-card text-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            {tab.label}
          </button>
        ))}
      </div>
      {children(active)}
    </div>
  )
}
