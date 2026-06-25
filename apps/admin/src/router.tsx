import { createBrowserRouter, Navigate } from 'react-router-dom'
import { AdminLayout } from '@/components/layout/AdminLayout'
import { AuthLayout } from '@/components/layout/AuthLayout'
import { Login } from '@/pages/Login'
import { Dashboard } from '@/pages/Dashboard'
import { Farms } from '@/pages/Farms'
import { FarmDetail } from '@/pages/FarmDetail'
import { Plots } from '@/pages/Plots'
import { PlotDetail } from '@/pages/PlotDetail'
import { Operations } from '@/pages/Operations'
import { Harvests } from '@/pages/Harvests'
import { HarvestDetail } from '@/pages/HarvestDetail'
import { Tenants } from '@/pages/Tenants'
import { Users } from '@/pages/Users'
import { NotFound } from '@/pages/NotFound'

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <AuthLayout />,
    children: [{ index: true, element: <Login /> }],
  },
  {
    path: '/',
    element: <AdminLayout />,
    children: [
      { index: true, element: <Dashboard /> },
      { path: 'farms', element: <Farms /> },
      { path: 'farms/:farmId', element: <FarmDetail /> },
      { path: 'plots', element: <Plots /> },
      { path: 'plots/:plotId', element: <PlotDetail /> },
      { path: 'operations', element: <Operations /> },
      { path: 'harvests', element: <Harvests /> },
      { path: 'harvests/:harvestId', element: <HarvestDetail /> },
      { path: 'tenants', element: <Tenants /> },
      { path: 'users', element: <Users /> },
      { path: '404', element: <NotFound /> },
      { path: '*', element: <Navigate to="/404" replace /> },
    ],
  },
])
