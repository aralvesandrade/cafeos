import { getDb } from '../db/schema'
import { apiRequest, getToken } from '../api/client'

export interface SyncItem {
  id: string
  event_type: string
  payload: string
  client_timestamp: string
  status: string
  retry_count: number
}

async function getPendingItems(): Promise<SyncItem[]> {
  const db = await getDb()
  const rows = await db.getAllAsync<SyncItem>(
    'SELECT * FROM sync_queue WHERE status = ? ORDER BY created_at ASC LIMIT 50',
    'pending'
  )
  return rows
}

async function markSynced(id: string): Promise<void> {
  const db = await getDb()
  await db.runAsync('UPDATE sync_queue SET status = ? WHERE id = ?', 'synced', id)
}

async function markFailed(id: string): Promise<void> {
  const db = await getDb()
  await db.runAsync('UPDATE sync_queue SET status = ?, retry_count = retry_count + 1 WHERE id = ?', 'failed', id)
}

async function incrementRetry(id: string): Promise<void> {
  const db = await getDb()
  await db.runAsync(
    'UPDATE sync_queue SET retry_count = retry_count + 1, status = ? WHERE id = ?',
    'pending',
    id
  )
}

export async function enqueue(eventType: string, payload: unknown): Promise<void> {
  const db = await getDb()
  const id = `${Date.now()}-${Math.random().toString(36).slice(2, 9)}`
  await db.runAsync(
    'INSERT INTO sync_queue (id, event_type, payload, client_timestamp) VALUES (?, ?, ?, ?)',
    id,
    eventType,
    JSON.stringify(payload),
    new Date().toISOString()
  )
}

export async function syncAll(): Promise<{ synced: number; failed: number }> {
  const token = await getToken()
  if (!token) return { synced: 0, failed: 0 }

  const items = await getPendingItems()
  if (items.length === 0) return { synced: 0, failed: 0 }

  const batch = items.map((item) => ({
    client_id: item.id,
    event_type: item.event_type,
    payload: JSON.parse(item.payload),
    client_timestamp: item.client_timestamp,
  }))

  try {
    const resp = await apiRequest<{ accepted: number; errors: Array<{ client_id: string }> }>(
      '/sync', { method: 'POST', body: { batch } }
    )

    const failedIds = new Set(resp.errors?.map((e) => e.client_id) || [])

    for (const item of items) {
      if (failedIds.has(item.id)) {
        if (item.retry_count >= 2) await markFailed(item.id)
        else await incrementRetry(item.id)
      } else {
        await markSynced(item.id)
      }
    }

    return { synced: resp.accepted || 0, failed: failedIds.size }
  } catch {
    for (const item of items) {
      if (item.retry_count >= 2) await markFailed(item.id)
      else await incrementRetry(item.id)
    }
    return { synced: 0, failed: items.length }
  }
}
