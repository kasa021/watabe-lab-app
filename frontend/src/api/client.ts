import axios from 'axios'

// APIクライアントの作成
export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// リクエストインターセプター
apiClient.interceptors.request.use(
  (config) => {
    // ローカルストレージからトークンを取得
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// レスポンスインターセプター
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // 401エラーの場合はログアウト処理
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/attendance/login'
    }
    return Promise.reject(error)
  }
)

