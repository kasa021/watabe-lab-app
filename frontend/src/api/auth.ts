import { apiClient } from './client'
import { User } from '../types'

export interface LoginResponse {
  token: string
  user: User
  expires_at: string
}

export const authApi = {
  login: async (username: string, password: string): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>('/api/v1/auth/login', {
      username,
      password,
    })
    return response.data
  },

  getMe: async (): Promise<User> => {
    const response = await apiClient.get<User>('/api/v1/auth/me')
    return response.data
  },
}
