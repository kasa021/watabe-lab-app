import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { rankingApi, UserRanking } from '../api/ranking'

const RankingPage = () => {
  const { t } = useTranslation()
  const [rankings, setRankings] = useState<UserRanking[]>([])
  const [loading, setLoading] = useState(false)
  const [period, setPeriod] = useState<'weekly' | 'monthly' | 'total'>('weekly')

  useEffect(() => {
    const fetchRankings = async () => {
      setLoading(true)
      try {
        const data = await rankingApi.getRankings(period)
        setRankings(data)
      } catch (error) {
        console.error('Failed to fetch rankings:', error)
      } finally {
        setLoading(false)
      }
    }
    fetchRankings()
  }, [period])

  const formatDuration = (minutes: number) => {
    const hours = Math.floor(minutes / 60)
    const mins = minutes % 60
    return t('ranking.duration', { hours, minutes: mins })
  }

  const getMedalColor = (index: number) => {
    switch (index) {
      case 0:
        return 'bg-yellow-100 text-yellow-800 border-yellow-200' // Gold
      case 1:
        return 'bg-gray-100 text-gray-800 border-gray-200' // Silver
      case 2:
        return 'bg-orange-100 text-orange-800 border-orange-200' // Bronze
      default:
        return 'bg-white text-gray-800 border-gray-100'
    }
  }

  return (
    <div className="space-y-6">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">{t('ranking.title')}</h1>
        <p className="text-gray-500">{t('ranking.description')}</p>
      </div>

      {/* Period Selector */}
      <div className="flex justify-center space-x-4 mb-8">
        {(['weekly', 'monthly', 'total'] as const).map((p) => (
          <button
            key={p}
            onClick={() => setPeriod(p)}
            className={`
              px-6 py-2 rounded-full font-medium transition-all
              ${
                period === p
                  ? 'bg-primary-600 text-white shadow-md transform scale-105'
                  : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
              }
            `}
          >
            {p === 'weekly' && t('ranking.weekly')}
            {p === 'monthly' && t('ranking.monthly')}
            {p === 'total' && t('ranking.total')}
          </button>
        ))}
      </div>

      {/* Ranking List */}
      <div className="max-w-3xl mx-auto space-y-4">
        {loading ? (
            <div className="text-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
                <p className="mt-4 text-gray-500">{t('ranking.loading')}</p>
            </div>
        ) : rankings.length > 0 ? (
          rankings.map((user, index) => (
            <div
              key={user.user_id}
              className={`
                flex items-center justify-between p-4 rounded-xl border shadow-sm transition-all hover:shadow-md
                ${getMedalColor(index)}
              `}
            >
              <div className="flex items-center space-x-4">
                <div className="flex-shrink-0 w-12 h-12 flex items-center justify-center font-bold text-xl rounded-full bg-white bg-opacity-50">
                  {index + 1}
                </div>
                <div>
                  <h3 className="font-bold text-lg">{user.display_name}</h3>
                  <p className="text-xs opacity-75">@{user.username}</p>
                </div>
              </div>
              <div className="text-right">
                <span className="block text-2xl font-bold font-mono">
                  {formatDuration(user.total_duration)}
                </span>
              </div>
            </div>
          ))
        ) : (
          <div className="text-center py-12 bg-white rounded-xl border border-dashed border-gray-300">
            <p className="text-gray-500">{t('ranking.no_data')}</p>
          </div>
        )}
      </div>
    </div>
  )
}

export default RankingPage
