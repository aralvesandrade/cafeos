import { useState, useEffect, useCallback, useRef } from 'react'
import NetInfo from '@react-native-community/netinfo'
import { syncAll } from '../sync/engine'

export function useNetworkStatus() {
  const [isConnected, setIsConnected] = useState(true)
  const lastOnline = useRef(false)

  const triggerSync = useCallback(async () => {
    try {
      const result = await syncAll()
      if (result.synced > 0 || result.failed > 0) {
        console.log(`[SYNC] synced=${result.synced} failed=${result.failed}`)
      }
    } catch (err) {
      console.error('[SYNC] error:', err)
    }
  }, [])

  useEffect(() => {
    const unsub = NetInfo.addEventListener((state) => {
      const online = state.isConnected ?? false
      setIsConnected(online)

      if (online && !lastOnline.current) {
        console.log('[NET] online, triggering sync')
        triggerSync()
      }
      lastOnline.current = online
    })

    return () => unsub()
  }, [triggerSync])

  useEffect(() => {
    if (isConnected) triggerSync()
  }, [isConnected, triggerSync])

  return { isConnected }
}
