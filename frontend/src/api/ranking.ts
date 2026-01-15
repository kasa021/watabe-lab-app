import { apiClient } from './client'


export interface UserRanking {
  user_id: number
  display_name: string
  username: string
  total_duration: number
}

export const rankingApi = {
  getRankings: async (type: 'weekly' | 'monthly' | 'total' = 'weekly'): Promise<UserRanking[]> => {
    const response = await apiClient.get<{ rankings: UserRanking[] }>('/api/v1/rankings', {
      params: { type },
    })
    return response.data.rankings
  },
}
