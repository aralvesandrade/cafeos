const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080'

interface RequestOptions {
  method?: string
  body?: unknown
  params?: Record<string, string>
}

function getToken(): string | null {
  return localStorage.getItem('cafeos_token')
}

function getTenantId(): string | null {
  return localStorage.getItem('cafeos_tenant_id')
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const token = getToken()
  const tenantId = getTenantId()

  let url = `${API_BASE}/api/v1${tenantId ? `/${tenantId}` : ''}${path}`

  if (options.params) {
    const searchParams = new URLSearchParams(options.params)
    url += `?${searchParams.toString()}`
  }

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(url, {
    method: options.method || 'GET',
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: response.statusText }))
    throw new Error(error.message || `Erro ${response.status}`)
  }

  if (response.status === 204) return undefined as T

  return response.json()
}

export function setAuthData(token: string, tenantId: string) {
  localStorage.setItem('cafeos_token', token)
  localStorage.setItem('cafeos_tenant_id', tenantId)
}

export function clearAuthData() {
  localStorage.removeItem('cafeos_token')
  localStorage.removeItem('cafeos_tenant_id')
}
