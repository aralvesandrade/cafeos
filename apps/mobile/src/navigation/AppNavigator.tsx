import { createBottomTabNavigator } from '@react-navigation/bottom-tabs'
import { NavigationContainer } from '@react-navigation/native'
import { OperationsScreen } from '../screens/OperationsScreen'
import { PendingSyncScreen } from '../screens/PendingSyncScreen'
import { View, Text, StyleSheet } from 'react-native'

const Tab = createBottomTabNavigator()

function DashboardScreen() {
  return (
    <View style={styles.center}>
      <Text style={styles.title}>CafeOS</Text>
      <Text style={styles.subtitle}>App offline para operações de campo</Text>
    </View>
  )
}

export function AppNavigator() {
  return (
    <NavigationContainer>
      <Tab.Navigator
        screenOptions={{
          headerStyle: { backgroundColor: '#2E7D32' },
          headerTintColor: '#fff',
          tabBarActiveTintColor: '#2E7D32',
          tabBarInactiveTintColor: '#999',
        }}
      >
        <Tab.Screen
          name="Dashboard"
          component={DashboardScreen}
          options={{ tabBarLabel: 'Início', tabBarIcon: ({ color }) => <Text style={{ color, fontSize: 20 }}>🏠</Text> }}
        />
        <Tab.Screen
          name="Operations"
          component={OperationsScreen}
          options={{ tabBarLabel: 'Operações', tabBarIcon: ({ color }) => <Text style={{ color, fontSize: 20 }}>🚜</Text> }}
        />
        <Tab.Screen
          name="PendingSync"
          component={PendingSyncScreen}
          options={{ tabBarLabel: 'Pendências', tabBarIcon: ({ color }) => <Text style={{ color, fontSize: 20 }}>📤</Text> }}
        />
      </Tab.Navigator>
    </NavigationContainer>
  )
}

const styles = StyleSheet.create({
  center: { flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: '#F5F0EB' },
  title: { fontSize: 28, fontWeight: 'bold', color: '#2E7D32' },
  subtitle: { fontSize: 14, color: '#666', marginTop: 8 },
})
