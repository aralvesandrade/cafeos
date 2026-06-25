// Web fallback — uses array storage (no persistence, for dev only)
// Native builds use expo-sqlite

interface Row { [key: string]: unknown }

const tables: Record<string, Row[]> = {
  sync_queue: [],
  operations: [],
  plots: [],
  farms: [],
}

let idCounter = 0

function query(table: string, sql: string, ...params: unknown[]): Row[] {
  // Minimal SQL parser for SELECT/INSERT/UPDATE patterns used in the app
  const rows = tables[table] || []

  if (sql.startsWith('SELECT')) {
    const whereIdx = sql.toUpperCase().indexOf('WHERE')
    let filtered = rows
    if (whereIdx >= 0) {
      const cond = sql.slice(whereIdx + 5).trim()
      const parts = cond.split('=')
      if (parts.length === 2) {
        const key = parts[0].trim()
        const val = params[0]
        filtered = rows.filter((r) => String(r[key]) === String(val))
      }
    }
    const orderIdx = sql.toUpperCase().indexOf('ORDER BY')
    if (orderIdx >= 0) {
      const order = sql.slice(orderIdx + 8).trim()
      const [col, dir] = order.split(' ')
      filtered = [...filtered].sort((a, b) => {
        const cmp = String(a[col] || '').localeCompare(String(b[col] || ''))
        return dir?.toUpperCase() === 'DESC' ? -cmp : cmp
      })
    }
    return filtered
  }

  if (sql.startsWith('INSERT')) {
    idCounter++
    const row: Row = {}
    for (let i = 0; i < params.length; i++) {
      const col = params[i] as string
      i++
      row[col] = params[i]
    }
    tables[table].push(row)
    return []
  }

  if (sql.startsWith('UPDATE')) {
    // Simple update: SET col = ? WHERE col = ?
    // Params: [val, ...conditions]
    const valIdx = sql.indexOf('SET') + 3
    const setPart = sql.slice(valIdx, sql.indexOf('WHERE')).trim()
    const setCol = setPart.split('=')[0].trim()
    const setVal = params[0]

    const wherePart = sql.slice(sql.indexOf('WHERE') + 5).trim()
    const whereCol = wherePart.split('=')[0].trim()
    const whereVal = params[1]

    for (const row of rows) {
      if (String(row[whereCol]) === String(whereVal)) {
        row[setCol] = setVal
      }
    }
    return []
  }

  return []
}

export async function getDb(): Promise<{ getAllAsync: (sql: string, ...params: unknown[]) => Row[]; runAsync: (sql: string, ...params: unknown[]) => void }> {
  return {
    getAllAsync(sql: string, ...params: unknown[]): Row[] {
      const tableMatch = sql.match(/FROM\s+(\w+)/i)
      return query(tableMatch?.[1] || '', sql, ...params)
    },
    runAsync(sql: string, ...params: unknown[]) {
      const tableMatch = sql.match(/INTO\s+(\w+)|UPDATE\s+(\w+)/i)
      query(tableMatch?.[1] || tableMatch?.[2] || '', sql, ...params)
    },
  }
}
