import React, { useState } from 'react'
import type { DefectFilters } from '../types'
import { useAuthStore } from '../store/authStore'
import { DefectsReport } from '../components/ui/DefectsReport'
import { ProjectsReport } from '../components/ui/ProjectsReport'
import { ExportPanel } from '../components/ui/ExportPanel'

export const ReportsPage: React.FC = () => {
  const { user } = useAuthStore()
  const [activeTab, setActiveTab] = useState<'defects' | 'projects'>('defects')
  const [filters, setFilters] = useState<DefectFilters>({})
  const [error] = useState<string>('')

  // Проверяем права доступа
  if (user?.role_name !== 'manager' && user?.role_name !== 'observer') {
    return (
      <div className='text-center py-12'>
        <div className='text-gray-400 text-6xl mb-4'>🚫</div>
        <h3 className='text-lg font-medium text-gray-900 mb-2'>
          Доступ запрещен
        </h3>
        <p className='text-gray-600'>
          Только менеджеры и наблюдатели имеют доступ к отчетам
        </p>
      </div>
    )
  }

  return (
    <div className='space-y-6'>
      <div className='flex justify-between items-center'>
        <div>
          <h1 className='text-2xl font-bold text-gray-900'>
            Аналитика и отчеты
          </h1>
          <p className='text-gray-600 mt-1'>
            Статистика по дефектам и проектам для принятия управленческих
            решений
          </p>
        </div>

        <ExportPanel filters={filters} />
      </div>

      {error && (
        <div className='bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded'>
          {error}
        </div>
      )}

      {/* Табы */}
      <div className='border-b border-gray-200'>
        <nav className='-mb-px flex space-x-8'>
          <button
            onClick={() => setActiveTab('defects')}
            className={`
              py-2 px-1 border-b-2 font-medium text-sm
              ${
                activeTab === 'defects'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }
            `}
          >
            Статистика по дефектам
          </button>
          <button
            onClick={() => setActiveTab('projects')}
            className={`
              py-2 px-1 border-b-2 font-medium text-sm
              ${
                activeTab === 'projects'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }
            `}
          >
            Отчеты по проектам
          </button>
        </nav>
      </div>

      {/* Контент табов */}
      <div className='min-h-96'>
        {activeTab === 'defects' && (
          <DefectsReport filters={filters} onFiltersChange={setFilters} />
        )}

        {activeTab === 'projects' && <ProjectsReport />}
      </div>
    </div>
  )
}
