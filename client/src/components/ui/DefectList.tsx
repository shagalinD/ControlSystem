import React from 'react'
import { Link } from 'react-router-dom'
import type { Defect, User } from '../../types'
import { DEFECT_STATUSES, DEFECT_PRIORITIES } from '../../constants'
import { Button } from '../ui/Button'

interface DefectListProps {
  defects: Defect[]
  onStatusChange: (defectId: number, newStatus: string) => void
  currentUser: User | null
}

export const DefectList: React.FC<DefectListProps> = ({
  defects,
  onStatusChange,
  currentUser,
}) => {
  const canEditDefect = (defect: Defect): boolean => {
    if (!currentUser) return false
    if (currentUser.role_name === 'manager') return true
    if (
      currentUser.role_name === 'engineer' &&
      defect.author_id === currentUser.id
    )
      return true
    return false
  }

  const getNextStatus = (currentStatus: string): string | null => {
    const statusFlow: Record<string, string> = {
      new: 'in_progress',
      in_progress: 'on_review',
      on_review: 'closed',
    }
    return statusFlow[currentStatus] || null
  }

  return (
    <div className='bg-white shadow-sm rounded-lg border border-gray-200'>
      <div className='overflow-x-auto'>
        <table className='min-w-full divide-y divide-gray-200'>
          <thead className='bg-gray-50'>
            <tr>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Дефект
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Проект
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Статус
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Приоритет
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Исполнитель
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Срок
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Действия
              </th>
            </tr>
          </thead>
          <tbody className='bg-white divide-y divide-gray-200'>
            {defects.map((defect) => {
              const statusInfo = DEFECT_STATUSES[defect.status]
              const priorityInfo = DEFECT_PRIORITIES[defect.priority]
              const nextStatus = getNextStatus(defect.status)
              const canEdit = canEditDefect(defect)

              return (
                <tr key={defect.id} className='hover:bg-gray-50'>
                  <td className='px-6 py-4 whitespace-nowrap'>
                    <div>
                      <Link
                        to={`/defects/${defect.id}`}
                        className='text-sm font-medium text-blue-600 hover:text-blue-900'
                      >
                        {defect.title}
                      </Link>
                      <p className='text-sm text-gray-500 truncate max-w-xs'>
                        {defect.description}
                      </p>
                    </div>
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-900'>
                    {defect.project.name}
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap'>
                    <span
                      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${statusInfo.color}`}
                    >
                      {statusInfo.label}
                    </span>
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap'>
                    <span
                      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${priorityInfo.color}`}
                    >
                      {priorityInfo.label}
                    </span>
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-900'>
                    {defect.assignee?.full_name || 'Не назначен'}
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-900'>
                    {defect.deadline
                      ? new Date(defect.deadline).toLocaleDateString('ru-RU')
                      : 'Не установлен'}
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2'>
                    <Link
                      to={`/defects/${defect.id}`}
                      className='text-blue-600 hover:text-blue-900'
                    >
                      Просмотр
                    </Link>

                    {canEdit && nextStatus && (
                      <Button
                        size='sm'
                        variant='secondary'
                        onClick={() => onStatusChange(defect.id, nextStatus)}
                      >
                        {
                          DEFECT_STATUSES[
                            nextStatus as keyof typeof DEFECT_STATUSES
                          ].label
                        }
                      </Button>
                    )}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}
