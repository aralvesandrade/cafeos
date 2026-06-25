import { Platform } from 'react-native'

const isWeb = Platform.OS === 'web'

async function getItem(key: string): Promise<string | null> {
  if (isWeb) {
    return localStorage.getItem(key)
  }
  const SecureStore = await import('expo-secure-store')
  return SecureStore.getItemAsync(key)
}

async function setItem(key: string, value: string): Promise<void> {
  if (isWeb) {
    localStorage.setItem(key, value)
    return
  }
  const SecureStore = await import('expo-secure-store')
  await SecureStore.setItemAsync(key, value)
}

async function removeItem(key: string): Promise<void> {
  if (isWeb) {
    localStorage.removeItem(key)
    return
  }
  const SecureStore = await import('expo-secure-store')
  await SecureStore.deleteItemAsync(key)
}

export const storage = { getItem, setItem, removeItem }
