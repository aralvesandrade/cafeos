const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:5001'

export interface PlanFeature {
  label: string
  included: boolean
}

export interface Plan {
  id: string
  name: string
  slug: string
  description: string
  price_cents: number
  billing_interval: string
  max_farms: number
  max_users: number
  features: PlanFeature[]
  active: boolean
  featured: boolean
  display_order: number
}

export async function fetchPlans(): Promise<Plan[]> {
  const res = await fetch(`${API_BASE}/api/v1/public/plans`)
  if (!res.ok) throw new Error('failed to fetch plans')
  const plans: Plan[] = await res.json()
  return [...plans].sort((a, b) => a.display_order - b.display_order)
}
