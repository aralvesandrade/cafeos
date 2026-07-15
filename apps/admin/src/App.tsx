import { RouterProvider } from 'react-router-dom'
import { AuthProvider } from '@/lib/auth'
import { PermissionsProvider } from '@/lib/permissions'
import { ThemeProvider } from '@/lib/theme'
import { ToastProvider } from '@/lib/toast'
import { ConfirmProvider } from '@/lib/confirm'
import { router } from '@/router'

export default function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <PermissionsProvider>
          <ToastProvider>
            <ConfirmProvider>
              <RouterProvider router={router} />
            </ConfirmProvider>
          </ToastProvider>
        </PermissionsProvider>
      </AuthProvider>
    </ThemeProvider>
  )
}
