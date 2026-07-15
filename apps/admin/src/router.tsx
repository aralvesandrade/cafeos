import { createBrowserRouter, Navigate } from 'react-router-dom'
import { AdminLayout } from '@/components/layout/AdminLayout'
import { AuthLayout } from '@/components/layout/AuthLayout'
import { Login } from '@/pages/Login'
import { Dashboard } from '@/pages/Dashboard'
import { Farms } from '@/pages/Farms'
import { FarmDetail } from '@/pages/FarmDetail'
import { FarmEdit } from '@/pages/FarmEdit'
import { Plots } from '@/pages/Plots'
import { PlotDetail } from '@/pages/PlotDetail'
import { PlotEdit } from '@/pages/PlotEdit'
import { Operations } from '@/pages/Operations'
import { OperationDetail } from '@/pages/OperationDetail'
import { OperationTypes } from '@/pages/OperationTypes'
import { Harvests } from '@/pages/Harvests'
import { HarvestDetail } from '@/pages/HarvestDetail'
import { Organizations } from '@/pages/Organizations'
import { Users } from '@/pages/Users'
import { Financial } from '@/pages/Financial'
import { CostCenters } from '@/pages/CostCenters'
import { Budget } from '@/pages/Budget'
import { CostAllocations } from '@/pages/CostAllocations'
import { Stock } from '@/pages/Stock'
import { Fleet } from '@/pages/Fleet'
import { Labor } from '@/pages/Labor'
import { TeamUsers } from '@/pages/TeamUsers'
import { Permissions } from '@/pages/Permissions'
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
      { path: 'farms/new', element: <FarmEdit /> },
      { path: 'farms/:farmId', element: <FarmDetail /> },
      { path: 'farms/:farmId/edit', element: <FarmEdit /> },
      { path: 'plots', element: <Plots /> },
      { path: 'plots/new', element: <PlotEdit /> },
      { path: 'plots/:plotId', element: <PlotDetail /> },
      { path: 'plots/:plotId/edit', element: <PlotEdit /> },
      { path: 'operations', element: <Operations /> },
      { path: 'operations/:operationId', element: <OperationDetail /> },
      { path: 'operation-types', element: <OperationTypes /> },
      { path: 'harvests', element: <Harvests /> },
      { path: 'harvests/:harvestId', element: <HarvestDetail /> },
      { path: 'financial', element: <Financial /> },
      { path: 'cost-centers', element: <CostCenters /> },
      { path: 'budgets', element: <Budget /> },
      { path: 'cost-allocations', element: <CostAllocations /> },
      { path: 'stock', element: <Stock /> },
      { path: 'fleet', element: <Fleet /> },
      { path: 'labor', element: <Labor /> },
      { path: 'team', element: <TeamUsers /> },
      { path: 'organizations', element: <Organizations /> },
      { path: 'users', element: <Users /> },
      { path: 'permissions', element: <Permissions /> },
      { path: '404', element: <NotFound /> },
      { path: '*', element: <Navigate to="/404" replace /> },
    ],
  },
])
