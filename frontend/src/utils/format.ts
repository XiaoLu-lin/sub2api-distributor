export function formatCurrency(amount: number | null | undefined): string {
  return `¥${(amount ?? 0).toFixed(2)}`
}

export function buildInviteLink(affCode: string | null | undefined, mainAppBaseUrl: string | null | undefined): string {
  const normalizedAffCode = affCode?.trim()
  const normalizedBaseUrl = mainAppBaseUrl?.trim()

  if (!normalizedAffCode || !normalizedBaseUrl) {
    return ''
  }

  return `${normalizedBaseUrl.replace(/\/+$/, '')}/register?aff=${encodeURIComponent(normalizedAffCode)}`
}

export function formatDateTime(value: string | null | undefined): string {
  if (!value) {
    return '-'
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function statusText(status: string): string {
  return (
    {
      paying: '打款中',
      paid: '已打款',
      cancelled: '已取消',
      active: '启用',
      disabled: '停用',
    }[status] ?? status
  )
}

export function eventActionText(action: string): string {
  return (
    {
      create: '创建申请',
      mark_paid: '标记已打款',
      cancel: '取消申请',
    }[action] ?? action
  )
}

function parseEventDetail(detail: string): Record<string, unknown> | null {
  if (!detail) {
    return null
  }

  try {
    const parsed = JSON.parse(detail) as unknown
    if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
      return parsed as Record<string, unknown>
    }
  } catch {
    return null
  }

  return null
}

function normalizeTextValue(value: unknown): string {
  if (typeof value === 'number') {
    return String(value)
  }

  if (typeof value === 'string' && value.trim()) {
    return value.trim()
  }

  return ''
}

export function formatEventDetail(action: string, detail: string): string {
  const parsed = parseEventDetail(detail)
  if (!parsed) {
    return detail?.trim() || '暂无操作说明'
  }

  if (action === 'create') {
    const amountValue =
      typeof parsed.amount === 'number'
        ? parsed.amount
        : typeof parsed.amount === 'string'
          ? Number(parsed.amount)
          : Number.NaN
    const amount = Number.isFinite(amountValue) ? formatCurrency(amountValue) : ''
    const remark = normalizeTextValue(parsed.remark)

    if (amount && remark) {
      return `申请金额：${amount}，备注：${remark}`
    }
    if (amount) {
      return `申请金额：${amount}`
    }
    if (parsed.seed === true) {
      return '系统已写入演示提现申请'
    }
    if (remark) {
      return `备注：${remark}`
    }
    return '已创建提现申请'
  }

  if (action === 'mark_paid') {
    const channelMap: Record<string, string> = {
      bank: '银行卡',
      alipay: '支付宝',
      wechat: '微信',
      usdt: 'USDT',
      manual: '人工协商',
    }

    const parts = [
      normalizeTextValue(parsed.paid_channel) ? `打款渠道：${channelMap[normalizeTextValue(parsed.paid_channel)] ?? normalizeTextValue(parsed.paid_channel)}` : '',
      normalizeTextValue(parsed.paid_reference_no) ? `打款流水号：${normalizeTextValue(parsed.paid_reference_no)}` : '',
      normalizeTextValue(parsed.paid_remark) ? `打款备注：${normalizeTextValue(parsed.paid_remark)}` : '',
    ].filter(Boolean)

    return parts.length ? parts.join('，') : '运营已标记线下打款完成'
  }

  if (action === 'cancel') {
    const reason = normalizeTextValue(parsed.reason)
    const remark = normalizeTextValue(parsed.paid_remark)
    if (parsed.seed === true) {
      return '系统已取消该演示申请'
    }

    if (reason) {
      return `取消原因：${reason}`
    }
    if (remark) {
      return `取消备注：${remark}`
    }
    return '该提现申请已被取消'
  }

  return '已记录操作变更'
}
