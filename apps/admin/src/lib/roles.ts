export type UserRole =
  | 'platform_owner'
  | 'organization_admin'
  | 'proprietario'
  | 'gerente_agricola'
  | 'engenheiro_agronomo'
  | 'tecnico_agricola'
  | 'operador_campo'
  | 'financeiro'
  | 'consultor_externo'
  | 'auditor'

export const ROLE_LABELS: Record<UserRole, string> = {
  platform_owner: 'Platform Owner',
  organization_admin: 'Admin',
  proprietario: 'Proprietário',
  gerente_agricola: 'Gerente Agrícola',
  engenheiro_agronomo: 'Eng. Agrônomo',
  tecnico_agricola: 'Técnico Agrícola',
  operador_campo: 'Operador de Campo',
  financeiro: 'Financeiro',
  consultor_externo: 'Consultor',
  auditor: 'Auditor',
}

export const ALL_ROLES: UserRole[] = Object.keys(ROLE_LABELS) as UserRole[]
