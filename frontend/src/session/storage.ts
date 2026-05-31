import type { PortalRole } from '../types'

const TOKEN_KEY = 'distributor_token'
const ROLE_KEY = 'distributor_role'
const EMAIL_KEY = 'distributor_email'

export interface StoredSession {
  token: string
  portalRole: PortalRole | ''
  email: string
}

export function readStoredSession(): StoredSession {
  return {
    token: localStorage.getItem(TOKEN_KEY) || '',
    portalRole: (localStorage.getItem(ROLE_KEY) as PortalRole | null) || '',
    email: localStorage.getItem(EMAIL_KEY) || '',
  }
}

export function writeStoredSession(session: StoredSession): void {
  localStorage.setItem(TOKEN_KEY, session.token)
  localStorage.setItem(ROLE_KEY, session.portalRole)
  localStorage.setItem(EMAIL_KEY, session.email)
}

export function clearStoredSession(): void {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(ROLE_KEY)
  localStorage.removeItem(EMAIL_KEY)
}
