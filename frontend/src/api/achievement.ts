import { apiClient } from './client'
import { User } from '../types'

export interface Achievement {
  id: number
  code: string
  name: string
  description: string
  icon_url?: string
  category: string
  points_reward: number
  is_active: boolean
}

export interface UserAchievement {
  id: number
  user_id: number
  achievement_id: number
  achieved_at: string
  achievement: Achievement
  user?: User
}

export const achievementApi = {
  getAchievements: async (): Promise<Achievement[]> => {
    const response = await apiClient.get<{ achievements: Achievement[] }>('/api/v1/achievements')
    return response.data.achievements
  },

  getMyAchievements: async (): Promise<UserAchievement[]> => {
    const response = await apiClient.get<{ user_achievements: UserAchievement[] }>('/api/v1/achievements/my')
    return response.data.user_achievements
  },
}
