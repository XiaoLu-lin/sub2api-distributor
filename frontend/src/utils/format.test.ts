import { describe, expect, it } from 'vitest'
import {
  buildInviteLink,
  eventActionText,
  formatCurrency,
  formatDateTime,
  formatEventDetail,
  statusText,
} from './format'

describe('format helpers', () => {
  it('formats currency with a fallback zero value', () => {
    expect(formatCurrency(undefined)).toBe('¥0.00')
    expect(formatCurrency(12.345)).toBe('¥12.35')
  })

  it('formats valid datetimes for zh-CN locale output', () => {
    expect(formatDateTime('2026-05-31T09:30:00Z')).toContain('2026')
  })

  it('returns a dash for empty datetime values', () => {
    expect(formatDateTime('')).toBe('-')
    expect(formatDateTime(undefined)).toBe('-')
  })

  it('falls back to the original string when the datetime is invalid', () => {
    expect(formatDateTime('not-a-date')).toBe('not-a-date')
  })

  it('maps known status and event labels to Chinese text', () => {
    expect(statusText('paying')).toBe('打款中')
    expect(statusText('disabled')).toBe('停用')
    expect(eventActionText('mark_paid')).toBe('标记已打款')
  })

  it('returns original text for unknown status and event labels', () => {
    expect(statusText('archived')).toBe('archived')
    expect(eventActionText('custom')).toBe('custom')
  })

  it('builds an invite link from the configured main app base url', () => {
    expect(buildInviteLink('AFF2026', 'https://main.example.com')).toBe(
      'https://main.example.com/register?aff=AFF2026',
    )
    expect(buildInviteLink('AFF2026', 'https://main.example.com/')).toBe(
      'https://main.example.com/register?aff=AFF2026',
    )
  })

  it('does not generate a fake invite link when configuration is missing', () => {
    expect(buildInviteLink('AFF2026', '')).toBe('')
    expect(buildInviteLink('AFF2026', '   ')).toBe('')
    expect(buildInviteLink('', 'https://main.example.com')).toBe('')
  })

  it('formats event detail into readable Chinese text', () => {
    expect(formatEventDetail('create', '{"amount":1.23,"remark":"测试申请"}')).toBe(
      '申请金额：¥1.23，备注：测试申请',
    )
    expect(formatEventDetail('mark_paid', '{"paid_channel":"alipay","paid_reference_no":"ALI-1"}')).toBe(
      '打款渠道：支付宝，打款流水号：ALI-1',
    )
    expect(formatEventDetail('create', '{"seed":true}')).toBe('系统已写入演示提现申请')
    expect(formatEventDetail('custom', '{"seed":true}')).toBe('已记录操作变更')
    expect(formatEventDetail('custom', '')).toBe('暂无操作说明')
  })
})
