import { create } from 'zustand'
import { CheckInLog } from '../types'
import { apiClient } from '../api/client'

interface OccupancyState {
  activeUsers: CheckInLog[]
  isConnected: boolean
  fetchActiveUsers: () => Promise<void>
  connect: () => void
  disconnect: () => void
}

export const useOccupancyStore = create<OccupancyState>((set, get) => {
  let socket: WebSocket | null = null

  return {
    activeUsers: [],
    isConnected: false,

    fetchActiveUsers: async () => {
      try {
        const response = await apiClient.get<{ active_users: CheckInLog[] }>('/api/v1/attendance/active')
        // backend returns { active_users: [...] } based on handler implementation
        set({ activeUsers: response.data.active_users || [] })
      } catch (error) {
        console.error('Failed to fetch active users', error)
      }
    },

    connect: () => {
      if (socket) return


      const apiBase = import.meta.env.VITE_API_BASE_URL || ''
      const wsUrl = window.location.origin.replace(/^http/, 'ws') + apiBase + '/api/v1/ws'

      console.log('Connecting to WebSocket:', wsUrl)
      socket = new WebSocket(wsUrl)

      socket.onopen = () => {
        console.log('WebSocket connected')
        set({ isConnected: true })
        get().fetchActiveUsers()
      }

      socket.onclose = () => {
        console.log('WebSocket disconnected')
        set({ isConnected: false })
        socket = null
        // Simple reconnect logic
        setTimeout(() => {
            if (!get().isConnected) {
                get().connect()
            }
        }, 5000)
      }

      socket.onerror = (error) => {
          console.error('WebSocket error:', error)
      }

      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          const { type, payload } = data

          set((state) => {
            if (type === 'check_in') {
                // Remove existing if any (to update)
                const others = state.activeUsers.filter(u => u.user_id !== payload.user_id)
                return { activeUsers: [...others, payload] }
            } else if (type === 'check_out') {
                return { activeUsers: state.activeUsers.filter(u => u.user_id !== payload.user_id) }
            }
            return state
          })
        } catch (e) {
          console.error('Failed to parse WS message', e)
        }
      }
    },

    disconnect: () => {
      if (socket) {
        socket.close()
        socket = null
        set({ isConnected: false })
      }
    }
  }
})
