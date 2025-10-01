import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import type { Defect, DefectFilters } from '../types'
import { defectService } from '../services/defectService'
import { useAuthStore } from '../store/authStore'
import { Button } from '../components/ui/Button'
import { LoadingSpinner } from '../components/ui/LoadingSpinner'
import { DefectFilters as DefectFiltersComponent } from '../components/ui/DefectFilters'
import { DefectList } from '../components/ui/DefectList'

export const DefectsPage: React.FC = () => {
  const { user } = useAuthStore()
  const [defects, setDefects] = useState<Defect[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string>('')
  const [filters, setFilters] = useState<DefectFilters>({
    page: 1,
    page_size: 20,
  })

  const loadDefects = async () => {
    setIsLoading(true)
    setError('')

    try {
      const response = await defectService.getDefects(filters)
      setDefects(response.data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка загрузки дефектов')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    loadDefects()
  }, [filters])

  const handleFiltersChange = (newFilters: DefectFilters) => {
    setFilters({ ...newFilters, page: 1 })
  }

  const handleStatusChange = async (defectId: number, newStatus: string) => {
    try {
      await defectService.updateDefectStatus(defectId, newStatus as any)
      await loadDefects() // Перезагружаем список
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка обновления статуса')
    }
  }

  if (isLoading) {
    return (
      <div className='flex justify-center items-center min-h-64'>
        <LoadingSpinner size='lg' />
      </div>
    )
  }

  return (
    <div className='space-y-6'>
      <div className='flex justify-between items-center'>
        <div>
          <h1 className='text-2xl font-bold text-gray-900'>Дефекты</h1>
          <p className='text-gray-600 mt-1'>
            Список всех дефектов строительных объектов
          </p>
        </div>

        {(user?.role_name === 'engineer' || user?.role_name === 'manager') && (
          <Link to='/defects/create'>
            <Button>Создать дефект</Button>
          </Link>
        )}
      </div>

      {error && (
        <div className='bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded'>
          {error}
        </div>
      )}

      <DefectFiltersComponent
        filters={filters}
        onFiltersChange={handleFiltersChange}
      />

      <DefectList
        defects={defects}
        onStatusChange={handleStatusChange}
        currentUser={user}
      />

      {defects.length === 0 && !isLoading && (
        <div className='text-center py-12'>
          <div className='text-gray-400 text-6xl mb-4'>🐛</div>
          <h3 className='text-lg font-medium text-gray-900 mb-2'>
            Дефекты не найдены
          </h3>
          <p className='text-gray-600 mb-4'>
            Попробуйте изменить параметры фильтрации или создайте первый дефект
          </p>
          {(user?.role_name === 'engineer' ||
            user?.role_name === 'manager') && (
            <Link to='/defects/create'>
              <Button>Создать первый дефект</Button>
            </Link>
          )}
        </div>
      )}
    </div>
  )
}
