import { apiClient } from './client'

export interface User {
  id: number
  username: string
  display_name: string
  role: string
  is_active: boolean
}

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
  
  getMe: async (): Promise<{ user: User }> => {
    const response = await apiClient.get<{ user: User }>('/api/v1/auth/me')
    return response.data
  },
}
