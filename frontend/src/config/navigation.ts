import type { NavItem } from '../types'

export const distributorNavItems: NavItem[] = [
  { key: 'dashboard', label: '概览', path: '/portal/dashboard', role: 'distributor' },
  { key: 'invitees', label: '邀请用户', path: '/portal/invitees', role: 'distributor' },
  { key: 'rebates', label: '返利明细', path: '/portal/rebates', role: 'distributor' },
  { key: 'withdrawals', label: '提现申请', path: '/portal/withdrawals', role: 'distributor' },
  { key: 'settlement', label: '收款信息', path: '/portal/settlement', role: 'distributor' },
]

export const operatorNavItems: NavItem[] = [
  { key: 'distributors', label: '分销商管理', path: '/ops/distributors', role: 'operator' },
  { key: 'withdrawals', label: '提现管理', path: '/ops/withdrawals', role: 'operator' },
]
