import { apiClient } from './client'

export interface CheckInRequest {
  wifi_ssid: string
  check_in_method: string
  gps_latitude?: number
  gps_longitude?: number
}

export const attendanceApi = {
  checkIn: async (data: CheckInRequest): Promise<void> => {
    await apiClient.post('/api/v1/attendance/checkin', data)
  },

  checkOut: async (): Promise<void> => {
    await apiClient.post('/api/v1/attendance/checkout')
  },
}
