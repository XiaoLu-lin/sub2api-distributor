export type PortalRole = 'distributor' | 'operator'
export type WithdrawStatus = 'paying' | 'paid' | 'cancelled'
export type UserRole = PortalRole | ''

export interface LoginResponse {
  token: string
  user: {
    id: number
    email: string
    role: string
  }
  portal_role: PortalRole
}

export interface MeResponse {
  user_id: number
  email: string
  portal_role: PortalRole
}

export interface DashboardSummary {
  total_earned: number
  frozen_amount: number
  internal_transferred_amount: number
  paying_amount: number
  paid_amount: number
  withdrawable_amount: number
}

export interface InviteeItem {
  user_id: number
  email: string
  username: string
  total_rebate: number
  created_at: string
}

export interface InviteMeta {
  aff_code: string
}

export interface RebateItem {
  ledger_id: number
  amount: number
  source_user_id: number
  source_email: string
  source_order_id?: number
  created_at: string
}

export interface WithdrawalItem {
  id: number
  request_no: string
  user_id: number
  user_email?: string
  amount: number
  status: WithdrawStatus
  applicant_remark: string
  paid_channel: string
  paid_reference_no: string
  paid_remark: string
  paid_at?: string | null
  created_at: string
  snapshot_withdrawable_before: number
  snapshot_withdrawable_after: number
}

export interface WithdrawalEventItem {
  id: number
  request_id: number
  action: string
  operator_user_id?: number | null
  detail: string
  created_at: string
}

export interface WithdrawalDetailItem extends WithdrawalItem {
  events: WithdrawalEventItem[]
}

export interface DistributorProfile {
  user_id: number
  status: string
  display_name: string
  settlement_channel: string
  settlement_account_name: string
  settlement_account_no: string
  settlement_account_extra: string
  notes: string
}

export interface UserLookupItem {
  id: number
  email: string
  username: string
  role: string
  status: string
  is_distributor: boolean
}

export interface RouteMetaConfig {
  title: string
  requiresAuth: boolean
  role?: PortalRole
}

export interface NavItem {
  key: string
  label: string
  path: string
  role: PortalRole
}

export interface TableColumn<T> {
  key: string
  title: string
  width?: string
  align?: 'left' | 'center' | 'right'
  render?: (item: T) => string
}
