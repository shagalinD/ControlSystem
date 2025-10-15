import React from 'react'
import { Link } from 'react-router-dom'
import type { Defect, User } from '../../types'
import { DEFECT_STATUSES, DEFECT_PRIORITIES } from '../../constants'

interface DefectListProps {
  defects: Defect[]
  onStatusChange: (defectId: number, newStatus: string) => void
  currentUser: User | null
}

export const DefectList: React.FC<DefectListProps> = ({ defects }) => {
  return (
    <div className='bg-white shadow-sm rounded-lg border border-gray-200'>
      <div className='overflow-x-auto'>
        <table className='min-w-full divide-y divide-gray-200'>
          <thead className='bg-gray-50'>
            <tr>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Дефект
              </th>
              <th className='hidden xl:table-cell px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Проект
              </th>
              <th className='hidden lg:table-cell px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Статус
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Приоритет
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Срок
              </th>
            </tr>
          </thead>
          <tbody className='bg-white divide-y divide-gray-200'>
            {defects.map((defect) => {
              const statusInfo = DEFECT_STATUSES[defect.status]
              const priorityInfo = DEFECT_PRIORITIES[defect.priority]

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
                  <td className='hidden xl:table-cell px-6 py-4 whitespace-nowrap text-sm text-gray-900'>
                    {defect.project.name}
                  </td>
                  <td className='hidden lg:table-cell px-6 py-4 whitespace-nowrap'>
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
                    {defect.deadline
                      ? new Date(defect.deadline).toLocaleDateString('ru-RU')
                      : 'Не установлен'}
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
