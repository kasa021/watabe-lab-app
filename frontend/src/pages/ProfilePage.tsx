import { useEffect, useState } from 'react'
import CalendarHeatmap from 'react-calendar-heatmap'
import 'react-calendar-heatmap/dist/styles.css'
import { Tooltip } from 'react-tooltip'
import { userApi, HeatmapData, UpdateProfileRequest } from '../api/user'
import { authApi } from '../api/auth'
import { apiClient } from '../api/client' // For user data re-fetch if needed
import { User } from '../types'

const ProfilePage = () => {
  const [user, setUser] = useState<User | null>(null)
  const [heatmapData, setHeatmapData] = useState<HeatmapData[]>([])
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  
  // Form state
  const [displayName, setDisplayName] = useState('')
  const [isPresencePublic, setIsPresencePublic] = useState(true)
  const [message, setMessage] = useState<{ text: string; type: 'success' | 'error' } | null>(null)

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Fetch current user data from /auth/me (fresh data)
        const user = await authApi.getMe()
        setUser(user)
        setDisplayName(user.display_name)
        setIsPresencePublic(user.is_presence_public)
        
        // Fetch heatmap data
        const data = await userApi.getHeatmap('me')
        setHeatmapData(data)
      } catch (error) {
        console.error('Failed to fetch profile data', error)
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [])

  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    setMessage(null)
    
    try {
      const req: UpdateProfileRequest = {
        display_name: displayName,
        is_presence_public: isPresencePublic,
      }
      const updatedUser = await userApi.updateProfile(req)
      setUser(updatedUser)
      // Update local storage if needed, though we rely on API for truth
      localStorage.setItem('user', JSON.stringify(updatedUser))
      setMessage({ text: 'プロフィールを更新しました', type: 'success' })
    } catch (error) {
       console.error(error)
       setMessage({ text: '更新に失敗しました', type: 'error' })
    } finally {
      setSaving(false)
    }
  }

  // Heatmap helper
  const getTooltipDataAttrs = (value: any) => {
    if (!value || !value.date) {
      return { 'data-tooltip-content': 'No data' }
    }
    return {
      'data-tooltip-id': 'heatmap-tooltip',
      'data-tooltip-content': `${value.date}: ${value.duration}分 (${value.count}回)`,
    }
  }

  const classForValue = (value: any) => {
    if (!value || value.count === 0) {
      return 'color-empty'
    }
    // Simple scaling based on duration (minutes)
    // 0-60min: scale-1
    // 60-180min: scale-2
    // 180-300min: scale-3
    // 300+: scale-4
    if (value.duration < 60) return 'color-scale-1'
    if (value.duration < 180) return 'color-scale-2'
    if (value.duration < 300) return 'color-scale-3'
    return 'color-scale-4'
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">プロフィール設定</h1>
        <p className="text-gray-500 mt-2">あなたの活動履歴と公開設定を管理します。</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Profile Settings Card */}
        <div className="lg:col-span-1">
            <div className="bg-white shadow rounded-lg p-6">
                <h2 className="text-xl font-bold text-gray-900 mb-4">アカウント設定</h2>
                
                <form onSubmit={handleUpdateProfile} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700">ユーザー名</label>
                        <input 
                            type="text" 
                            value={user?.username || ''} 
                            disabled 
                            className="mt-1 block w-full px-3 py-2 bg-gray-100 border border-gray-300 rounded-md shadow-sm sm:text-sm text-gray-500"
                        />
                    </div>
                    
                    <div>
                        <label className="block text-sm font-medium text-gray-700">表示名</label>
                        <input 
                            type="text" 
                            value={displayName} 
                            onChange={(e) => setDisplayName(e.target.value)}
                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                        />
                         <p className="text-xs text-gray-500 mt-1">ランキングや一覧に表示される名前です。</p>
                    </div>

                    <div className="flex items-start">
                        <div className="flex items-center h-5">
                            <input
                                id="is_public"
                                name="is_public"
                                type="checkbox"
                                checked={isPresencePublic}
                                onChange={(e) => setIsPresencePublic(e.target.checked)}
                                className="focus:ring-primary-500 h-4 w-4 text-primary-600 border-gray-300 rounded"
                            />
                        </div>
                        <div className="ml-3 text-sm">
                            <label htmlFor="is_public" className="font-medium text-gray-700">在室状況を公開する</label>
                            <p className="text-gray-500">オフにすると「現在のアクティブユーザー」に表示されなくなります。</p>
                        </div>
                    </div>

                    {message && (
                        <div className={`p-3 rounded-md text-sm ${message.type === 'success' ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'}`}>
                            {message.text}
                        </div>
                    )}

                    <button
                        type="submit"
                        disabled={saving}
                        className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
                    >
                        {saving ? '保存中...' : '設定を保存'}
                    </button>
                </form>
            </div>
        </div>

        {/* Heatmap & Stats */}
        <div className="lg:col-span-2 space-y-8">
            {/* Heatmap */}
            <div className="bg-white shadow rounded-lg p-6">
                <h2 className="text-xl font-bold text-gray-900 mb-6">Attendance Heatmap</h2>
                <div className="w-full overflow-x-auto">
                    <div className="min-w-[600px]">
                        <CalendarHeatmap
                            startDate={new Date(new Date().setFullYear(new Date().getFullYear() - 1))}
                            endDate={new Date()}
                            values={heatmapData}
                            classForValue={classForValue}
                            tooltipDataAttrs={getTooltipDataAttrs}
                            showWeekdayLabels={true}
                        />
                        <Tooltip id="heatmap-tooltip" />
                    </div>
                </div>
                <div className="mt-4 flex items-center justify-end text-sm text-gray-500 space-x-2">
                    <span>Less</span>
                    <div className="flex space-x-1">
                        <div className="w-3 h-3 bg-gray-100 rounded-sm"></div>
                        <div className="w-3 h-3 bg-green-200 rounded-sm"></div>
                        <div className="w-3 h-3 bg-green-400 rounded-sm"></div>
                        <div className="w-3 h-3 bg-green-600 rounded-sm"></div>
                        <div className="w-3 h-3 bg-green-800 rounded-sm"></div>
                    </div>
                    <span>More</span>
                </div>
            </div>
            
            {/* Stats (Placeholder or simple calc) */}
             <div className="grid grid-cols-2 gap-4">
                <div className="bg-white shadow rounded-lg p-6 text-center">
                    <p className="text-sm font-medium text-gray-500">過去1年間の出席日数</p>
                    <p className="text-3xl font-bold text-primary-600 mt-2">{heatmapData.length} 日</p>
                </div>
                <div className="bg-white shadow rounded-lg p-6 text-center">
                    <p className="text-sm font-medium text-gray-500">過去1年間の滞在時間</p>
                    <p className="text-3xl font-bold text-primary-600 mt-2">
                        {Math.round(heatmapData.reduce((acc, curr) => acc + curr.duration, 0) / 60)} 時間
                    </p>
                </div>
            </div>
        </div>
      </div>
    </div>
  )
}

export default ProfilePage
