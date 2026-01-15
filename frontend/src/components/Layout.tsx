import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom'

export const Layout = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const isLoggedIn = !!localStorage.getItem('token')

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    navigate('/login')
  }

  if (!isLoggedIn && location.pathname !== '/login') {
    // ログインページ以外で未ログインの場合はヘッダーをシンプルにするか、リダイレクトする
    // ここではシンプルにヘッダーを表示（ログインボタンのみなど）
  }

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <header className="bg-white shadow-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex">
              <Link to="/" className="flex-shrink-0 flex items-center">
                <span className="text-xl font-bold text-primary-600">Watabe Lab</span>
              </Link>
              {isLoggedIn && (
                <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
                  <Link
                    to="/"
                    className={`${
                      location.pathname === '/'
                        ? 'border-primary-500 text-gray-900'
                        : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
                    } inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium`}
                  >
                    ホーム
                  </Link>
                  <Link
                    to="/ranking"
                    className={`${
                      location.pathname === '/ranking'
                        ? 'border-primary-500 text-gray-900'
                        : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
                    } inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium`}
                  >
                    ランキング
                  </Link>
                  <Link
                    to="/achievements"
                    className={`${
                      location.pathname === '/achievements'
                        ? 'border-primary-500 text-gray-900'
                        : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
                    } inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium`}
                  >
                    実績
                  </Link>
                  <Link
                    to="/profile"
                    className={`${
                      location.pathname === '/profile'
                        ? 'border-primary-500 text-gray-900'
                        : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
                    } inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium`}
                  >
                    プロフィール
                  </Link>
                </div>
              )}
            </div>
            <div className="flex items-center">
              {isLoggedIn ? (
                <button
                  onClick={handleLogout}
                  className="text-gray-500 hover:text-gray-700 text-sm font-medium"
                >
                  ログアウト
                </button>
              ) : (
                <Link
                  to="/login"
                  className="text-primary-600 hover:text-primary-700 text-sm font-medium"
                >
                  ログイン
                </Link>
              )}
            </div>
          </div>
        </div>
      </header>

      <main className="flex-grow">
        <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
          <Outlet />
        </div>
      </main>

      <footer className="bg-white border-t border-gray-200 mt-auto">
        <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
          <p className="text-center text-sm text-gray-500">
            &copy; {new Date().getFullYear()} Watabe Lab Attendance System
          </p>
        </div>
      </footer>
    </div>
  )
}
