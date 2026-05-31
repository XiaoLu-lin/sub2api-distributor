import { describe, expect, it } from 'vitest'
import { distributorNavItems, operatorNavItems } from './navigation'

describe('navigation config', () => {
  it('exposes the expected distributor menu order', () => {
    expect(distributorNavItems.map((item) => item.path)).toEqual([
      '/portal/dashboard',
      '/portal/invitees',
      '/portal/rebates',
      '/portal/withdrawals',
      '/portal/settlement',
    ])
  })

  it('exposes the expected operator menu order', () => {
    expect(operatorNavItems.map((item) => item.path)).toEqual([
      '/ops/distributors',
      '/ops/withdrawals',
    ])
  })

  it('keeps all nav items aligned to their owning role namespace', () => {
    for (const item of distributorNavItems) {
      expect(item.role).toBe('distributor')
      expect(item.path.startsWith('/portal/')).toBe(true)
    }

    for (const item of operatorNavItems) {
      expect(item.role).toBe('operator')
      expect(item.path.startsWith('/ops/')).toBe(true)
    }
  })
})
