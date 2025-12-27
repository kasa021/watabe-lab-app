import { useEffect, useState } from 'react'
import { apiClient } from '../api/client'

interface HealthResponse {
  status: string
  message: string
}

const HomePage = () => {
  const [health, setHealth] = useState<HealthResponse | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const checkHealth = async () => {
      try {
        const response = await apiClient.get<HealthResponse>('/health')
        setHealth(response.data)
      } catch (error) {
        console.error('Health check failed:', error)
      } finally {
        setLoading(false)
      }
    }

    checkHealth()
  }, [])

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold text-gray-900 mb-8">
          研究室出席管理システム
        </h1>

        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <h2 className="text-2xl font-semibold text-gray-800 mb-4">
            システムステータス
          </h2>
          {loading ? (
            <p className="text-gray-600">確認中...</p>
          ) : health ? (
            <div className="space-y-2">
              <p className="text-green-600 font-medium">
                ✅ {health.message}
              </p>
              <p className="text-gray-600">
                ステータス: <span className="font-mono">{health.status}</span>
              </p>
            </div>
          ) : (
            <p className="text-red-600">
              ❌ バックエンドサーバーに接続できません
            </p>
          )}
        </div>

        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-semibold text-gray-800 mb-4">
            機能一覧
          </h2>
          <ul className="space-y-3 text-gray-700">
            <li className="flex items-center">
              <span className="mr-2">📝</span>
              <span>チェックイン/チェックアウト機能</span>
            </li>
            <li className="flex items-center">
              <span className="mr-2">👥</span>
              <span>リアルタイム在室者表示</span>
            </li>
            <li className="flex items-center">
              <span className="mr-2">🏆</span>
              <span>ランキング機能</span>
            </li>
            <li className="flex items-center">
              <span className="mr-2">🎖️</span>
              <span>称号システム</span>
            </li>
            <li className="flex items-center">
              <span className="mr-2">📊</span>
              <span>出席履歴の可視化</span>
            </li>
          </ul>
        </div>

        <div className="mt-8 text-center">
          <a
            href="/login"
            className="inline-block bg-primary-600 hover:bg-primary-700 text-white font-semibold px-8 py-3 rounded-lg transition-colors"
          >
            ログイン
          </a>
        </div>
      </div>
    </div>
  )
}

export default HomePage

