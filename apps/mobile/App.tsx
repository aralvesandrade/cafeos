import { useState, useEffect } from 'react'
import { StatusBar } from 'expo-status-bar'
import { getToken } from './src/api/client'
import { LoginScreen } from './src/screens/LoginScreen'
import { AppNavigator } from './src/navigation/AppNavigator'
import { useNetworkStatus } from './src/hooks/useNetworkStatus'

function AppContent() {
  const [authenticated, setAuthenticated] = useState(false)
  const [checking, setChecking] = useState(true)
  const { isConnected } = useNetworkStatus()

  useEffect(() => {
    getToken().then((token) => {
      setAuthenticated(!!token)
      setChecking(false)
    })
  }, [])

  if (checking) return null

  if (!authenticated) {
    return <LoginScreen onLogin={() => setAuthenticated(true)} />
  }

  return <AppNavigator />
}

export default function App() {
  return (
    <>
      <StatusBar style="auto" />
      <AppContent />
    </>
  )
}
