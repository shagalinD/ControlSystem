import React from 'react'
import type {
  DefectFilters as DefectFiltersType,
  DefectStatus,
  DefectPriority,
} from '../../types'
import { DEFECT_STATUSES, DEFECT_PRIORITIES } from '../../constants'
import { Button } from '../ui/Button'
import { Input } from '../ui/Input'

interface DefectFiltersProps {
  filters: DefectFiltersType
  onFiltersChange: (filters: DefectFiltersType) => void
}

export const DefectFilters: React.FC<DefectFiltersProps> = ({
  filters,
  onFiltersChange,
}) => {
  const handleStatusChange = (status: DefectStatus | '') => {
    onFiltersChange({
      ...filters,
      status: status || undefined,
      page: 1,
    })
  }

  const handlePriorityChange = (priority: DefectPriority | '') => {
    onFiltersChange({
      ...filters,
      priority: priority || undefined,
      page: 1,
    })
  }

  const handleSearchChange = (search: string) => {
    onFiltersChange({
      ...filters,
      search: search || undefined,
      page: 1,
    })
  }

  const clearFilters = () => {
    onFiltersChange({
      page: 1,
      page_size: 20,
    })
  }

  const hasActiveFilters = Boolean(
    filters.status || filters.priority || filters.search
  )

  return (
    <div className='bg-white p-4 rounded-lg shadow-sm border border-gray-200'>
      <div className='grid grid-cols-1 md:grid-cols-4 gap-4'>
        <div>
          <label className='block text-sm font-medium text-gray-700 mb-1'>
            Поиск
          </label>
          <Input
            type='text'
            placeholder='Поиск по названию...'
            value={filters.search || ''}
            onChange={(e) => handleSearchChange(e.target.value)}
          />
        </div>

        <div>
          <label className='block text-sm font-medium text-gray-700 mb-1'>
            Статус
          </label>
          <select
            className='w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500'
            value={filters.status || ''}
            onChange={(e) =>
              handleStatusChange(e.target.value as DefectStatus | '')
            }
          >
            <option value=''>Все статусы</option>
            {Object.values(DEFECT_STATUSES).map((status) => (
              <option key={status.value} value={status.value}>
                {status.label}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className='block text-sm font-medium text-gray-700 mb-1'>
            Приоритет
          </label>
          <select
            className='w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500'
            value={filters.priority || ''}
            onChange={(e) =>
              handlePriorityChange(e.target.value as DefectPriority | '')
            }
          >
            <option value=''>Все приоритеты</option>
            {Object.values(DEFECT_PRIORITIES).map((priority) => (
              <option key={priority.value} value={priority.value}>
                {priority.label}
              </option>
            ))}
          </select>
        </div>

        <div className='flex items-end'>
          {hasActiveFilters && (
            <Button
              variant='secondary'
              onClick={clearFilters}
              className='w-full'
            >
              Сбросить фильтры
            </Button>
          )}
        </div>
      </div>
    </div>
  )
}
