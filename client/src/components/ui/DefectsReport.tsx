import React, { useState, useEffect } from 'react'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts'
import type {
  DefectFilters,
  DefectsReport as DefectsReportType,
} from '../../types'
import { reportService } from '../../services/reportService'
import { LoadingSpinner } from '../ui/LoadingSpinner'
import { DEFECT_STATUSES, DEFECT_PRIORITIES } from '../../constants'

interface DefectsReportProps {
  filters: DefectFilters
  onFiltersChange: (filters: DefectFilters) => void
}

const COLORS = {
  status: ['#6B7280', '#3B82F6', '#F59E0B', '#10B981', '#EF4444'],
  priority: ['#6B7280', '#F59E0B', '#F97316', '#EF4444', '#7C2D12'],
}

export const DefectsReport: React.FC<DefectsReportProps> = ({ filters }) => {
  const [report, setReport] = useState<DefectsReportType | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string>('')

  const loadReport = async () => {
    setIsLoading(true)
    setError('')

    try {
      const response = await reportService.getDefectsReport(filters)
      setReport(response.report)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка загрузки отчета')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    loadReport()
  }, [filters])

  if (isLoading) {
    return (
      <div className='flex justify-center items-center min-h-64'>
        <LoadingSpinner size='lg' />
      </div>
    )
  }

  if (error) {
    return (
      <div className='bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded'>
        {error}
      </div>
    )
  }

  if (!report) {
    return (
      <div className='text-center py-12'>
        <p className='text-gray-500'>Нет данных для отображения</p>
      </div>
    )
  }

  // Подготовка данных для графиков
  const statusData = report.defects_by_status.map((item, index) => ({
    name:
      DEFECT_STATUSES[item.status as keyof typeof DEFECT_STATUSES]?.label ||
      item.status,
    value: item.count,
    color: COLORS.status[index % COLORS.status.length],
  }))

  const priorityData = report.defects_by_priority.map((item, index) => ({
    name:
      DEFECT_PRIORITIES[item.priority as keyof typeof DEFECT_PRIORITIES]
        ?.label || item.priority,
    value: item.count,
    color: COLORS.priority[index % COLORS.priority.length],
  }))

  return (
    <div className='space-y-8'>
      {/* Ключевые метрики */}
      <div className='grid grid-cols-1 md:grid-cols-4 gap-6'>
        <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200 text-center'>
          <div className='text-2xl font-bold text-gray-900'>
            {report.total_defects}
          </div>
          <div className='text-gray-600 text-sm mt-1'>Всего дефектов</div>
        </div>

        <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200 text-center'>
          <div className='text-2xl font-bold text-red-600'>
            {report.overdue_defects}
          </div>
          <div className='text-gray-600 text-sm mt-1'>Просрочено</div>
        </div>

        <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200 text-center'>
          <div className='text-2xl font-bold text-green-600'>
            {report.defects_by_status.find((s) => s.status === 'closed')
              ?.count || 0}
          </div>
          <div className='text-gray-600 text-sm mt-1'>Закрыто</div>
        </div>

        <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200 text-center'>
          <div className='text-2xl font-bold text-blue-600'>
            {report.avg_resolution_time.toFixed(1)}
          </div>
          <div className='text-gray-600 text-sm mt-1'>
            Ср. время решения (ч)
          </div>
        </div>
      </div>

      {/* Графики */}
      <div className='grid grid-cols-1 lg:grid-cols-2 gap-8'>
        {/* Распределение по статусам */}
        <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
          <h3 className='text-lg font-semibold text-gray-900 mb-4'>
            Распределение по статусам
          </h3>
          <div className='h-80'>
            <ResponsiveContainer width='100%' height='100%'>
              <PieChart>
                <Pie
                  data={statusData}
                  cx='50%'
                  cy='50%'
                  labelLine={false}
                  label={({ name, value, percent }) =>
                    `${name}: ${value} (${((percent as number) * 100).toFixed(
                      1
                    )}%)`
                  }
                  outerRadius={80}
                  fill='#8884d8'
                  dataKey='value'
                >
                  {statusData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Распределение по приоритетам */}
        <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
          <h3 className='text-lg font-semibold text-gray-900 mb-4'>
            Распределение по приоритетам
          </h3>
          <div className='h-80'>
            <ResponsiveContainer width='100%' height='100%'>
              <BarChart data={priorityData}>
                <CartesianGrid strokeDasharray='3 3' />
                <XAxis dataKey='name' />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey='value' name='Количество дефектов'>
                  {priorityData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* Детальная таблица статусов */}
      <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
        <h3 className='text-lg font-semibold text-gray-900 mb-4'>
          Детальная статистика по статусам
        </h3>
        <div className='overflow-x-auto'>
          <table className='min-w-full divide-y divide-gray-200'>
            <thead className='bg-gray-50'>
              <tr>
                <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                  Статус
                </th>
                <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                  Количество
                </th>
                <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                  Процент
                </th>
              </tr>
            </thead>
            <tbody className='bg-white divide-y divide-gray-200'>
              {report.defects_by_status.map((item) => {
                const statusInfo =
                  DEFECT_STATUSES[item.status as keyof typeof DEFECT_STATUSES]
                const percentage = (item.count / report.total_defects) * 100

                return (
                  <tr key={item.status}>
                    <td className='px-6 py-4 whitespace-nowrap'>
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          statusInfo?.color || 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {statusInfo?.label || item.status}
                      </span>
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-900'>
                      {item.count}
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-900'>
                      {percentage.toFixed(1)}%
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
