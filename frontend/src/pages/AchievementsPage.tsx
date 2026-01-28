import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { achievementApi, Achievement, UserAchievement } from '../api/achievement'
import { 
  Medal, 
  Lock, 
  Clock, 
  Sunrise, 
  Moon, 
  Trophy, 
  Flame,
  CalendarCheck,
  LucideIcon
} from 'lucide-react'

const AchievementsPage = () => {
  const { t } = useTranslation()
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

  const getAchievementIcon = (achievement: Achievement): LucideIcon => {
    const { code, category } = achievement
    
    // Specific implementation based on achievement code
    if (code.includes('early_bird')) return Sunrise
    if (code.includes('night_owl')) return Moon
    if (code.includes('streak')) return Flame
    
    // Category based implementation
    switch (category) {
      case 'time':
        return Clock
      case 'attendance':
        return CalendarCheck
      case 'special':
        return Trophy
      default:
        return Medal
    }
  }

  const totalPoints = myAchievements.reduce((sum, ua) => sum + ua.achievement.points_reward, 0)

  return (
    <div className="space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">{t('achievements.title')}</h1>
        <p className="text-gray-500">
          {t('achievements.total_points')}: <span className="text-2xl font-bold text-primary-600">{totalPoints} pt</span>
        </p>
      </div>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
          <p className="mt-4 text-gray-500">{t('achievements.loading')}</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {achievements.map((ach) => {
            const achieved = isAchieved(ach.id)
            const IconComponent = getAchievementIcon(ach)
            
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
                    {t('achievements.get')}
                  </div>
                )}
                
                <div className="flex items-start space-x-4">
                  <div className={`
                    w-12 h-12 flex-shrink-0 rounded-full flex items-center justify-center
                    ${achieved ? 'bg-gradient-to-br from-yellow-100 to-orange-100 text-yellow-600' : 'bg-gray-200 text-gray-400'}
                  `}>
                    {achieved ? (
                      <IconComponent size={28} strokeWidth={1.5} />
                    ) : (
                      <Lock size={24} />
                    )}
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
                                {t('achievements.unlocked_at', { date: getAchievedDate(ach.id) })}
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
