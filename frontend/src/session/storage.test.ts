import { beforeEach, describe, expect, it } from 'vitest'
import { clearStoredSession, readStoredSession, writeStoredSession } from './storage'

function createLocalStorageMock(): Storage {
  const store = new Map<string, string>()

  return {
    get length() {
      return store.size
    },
    clear() {
      store.clear()
    },
    getItem(key: string) {
      return store.has(key) ? store.get(key)! : null
    },
    key(index: number) {
      return Array.from(store.keys())[index] ?? null
    },
    removeItem(key: string) {
      store.delete(key)
    },
    setItem(key: string, value: string) {
      store.set(key, value)
    },
  }
}

describe('session storage helpers', () => {
  beforeEach(() => {
    Object.defineProperty(globalThis, 'localStorage', {
      configurable: true,
      value: createLocalStorageMock(),
    })
  })

  it('returns empty defaults when storage has no session', () => {
    expect(readStoredSession()).toEqual({
      token: '',
      portalRole: '',
      email: '',
    })
  })

  it('writes and reads a persisted session', () => {
    writeStoredSession({
      token: 'token-1',
      portalRole: 'distributor',
      email: 'demo@example.com',
    })

    expect(readStoredSession()).toEqual({
      token: 'token-1',
      portalRole: 'distributor',
      email: 'demo@example.com',
    })
  })

  it('clears all stored session fields', () => {
    writeStoredSession({
      token: 'token-2',
      portalRole: 'operator',
      email: 'ops@example.com',
    })

    clearStoredSession()

    expect(readStoredSession()).toEqual({
      token: '',
      portalRole: '',
      email: '',
    })
  })
})
