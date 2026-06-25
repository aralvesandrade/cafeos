import { Platform } from 'react-native'
import { storage } from './storage'

const defaultHost = Platform.OS === 'android' ? '10.0.2.2' : 'localhost'
const API_BASE = process.env.EXPO_PUBLIC_API_URL || `http://${defaultHost}:5001`

export interface RequestOptions {
  method?: string
  body?: unknown
  params?: Record<string, string>
}

export async function getToken(): Promise<string | null> {
  return storage.getItem('cafeos_token')
}

export async function setToken(token: string): Promise<void> {
  await storage.setItem('cafeos_token', token)
}

export async function clearToken(): Promise<void> {
  await storage.removeItem('cafeos_token')
}

export async function getTenantId(): Promise<string | null> {
  return storage.getItem('cafeos_tenant_id')
}

export async function setTenantId(id: string): Promise<void> {
  await storage.setItem('cafeos_tenant_id', id)
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const token = await getToken()
  const tenantId = await getTenantId()

  const isPublic = path.startsWith('/auth/')
  let url: string
  if (isPublic) {
    url = `${API_BASE}${path}`
  } else {
    url = `${API_BASE}/api/v1${tenantId ? `/${tenantId}` : ''}${path}`
  }

  if (options.params) {
    const searchParams = new URLSearchParams(options.params)
    url += `?${searchParams.toString()}`
  }

  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (token) headers['Authorization'] = `Bearer ${token}`

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

export async function loginRequest(email: string, password: string): Promise<{
  token: string; tenant_id: string; user: { id: string; email: string; name: string; role: string }
}> {
  return apiRequest('/auth/login', {
    method: 'POST',
    body: { email, password },
  })
}
