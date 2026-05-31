import type { DistributorProfile, UserLookupItem, WithdrawalDetailItem, WithdrawalItem } from '../types'
import { apiRequest } from './http'

export async function getOpsDistributors(token: string): Promise<DistributorProfile[]> {
  const resp = await apiRequest<{ items: DistributorProfile[] }>('/ops/distributors', {}, token)
  return resp.items
}

export async function lookupOpsUsers(token: string, keyword: string): Promise<UserLookupItem[]> {
  const query = new URLSearchParams({ q: keyword })
  const resp = await apiRequest<{ items: UserLookupItem[] }>(`/ops/users/lookup?${query.toString()}`, {}, token)
  return resp.items
}

export function updateOpsDistributorProfile(
  token: string,
  userId: number,
  profile: DistributorProfile,
): Promise<{ success: boolean }> {
  return apiRequest<{ success: boolean }>(
    `/ops/distributors/${userId}/profile`,
    {
      method: 'PUT',
      body: JSON.stringify(profile),
    },
    token,
  )
}

export async function getOpsWithdrawals(token: string): Promise<WithdrawalItem[]> {
  const resp = await apiRequest<{ items: WithdrawalItem[] }>('/ops/withdrawals', {}, token)
  return resp.items
}

export function getOpsWithdrawalDetail(token: string, id: number): Promise<WithdrawalDetailItem> {
  return apiRequest<WithdrawalDetailItem>(`/ops/withdrawals/${id}`, {}, token)
}

export function markPaid(
  token: string,
  id: number,
  paidChannel: string,
  paidReferenceNo: string,
  paidRemark: string,
): Promise<{ success: boolean }> {
  return apiRequest<{ success: boolean }>(
    `/ops/withdrawals/${id}/mark-paid`,
    {
      method: 'POST',
      body: JSON.stringify({
        paid_channel: paidChannel,
        paid_reference_no: paidReferenceNo,
        paid_remark: paidRemark,
      }),
    },
    token,
  )
}

export function cancelOpsWithdrawal(token: string, id: number): Promise<{ success: boolean }> {
  return apiRequest<{ success: boolean }>(
    `/ops/withdrawals/${id}/cancel`,
    {
      method: 'POST',
    },
    token,
  )
}
