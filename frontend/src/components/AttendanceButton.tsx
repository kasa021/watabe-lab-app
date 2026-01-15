import React, { useState, useEffect } from 'react'
import { attendanceApi } from '../api/attendance'
import { useOccupancyStore } from '../stores/useOccupancyStore'

interface AttendanceButtonProps {
  userId?: number
  onStatusChange?: (isCheckedIn: boolean) => void
}

export const AttendanceButton: React.FC<AttendanceButtonProps> = ({
  userId,
  onStatusChange,
}) => {
  const [isCheckedIn, setIsCheckedIn] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const { activeUsers } = useOccupancyStore()

  // activeUsersの変更を監視して、自分の状態を同期する
  useEffect(() => {
    if (userId) {
      const isUserActive = activeUsers.some(u => u.user_id === userId)
      setIsCheckedIn(isUserActive)
    }
  }, [activeUsers, userId])

  const handleCheckIn = async () => {
    setIsLoading(true)
    try {
      await attendanceApi.checkIn({
        wifi_ssid: 'WatabeLabWiFi', // TODO: 実際にSSIDを取得するのはブラウザでは難しいので固定か、PWA/Native化が必要
        check_in_method: 'web_manual',
      })
      // WebSocket経由で更新されるはずだが、念のためローカルも更新
      setIsCheckedIn(true)
      onStatusChange?.(true)
      alert('チェックインしました！')
    } catch (error: any) {
      if (error.response?.status === 409) {
        setIsCheckedIn(true)
        onStatusChange?.(true)
        alert('既にチェックイン済みです。ステータスを更新しました。')
      } else {
        console.error('Check-in failed:', error)
        alert('チェックインに失敗しました')
      }
    } finally {
      setIsLoading(false)
    }
  }

  const handleCheckOut = async () => {
    setIsLoading(true)
    try {
      await attendanceApi.checkOut()
      // WebSocket経由で更新されるはずだが、念のためローカルも更新
      setIsCheckedIn(false)
      onStatusChange?.(false)
      alert('チェックアウトしました！お疲れ様でした。')
    } catch (error) {
      console.error('Check-out failed:', error)
      alert('チェックアウトに失敗しました')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <button
      onClick={isCheckedIn ? handleCheckOut : handleCheckIn}
      disabled={isLoading}
      className={`
        w-64 h-64 rounded-full text-2xl font-bold text-white shadow-lg transition-all transform hover:scale-105
        ${
          isLoading
            ? 'bg-gray-400 cursor-not-allowed'
            : isCheckedIn
            ? 'bg-red-500 hover:bg-red-600 shadow-red-500/50'
            : 'bg-green-500 hover:bg-green-600 shadow-green-500/50'
        }
      `}
    >
      {isLoading ? (
        '処理中...'
      ) : isCheckedIn ? (
        <div className="flex flex-col items-center">
          <span>退室する</span>
          <span className="text-sm font-normal mt-2">現在: 在室中</span>
        </div>
      ) : (
        <div className="flex flex-col items-center">
          <span>入室する</span>
          <span className="text-sm font-normal mt-2">現在: 退室中</span>
        </div>
      )}
    </button>
  )
}
