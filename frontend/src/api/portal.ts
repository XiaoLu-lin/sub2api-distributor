import type {
  DashboardSummary,
  DistributorProfile,
  InviteMeta,
  InviteeItem,
  RebateItem,
  WithdrawalItem,
} from '../types'
import { apiRequest } from './http'

export function getDashboard(token: string): Promise<DashboardSummary> {
  return apiRequest<DashboardSummary>('/portal/dashboard', {}, token)
}

export async function getInvitees(token: string): Promise<InviteeItem[]> {
  const resp = await apiRequest<{ items: InviteeItem[] }>('/portal/invitees', {}, token)
  return resp.items
}

export function getInviteMeta(token: string): Promise<InviteMeta> {
  return apiRequest<InviteMeta>('/portal/invite-meta', {}, token)
}

export async function getRebates(token: string): Promise<RebateItem[]> {
  const resp = await apiRequest<{ items: RebateItem[] }>('/portal/rebates', {}, token)
  return resp.items
}

export async function getWithdrawals(token: string): Promise<WithdrawalItem[]> {
  const resp = await apiRequest<{ items: WithdrawalItem[] }>('/portal/withdrawals', {}, token)
  return resp.items
}

export function createWithdrawal(token: string, amount: number, remark: string): Promise<WithdrawalItem> {
  return apiRequest<WithdrawalItem>(
    '/portal/withdrawals',
    {
      method: 'POST',
      body: JSON.stringify({ amount, remark }),
    },
    token,
  )
}

export function cancelWithdrawal(token: string, id: number): Promise<{ success: boolean }> {
  return apiRequest<{ success: boolean }>(
    `/portal/withdrawals/${id}/cancel`,
    {
      method: 'POST',
    },
    token,
  )
}

export function getProfile(token: string): Promise<DistributorProfile> {
  return apiRequest<DistributorProfile>('/portal/settlement-profile', {}, token)
}

export function updateProfile(token: string, payload: DistributorProfile): Promise<{ success: boolean }> {
  return apiRequest<{ success: boolean }>(
    '/portal/settlement-profile',
    {
      method: 'PUT',
      body: JSON.stringify(payload),
    },
    token,
  )
}
