import { useState, useEffect, useCallback } from 'react'
import { getDb } from '../db/schema'
import { enqueue } from '../sync/engine'
import { apiRequest } from '../api/client'

interface Operation {
  id: string
  plot_id: string
  type: string
  date: string
  responsible: string
  product_used: string
  quantity: number
  cost: number
  notes: string
  synced: number
}

interface Plot {
  id: string
  farm_id: string
  name: string
}

export function useOfflineOperations() {
  const [operations, setOperations] = useState<Operation[]>([])
  const [plots, setPlots] = useState<Plot[]>([])
  const [loading, setLoading] = useState(true)

  const load = useCallback(async () => {
    try {
      const db = await getDb()
      const ops = await db.getAllAsync<Operation>(
        'SELECT * FROM operations ORDER BY date DESC'
      )
      setOperations(ops)

      const p = await db.getAllAsync<Plot>(
        'SELECT * FROM plots ORDER BY name'
      )
      setPlots(p)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load() }, [load])

  const refreshPlots = useCallback(async () => {
    try {
      const data = await apiRequest<Plot[]>('/plots')
      const db = await getDb()
      for (const plot of data) {
        await db.runAsync(
          'INSERT OR REPLACE INTO plots (id, farm_id, name, synced) VALUES (?, ?, ?, 1)',
          plot.id, plot.farm_id || '', plot.name
        )
      }
    } catch { /* offline, use cached */ }
  }, [])

  const createOperation = useCallback(async (
    plotId: string, type: string, date: string,
    responsible: string, productUsed: string,
    quantity: number, cost: number, notes: string
  ) => {
    const payload = { plot_id: plotId, type, date, responsible, product_used: productUsed, quantity, cost, notes }
    await enqueue('operation.created', payload)

    const db = await getDb()
    const id = `local-${Date.now()}`
    await db.runAsync(
      'INSERT INTO operations (id, plot_id, type, date, responsible, product_used, quantity, cost, notes, synced) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0)',
      id, plotId, type, date, responsible, productUsed, quantity, cost, notes
    )

    await load()
  }, [load])

  return { operations, plots, loading, refreshPlots, createOperation }
}
