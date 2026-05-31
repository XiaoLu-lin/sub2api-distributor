import { computed, reactive } from 'vue'
import { logout as logoutRequest, me } from '../api'
import type { MeResponse, PortalRole } from '../types'
import { clearStoredSession, readStoredSession, writeStoredSession } from './storage'

interface SessionState {
  token: string
  portalRole: PortalRole | ''
  email: string
  me: MeResponse | null
  hydrated: boolean
}

const stored = readStoredSession()

const state = reactive<SessionState>({
  token: stored.token,
  portalRole: stored.portalRole,
  email: stored.email,
  me: null,
  hydrated: false,
})

function persist(): void {
  writeStoredSession({
    token: state.token,
    portalRole: state.portalRole,
    email: state.email,
  })
}

export function useSession() {
  const isAuthenticated = computed(() => Boolean(state.token))

  async function hydrate(): Promise<void> {
    if (state.hydrated) {
      return
    }

    if (!state.token) {
      state.hydrated = true
      return
    }

    try {
      state.me = await me(state.token)
      state.portalRole = state.me.portal_role
    } catch {
      clear()
    } finally {
      state.hydrated = true
    }
  }

  function setSession(token: string, portalRole: PortalRole, email: string): void {
    state.token = token
    state.portalRole = portalRole
    state.email = email
    persist()
  }

  function setMeProfile(profile: MeResponse | null): void {
    state.me = profile
    if (profile) {
      state.portalRole = profile.portal_role
    }
  }

  function clear(): void {
    state.token = ''
    state.portalRole = ''
    state.email = ''
    state.me = null
    clearStoredSession()
  }

  async function logout(): Promise<void> {
    if (state.token) {
      try {
        await logoutRequest(state.token)
      } catch {
        // Best-effort logout for stateless JWT flow.
      }
    }

    clear()
  }

  return {
    state,
    isAuthenticated,
    hydrate,
    setSession,
    setMeProfile,
    clear,
    logout,
  }
}
