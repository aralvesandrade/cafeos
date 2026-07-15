import { RouterProvider } from 'react-router-dom'
import { AuthProvider } from '@/lib/auth'
import { PermissionsProvider } from '@/lib/permissions'
import { RolesProvider } from '@/lib/roles'
import { ThemeProvider } from '@/lib/theme'
import { ToastProvider } from '@/lib/toast'
import { ConfirmProvider } from '@/lib/confirm'
import { router } from '@/router'

export default function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <PermissionsProvider>
          <RolesProvider>
            <ToastProvider>
              <ConfirmProvider>
                <RouterProvider router={router} />
              </ConfirmProvider>
            </ToastProvider>
          </RolesProvider>
        </PermissionsProvider>
      </AuthProvider>
    </ThemeProvider>
  )
}
