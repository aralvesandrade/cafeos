import { useState } from 'react'
import { View, Text, FlatList, TouchableOpacity, TextInput, StyleSheet, Modal, Alert, ActivityIndicator } from 'react-native'
import { useOfflineOperations } from '../hooks/useOfflineOperations'
import { Ionicons } from '@expo/vector-icons'

const operationTypes = ['adubacao', 'pulverizacao', 'irrigacao', 'poda', 'colheita']
const typeLabels: Record<string, string> = {
  adubacao: 'Adubação', pulverizacao: 'Pulverização', irrigacao: 'Irrigação', poda: 'Poda', colheita: 'Colheita',
}

export function OperationsScreen() {
  const { operations, plots, loading, createOperation } = useOfflineOperations()
  const [modalVisible, setModalVisible] = useState(false)
  const [form, setForm] = useState({ plot_id: '', type: 'adubacao', date: '', responsible: '', product_used: '', quantity: '', cost: '', notes: '' })
  const [saving, setSaving] = useState(false)

  const handleSave = async () => {
    if (!form.plot_id && plots.length > 0) { Alert.alert('Erro', 'Selecione um talhão'); return }
    setSaving(true)
    try {
      await createOperation(
        form.plot_id || plots[0]?.id || '',
        form.type,
        form.date || new Date().toISOString().split('T')[0],
        form.responsible,
        form.product_used,
        parseFloat(form.quantity) || 0,
        parseFloat(form.cost) || 0,
        form.notes
      )
      setModalVisible(false)
      setForm({ plot_id: '', type: 'adubacao', date: '', responsible: '', product_used: '', quantity: '', cost: '', notes: '' })
    } catch (err) {
      Alert.alert('Erro', 'Falha ao salvar')
    } finally {
      setSaving(false)
    }
  }

  const getTypeColor = (type: string) => {
    const colors: Record<string, string> = { adubacao: '#1976D2', pulverizacao: '#F57C00', irrigacao: '#388E3C', poda: '#6D4C41', colheita: '#C62828' }
    return colors[type] || '#666'
  }

  if (loading) return <View style={styles.center}><ActivityIndicator size="large" color="#2E7D32" /></View>

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Operações</Text>
        <TouchableOpacity style={styles.addButton} onPress={() => setModalVisible(true)}>
          <Text style={styles.addButtonText}>+ Nova</Text>
        </TouchableOpacity>
      </View>

      <FlatList
        data={operations}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <View style={styles.card}>
            <View style={styles.cardHeader}>
              <View style={[styles.typeBadge, { backgroundColor: getTypeColor(item.type) }]}>
                <Text style={styles.typeText}>{typeLabels[item.type] || item.type}</Text>
              </View>
              {!item.synced && <Text style={styles.pendingBadge}>Pendente</Text>}
            </View>
            <Text style={styles.cardDate}>{item.date}</Text>
            <Text style={styles.cardDetail}>Talhão: {item.plot_id.slice(0, 8)}...</Text>
            {item.responsible ? <Text style={styles.cardDetail}>Resp: {item.responsible}</Text> : null}
            {item.cost > 0 ? <Text style={styles.cardCost}>R$ {item.cost.toFixed(2)}</Text> : null}
          </View>
        )}
        ListEmptyComponent={<Text style={styles.empty}>Nenhuma operação registrada</Text>}
      />

      <Modal visible={modalVisible} animationType="slide" transparent>
        <View style={styles.modalOverlay}>
          <View style={styles.modal}>
            <Text style={styles.modalTitle}>Nova Operação</Text>

            <Text style={styles.label}>Tipo</Text>
            <View style={styles.typeRow}>
              {operationTypes.map((t) => (
                <TouchableOpacity
                  key={t}
                  style={[styles.typeChip, form.type === t && styles.typeChipActive]}
                  onPress={() => setForm({ ...form, type: t })}
                >
                  <Text style={[styles.typeChipText, form.type === t && styles.typeChipTextActive]}>
                    {typeLabels[t]}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>

            <Text style={styles.label}>Data</Text>
            <TextInput style={styles.input} value={form.date} onChangeText={(d) => setForm({ ...form, date: d })} placeholder="YYYY-MM-DD" />

            <Text style={styles.label}>Responsável</Text>
            <TextInput style={styles.input} value={form.responsible} onChangeText={(v) => setForm({ ...form, responsible: v })} placeholder="Nome" />

            <View style={styles.row}>
              <View style={{ flex: 1, marginRight: 8 }}>
                <Text style={styles.label}>Quantidade</Text>
                <TextInput style={styles.input} value={form.quantity} onChangeText={(v) => setForm({ ...form, quantity: v })} keyboardType="decimal-pad" placeholder="0" />
              </View>
              <View style={{ flex: 1 }}>
                <Text style={styles.label}>Custo (R$)</Text>
                <TextInput style={styles.input} value={form.cost} onChangeText={(v) => setForm({ ...form, cost: v })} keyboardType="decimal-pad" placeholder="0" />
              </View>
            </View>

            <View style={styles.modalButtons}>
              <TouchableOpacity style={styles.cancelButton} onPress={() => setModalVisible(false)}>
                <Text style={styles.cancelButtonText}>Cancelar</Text>
              </TouchableOpacity>
              <TouchableOpacity style={styles.saveButton} onPress={handleSave} disabled={saving}>
                {saving ? <ActivityIndicator color="#fff" /> : <Text style={styles.saveButtonText}>Salvar</Text>}
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </Modal>
    </View>
  )
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#F5F0EB', padding: 16 },
  center: { flex: 1, justifyContent: 'center', alignItems: 'center' },
  header: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 },
  title: { fontSize: 24, fontWeight: 'bold', color: '#2E7D32' },
  addButton: { backgroundColor: '#2E7D32', paddingHorizontal: 16, paddingVertical: 8, borderRadius: 8 },
  addButtonText: { color: '#fff', fontWeight: '600' },
  card: { backgroundColor: '#fff', borderRadius: 12, padding: 14, marginBottom: 10, shadowColor: '#000', shadowOpacity: 0.05, shadowRadius: 5, elevation: 2 },
  cardHeader: { flexDirection: 'row', justifyContent: 'space-between', marginBottom: 8 },
  typeBadge: { paddingHorizontal: 10, paddingVertical: 3, borderRadius: 12 },
  typeText: { color: '#fff', fontSize: 12, fontWeight: '600' },
  pendingBadge: { fontSize: 11, color: '#F57C00', fontWeight: '600' },
  cardDate: { fontSize: 13, color: '#666', marginBottom: 2 },
  cardDetail: { fontSize: 13, color: '#333' },
  cardCost: { fontSize: 15, fontWeight: '700', color: '#2E7D32', marginTop: 4 },
  empty: { textAlign: 'center', color: '#999', marginTop: 40, fontSize: 15 },
  modalOverlay: { flex: 1, justifyContent: 'flex-end', backgroundColor: 'rgba(0,0,0,0.5)' },
  modal: { backgroundColor: '#fff', borderTopLeftRadius: 20, borderTopRightRadius: 20, padding: 20, maxHeight: '80%' },
  modalTitle: { fontSize: 20, fontWeight: 'bold', color: '#2E7D32', marginBottom: 16 },
  label: { fontSize: 13, fontWeight: '600', color: '#555', marginTop: 10, marginBottom: 4 },
  input: { borderWidth: 1, borderColor: '#ddd', borderRadius: 8, padding: 10, fontSize: 14 },
  row: { flexDirection: 'row' },
  typeRow: { flexDirection: 'row', flexWrap: 'wrap', gap: 6 },
  typeChip: { paddingHorizontal: 12, paddingVertical: 6, borderRadius: 16, borderWidth: 1, borderColor: '#ddd' },
  typeChipActive: { backgroundColor: '#2E7D32', borderColor: '#2E7D32' },
  typeChipText: { fontSize: 12, color: '#666' },
  typeChipTextActive: { color: '#fff', fontWeight: '600' },
  modalButtons: { flexDirection: 'row', justifyContent: 'flex-end', marginTop: 20, gap: 10 },
  cancelButton: { paddingHorizontal: 20, paddingVertical: 12, borderRadius: 8 },
  cancelButtonText: { color: '#666', fontWeight: '600' },
  saveButton: { backgroundColor: '#2E7D32', paddingHorizontal: 24, paddingVertical: 12, borderRadius: 8, minWidth: 80, alignItems: 'center' },
  saveButtonText: { color: '#fff', fontWeight: '600' },
})
