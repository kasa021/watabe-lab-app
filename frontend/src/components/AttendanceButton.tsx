import React, { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
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
  const { t } = useTranslation()
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

  const handleCheckIn = () => {
    setIsLoading(true)

    if (!navigator.geolocation) {
      alert("お使いのブラウザは位置情報をサポートしていません")
      setIsLoading(false)
      return
    }

    navigator.geolocation.getCurrentPosition(
      async (position) => {
        try {
          await attendanceApi.checkIn({
            wifi_ssid: 'WatabeLabWiFi',
            check_in_method: 'web_manual',
            gps_latitude: position.coords.latitude,
            gps_longitude: position.coords.longitude,
          })
          
          setIsCheckedIn(true)
          onStatusChange?.(true)
          alert(t('attendance.checkin_success'))
        } catch (error: any) {
          if (error.response?.status === 409) {
            setIsCheckedIn(true)
            onStatusChange?.(true)
            alert(t('attendance.already_checked_in'))
          } else if (error.response?.status === 403) {
            // Restriction violation (IP or GPS)
            console.log('403 Forbidden caught', error.response.data)
            const message = error.response?.data?.error || "アクセスが拒否されました"
            alert(message)
          } else {
            console.error('Check-in failed:', error)
            const message = error.response?.data?.error || t('attendance.checkin_failed')
            alert(message)
          }
        } finally {
          setIsLoading(false)
        }
      },
      (error) => {
        console.error('Geolocation error:', error)
        alert("位置情報の取得に失敗しました。位置情報の利用を許可してください。")
        setIsLoading(false)
      }
    )
  }

  const handleCheckOut = async () => {
    setIsLoading(true)
    try {
      await attendanceApi.checkOut()
      // WebSocket経由で更新されるはずだが、念のためローカルも更新
      setIsCheckedIn(false)
      onStatusChange?.(false)
      alert(t('attendance.checkout_success'))
    } catch (error) {
      console.error('Check-out failed:', error)
      alert(t('attendance.checkout_failed'))
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
        t('attendance.processing')
      ) : isCheckedIn ? (
        <div className="flex flex-col items-center">
          <span>{t('attendance.exit')}</span>
          <span className="text-sm font-normal mt-2">{t('attendance.current_status_in')}</span>
        </div>
      ) : (
        <div className="flex flex-col items-center">
          <span>{t('attendance.enter')}</span>
          <span className="text-sm font-normal mt-2">{t('attendance.current_status_out')}</span>
        </div>
      )}
    </button>
  )
}
