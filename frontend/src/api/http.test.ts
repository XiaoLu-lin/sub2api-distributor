import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { apiRequest } from './http'

describe('apiRequest', () => {
  const fetchMock = vi.fn()

  beforeEach(() => {
    vi.stubGlobal('fetch', fetchMock)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
    fetchMock.mockReset()
  })

  it('adds json and authorization headers for authenticated requests', async () => {
    fetchMock.mockResolvedValue({
      ok: true,
      json: async () => ({ success: true }),
    })

    await expect(apiRequest<{ success: boolean }>('/demo', { method: 'POST' }, 'token-123')).resolves.toEqual({
      success: true,
    })

    expect(fetchMock).toHaveBeenCalledTimes(1)
    const [url, options] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(url).toBe('http://127.0.0.1:8091/api/demo')

    const headers = options.headers as Headers
    expect(headers.get('Content-Type')).toBe('application/json')
    expect(headers.get('Authorization')).toBe('Bearer token-123')
  })

  it('throws the backend message when the request fails', async () => {
    fetchMock.mockResolvedValue({
      ok: false,
      status: 403,
      json: async () => ({ message: '无权访问当前页面' }),
    })

    await expect(apiRequest('/denied')).rejects.toThrow('无权访问当前页面')
  })

  it('falls back to an http status error when the response body is not json', async () => {
    fetchMock.mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error('bad json')
      },
    })

    await expect(apiRequest('/broken')).rejects.toThrow('请求失败，状态码：500')
  })
})
