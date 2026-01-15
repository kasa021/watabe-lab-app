import { useEffect, useState } from 'react'
import { achievementApi, Achievement, UserAchievement } from '../api/achievement'

const AchievementsPage = () => {
  const [achievements, setAchievements] = useState<Achievement[]>([])
  const [myAchievements, setMyAchievements] = useState<UserAchievement[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true)
      try {
        const [allData, myData] = await Promise.all([
          achievementApi.getAchievements(),
          achievementApi.getMyAchievements(),
        ])
        setAchievements(allData)
        setMyAchievements(myData)
      } catch (error) {
        console.error('Failed to fetch achievements:', error)
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [])

  const isAchieved = (achievementId: number) => {
    return myAchievements.some((ua) => ua.achievement_id === achievementId)
  }

  const getAchievedDate = (achievementId: number) => {
    const ua = myAchievements.find((ua) => ua.achievement_id === achievementId)
    if (ua) {
      return new Date(ua.achieved_at).toLocaleDateString()
    }
    return null
  }

  const totalPoints = myAchievements.reduce((sum, ua) => sum + ua.achievement.points_reward, 0)

  return (
    <div className="space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">🏆 実績リスト</h1>
        <p className="text-gray-500">
          合計獲得ポイント: <span className="text-2xl font-bold text-primary-600">{totalPoints} pt</span>
        </p>
      </div>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
          <p className="mt-4 text-gray-500">読み込み中...</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {achievements.map((ach) => {
            const achieved = isAchieved(ach.id)
            return (
              <div
                key={ach.id}
                className={`
                  relative p-6 rounded-xl border transition-all duration-300
                  ${
                    achieved
                      ? 'bg-white border-yellow-400 shadow-md transform hover:-translate-y-1'
                      : 'bg-gray-50 border-gray-200 opacity-75 grayscale'
                  }
                `}
              >
                {achieved && (
                  <div className="absolute top-0 right-0 -mt-2 -mr-2 bg-yellow-400 text-white text-xs font-bold px-3 py-1 rounded-full shadow-sm">
                    GET!
                  </div>
                )}
                
                <div className="flex items-start space-x-4">
                  <div className={`
                    w-12 h-12 flex-shrink-0 rounded-full flex items-center justify-center text-2xl
                    ${achieved ? 'bg-yellow-100' : 'bg-gray-200'}
                  `}>
                    {achieved ? '🥇' : '🔒'}
                  </div>
                  <div>
                    <h3 className={`font-bold text-lg ${achieved ? 'text-gray-900' : 'text-gray-500'}`}>
                      {ach.name}
                    </h3>
                    <p className="text-sm text-gray-500 mt-1 mb-2 min-h-[40px]">
                      {ach.description}
                    </p>
                    <div className="flex items-center justify-between text-xs">
                        <span className="bg-primary-50 text-primary-700 px-2 py-1 rounded">
                            {ach.points_reward} pt
                        </span>
                        {achieved && (
                            <span className="text-green-600 font-medium">
                                {getAchievedDate(ach.id)} 解除
                            </span>
                        )}
                    </div>
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}

export default AchievementsPage
