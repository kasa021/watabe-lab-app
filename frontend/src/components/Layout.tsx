import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useState } from 'react'

export const Layout = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const isLoggedIn = !!localStorage.getItem('token')
  const [isMenuOpen, setIsMenuOpen] = useState(false)

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    navigate('/login')
  }

  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen)
  }

  const closeMenu = () => {
    setIsMenuOpen(false)
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
              <Link to="/" className="flex-shrink-0 flex items-center" onClick={closeMenu}>
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
            <div className="hidden sm:flex sm:items-center">
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
            <div className="-mr-2 flex items-center sm:hidden">
              <button
                onClick={toggleMenu}
                type="button"
                className="bg-white inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                aria-controls="mobile-menu"
                aria-expanded="false"
              >
                <span className="sr-only">Open main menu</span>
                {!isMenuOpen ? (
                  <svg
                    className="block h-6 w-6"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    aria-hidden="true"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M4 6h16M4 12h16M4 18h16"
                    />
                  </svg>
                ) : (
                  <svg
                    className="block h-6 w-6"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    aria-hidden="true"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                )}
              </button>
            </div>
          </div>
        </div>

        {/* Mobile menu, show/hide based on menu state. */}
        {isLoggedIn && isMenuOpen && (
          <div className="sm:hidden" id="mobile-menu">
            <div className="pt-2 pb-3 space-y-1">
              <Link
                to="/"
                onClick={closeMenu}
                className={`${
                  location.pathname === '/'
                    ? 'bg-primary-50 border-primary-500 text-primary-700'
                    : 'border-transparent text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700'
                } block pl-3 pr-4 py-2 border-l-4 text-base font-medium`}
              >
                ホーム
              </Link>
              <Link
                to="/ranking"
                onClick={closeMenu}
                className={`${
                  location.pathname === '/ranking'
                    ? 'bg-primary-50 border-primary-500 text-primary-700'
                    : 'border-transparent text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700'
                } block pl-3 pr-4 py-2 border-l-4 text-base font-medium`}
              >
                ランキング
              </Link>
              <Link
                to="/achievements"
                onClick={closeMenu}
                className={`${
                  location.pathname === '/achievements'
                    ? 'bg-primary-50 border-primary-500 text-primary-700'
                    : 'border-transparent text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700'
                } block pl-3 pr-4 py-2 border-l-4 text-base font-medium`}
              >
                実績
              </Link>
              <Link
                to="/profile"
                onClick={closeMenu}
                className={`${
                  location.pathname === '/profile'
                    ? 'bg-primary-50 border-primary-500 text-primary-700'
                    : 'border-transparent text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700'
                } block pl-3 pr-4 py-2 border-l-4 text-base font-medium`}
              >
                プロフィール
              </Link>
            </div>
            <div className="pt-4 pb-4 border-t border-gray-200">
               <div className="mt-3 space-y-1">
                <button
                  onClick={() => {
                      handleLogout()
                      closeMenu()
                  }}
                  className="block w-full text-left pl-3 pr-4 py-2 border-l-4 border-transparent text-base font-medium text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700"
                >
                  ログアウト
                </button>
              </div>
            </div>
          </div>
        )}
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
