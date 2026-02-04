import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { apiClient } from '../api/client'
import { AttendanceButton } from '../components/AttendanceButton'
import { ActiveUsersList } from '../components/ActiveUsersList'
import { useOccupancyStore } from '../stores/useOccupancyStore'
import { User } from '../types'

interface HealthResponse {
  status: string
  message: string
}

const HomePage = () => {
  const { t } = useTranslation()
  const [health, setHealth] = useState<HealthResponse | null>(null)
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [user, setUser] = useState<User | null>(null)

  const { connect, disconnect } = useOccupancyStore()

  useEffect(() => {
    // Check login status
    const token = localStorage.getItem('token')
    const userStr = localStorage.getItem('user')
    if (token && userStr) {
      setIsLoggedIn(true)
      setUser(JSON.parse(userStr))
      // Connect to WebSocket when logged in
      connect()
    }

    const checkHealth = async () => {
      try {
        const response = await apiClient.get<HealthResponse>('/api/v1/health')
        setHealth(response.data)
      } catch (error) {
        console.error('Health check failed:', error)
      }
    }
    checkHealth()

    return () => {
        disconnect()
    }
  }, [])

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto text-center">
        <h1 className="text-4xl font-extrabold text-gray-900 sm:text-5xl md:text-6xl mb-8">
          {t('home.title')}
        </h1>

        {isLoggedIn && user ? (
          <div className="space-y-8">
            <div className="bg-white rounded-xl shadow-2xl p-8 max-w-2xl mx-auto transform transition-all">
                <h2 className="text-2xl font-bold text-gray-800 mb-8">
                {t('home.welcome', { name: user.display_name })}
                </h2>
                <div className="flex justify-center mb-8">
                <AttendanceButton userId={user.id} />
                </div>
                <p className="text-gray-500 text-sm">
                {t('home.instruction')}
                </p>
            </div>
            
            <ActiveUsersList />
          </div>
        ) : (
          <div className="bg-white rounded-xl shadow-xl p-8 max-w-lg mx-auto">
            <p className="text-lg text-gray-600 mb-6">
              {t('home.guest_message')}
            </p>
            {/* System Status - Only show when not logged in or as footer */}
            <div className="mb-6 p-4 bg-gray-50 rounded-lg">
                <p className="text-sm text-gray-500">
                  {t('home.system_status')}: {health ? <span className="text-green-600 font-bold">{t('home.online')}</span> : <span className="text-red-500">{t('home.connecting')}</span>}
                </p>
            </div>
            
            <a
              href="/attendance/login"
              className="inline-block w-full bg-primary-600 hover:bg-primary-700 text-white font-bold py-4 rounded-lg shadow-lg hover:shadow-xl transition-all transform hover:-translate-y-1"
            >
              {t('home.login_start')}
            </a>
          </div>
        )}
      </div>
    </div>
  )
}

export default HomePage
