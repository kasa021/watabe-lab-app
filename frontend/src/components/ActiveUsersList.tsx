import { useTranslation } from 'react-i18next'
import { useOccupancyStore } from '../stores/useOccupancyStore'

export const ActiveUsersList = () => {
    const { t } = useTranslation()
    const { activeUsers, isConnected } = useOccupancyStore()

    return (
        <div className="mt-12 w-full max-w-4xl mx-auto">
            <div className="flex items-center justify-between mb-6">
                <h3 className="text-2xl font-bold text-gray-800">
                    {t('active_users.title')}
                </h3>
                <span className={`px-3 py-1 rounded-full text-xs font-medium ${isConnected ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                    {isConnected ? t('active_users.live') : t('active_users.offline')}
                </span>
            </div>
            
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
                {activeUsers.map(log => (
                    <div key={log.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex items-center space-x-4 hover:shadow-md transition-shadow">
                        <div className="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-bold text-xl">
                            {log.user?.display_name?.charAt(0) || '?'}
                        </div>
                        <div className="text-left">
                            <p className="font-bold text-gray-900">{log.user?.display_name || t('active_users.unknown')}</p>
                            <p className="text-xs text-gray-500">
                                {new Date(log.check_in_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })} 〜
                            </p>
                        </div>
                    </div>
                ))}
            </div>

            {activeUsers.length === 0 && (
                <div className="text-center py-12 bg-white rounded-xl border border-dashed border-gray-300">
                    <p className="text-gray-500">{t('active_users.empty')}</p>
                </div>
            )}
        </div>
    )
}
