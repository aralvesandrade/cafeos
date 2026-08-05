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

export interface RegisterPayload {
  name: string
  email: string
  password: string
  plan_slug?: string
}

export interface RegisterResponse {
  user_id: string
}

export async function registerSignup(payload: RegisterPayload): Promise<RegisterResponse> {
  const res = await fetch(`${API_BASE}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error || 'falha ao realizar cadastro')
  }
  return res.json()
}
