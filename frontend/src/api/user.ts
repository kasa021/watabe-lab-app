import { apiClient } from './client'
import { User } from '../types'

export interface UpdateProfileRequest {
  display_name: string
  is_presence_public: boolean
}

export interface HeatmapData {
  date: string
  count: number
  duration: number
}

export const userApi = {
  updateProfile: async (data: UpdateProfileRequest): Promise<User> => {
    const response = await apiClient.put<User>('/api/v1/users/me', data)
    return response.data
  },

  getHeatmap: async (userId: string | number = 'me'): Promise<HeatmapData[]> => {
    const response = await apiClient.get<{ heatmap: HeatmapData[] }>(`/api/v1/users/${userId}/heatmap`)
    return response.data.heatmap
  },
}
