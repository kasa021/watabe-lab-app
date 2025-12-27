// ユーザー型
export interface User {
  id: number
  username: string
  display_name: string
  email?: string
  role: 'student' | 'teacher' | 'admin'
  is_presence_public: boolean
  created_at: string
  updated_at: string
  last_login_at?: string
  is_active: boolean
}

// チェックインログ型
export interface CheckInLog {
  id: number
  user_id: number
  check_in_at: string
  check_out_at?: string
  duration_minutes?: number
  check_in_method?: string
  wifi_ssid?: string
  gps_latitude?: number
  gps_longitude?: number
  created_at: string
  updated_at: string
  user?: User
}

// 日次出席記録型
export interface DailyAttendance {
  id: number
  user_id: number
  attendance_date: string
  total_duration_minutes: number
  check_in_count: number
  first_check_in_at?: string
  last_check_out_at?: string
  points: number
  is_holiday: boolean
  created_at: string
  updated_at: string
  user?: User
}

// 称号型
export interface Achievement {
  id: number
  code: string
  name: string
  description?: string
  icon_url?: string
  category?: string
  condition_type: string
  condition_value?: Record<string, unknown>
  points_reward: number
  is_active: boolean
  display_order: number
  created_at: string
  updated_at: string
}

// ユーザー称号型
export interface UserAchievement {
  id: number
  user_id: number
  achievement_id: number
  achieved_at: string
  user?: User
  achievement?: Achievement
}

// API レスポンス型
export interface ApiResponse<T> {
  data: T
  message?: string
}

export interface ApiError {
  error: {
    code: string
    message: string
  }
}

