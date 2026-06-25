import * as SQLite from 'expo-sqlite'

let db: SQLite.SQLiteDatabase | null = null

export async function getDb(): Promise<SQLite.SQLiteDatabase> {
  if (db) return db
  db = await SQLite.openDatabaseAsync('cafeos_offline.db')
  await migrate(db)
  return db
}

async function migrate(db: SQLite.SQLiteDatabase): Promise<void> {
  await db.execAsync(`
    CREATE TABLE IF NOT EXISTS sync_queue (
      id TEXT PRIMARY KEY,
      event_type TEXT NOT NULL,
      payload TEXT NOT NULL,
      client_timestamp TEXT NOT NULL,
      status TEXT DEFAULT 'pending',
      retry_count INTEGER DEFAULT 0,
      created_at TEXT DEFAULT (datetime('now'))
    );

    CREATE TABLE IF NOT EXISTS operations (
      id TEXT PRIMARY KEY,
      plot_id TEXT NOT NULL,
      type TEXT NOT NULL,
      date TEXT NOT NULL,
      responsible TEXT DEFAULT '',
      product_used TEXT DEFAULT '',
      quantity REAL DEFAULT 0,
      cost REAL DEFAULT 0,
      notes TEXT DEFAULT '',
      synced INTEGER DEFAULT 0,
      created_at TEXT DEFAULT (datetime('now'))
    );

    CREATE TABLE IF NOT EXISTS plots (
      id TEXT PRIMARY KEY,
      farm_id TEXT NOT NULL,
      name TEXT NOT NULL,
      area_ha REAL DEFAULT 0,
      cultivar TEXT DEFAULT '',
      synced INTEGER DEFAULT 0
    );

    CREATE TABLE IF NOT EXISTS farms (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      owner TEXT DEFAULT '',
      location TEXT DEFAULT '',
      synced INTEGER DEFAULT 0
    );
  `)
}
