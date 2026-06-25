import { useState, useEffect, useCallback } from 'react'
import { View, Text, FlatList, TouchableOpacity, StyleSheet, ActivityIndicator } from 'react-native'
import { getDb } from '../db/schema'
import { syncAll } from '../sync/engine'

interface SyncItem {
  id: string
  event_type: string
  status: string
  retry_count: number
  created_at: string
}

export function PendingSyncScreen() {
  const [items, setItems] = useState<SyncItem[]>([])
  const [loading, setLoading] = useState(true)
  const [syncing, setSyncing] = useState(false)

  const load = useCallback(async () => {
    try {
      const db = await getDb()
      const rows = await db.getAllAsync<SyncItem>(
        'SELECT * FROM sync_queue ORDER BY created_at DESC'
      )
      setItems(rows)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load() }, [load])

  const handleSync = async () => {
    setSyncing(true)
    try {
      const result = await syncAll()
      await load()
    } finally {
      setSyncing(false)
    }
  }

  const statusLabel: Record<string, string> = { pending: 'Pendente', synced: 'Sincronizado', failed: 'Falhou' }
  const statusColor: Record<string, string> = { pending: '#F57C00', synced: '#2E7D32', failed: '#C62828' }

  if (loading) return <View style={styles.center}><ActivityIndicator size="large" color="#2E7D32" /></View>

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Pendências</Text>
        <TouchableOpacity style={styles.syncButton} onPress={handleSync} disabled={syncing}>
          {syncing ? (
            <ActivityIndicator color="#fff" size="small" />
          ) : (
            <Text style={styles.syncButtonText}>Sincronizar</Text>
          )}
        </TouchableOpacity>
      </View>

      <FlatList
        data={items}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <View style={styles.card}>
            <View style={styles.cardRow}>
              <Text style={styles.eventType}>{item.event_type}</Text>
              <View style={[styles.statusBadge, { backgroundColor: statusColor[item.status] + '20' }]}>
                <Text style={[styles.statusText, { color: statusColor[item.status] }]}>
                  {statusLabel[item.status]}
                </Text>
              </View>
            </View>
            <Text style={styles.cardMeta}>ID: {item.id.slice(0, 16)}...</Text>
            {item.retry_count > 0 && <Text style={styles.cardMeta}>Tentativas: {item.retry_count}</Text>}
            <Text style={styles.cardMeta}>{item.created_at}</Text>
          </View>
        )}
        ListEmptyComponent={<Text style={styles.empty}>Nenhum item pendente</Text>}
      />
    </View>
  )
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#F5F0EB', padding: 16 },
  center: { flex: 1, justifyContent: 'center', alignItems: 'center' },
  header: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 },
  title: { fontSize: 24, fontWeight: 'bold', color: '#2E7D32' },
  syncButton: { backgroundColor: '#2E7D32', paddingHorizontal: 16, paddingVertical: 8, borderRadius: 8, minWidth: 100, alignItems: 'center' },
  syncButtonText: { color: '#fff', fontWeight: '600' },
  card: { backgroundColor: '#fff', borderRadius: 12, padding: 14, marginBottom: 8 },
  cardRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' },
  eventType: { fontSize: 15, fontWeight: '600', color: '#333' },
  statusBadge: { paddingHorizontal: 10, paddingVertical: 3, borderRadius: 12 },
  statusText: { fontSize: 12, fontWeight: '600' },
  cardMeta: { fontSize: 12, color: '#999', marginTop: 2 },
  empty: { textAlign: 'center', color: '#999', marginTop: 40, fontSize: 15 },
})
