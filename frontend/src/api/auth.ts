import type { LoginResponse, MeResponse } from '../types'
import { apiRequest } from './http'

export function login(email: string, password: string): Promise<LoginResponse> {
  return apiRequest<LoginResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export function me(token: string): Promise<MeResponse> {
  return apiRequest<MeResponse>('/me', {}, token)
}

export function logout(token: string): Promise<{ success: boolean }> {
  return apiRequest<{ success: boolean }>(
    '/auth/logout',
    {
      method: 'POST',
    },
    token,
  )
}
