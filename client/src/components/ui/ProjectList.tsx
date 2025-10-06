import React from 'react'
import { Link } from 'react-router-dom'
import type { Project } from '../../types'
import { useAuthStore } from '../../store/authStore'

interface ProjectListProps {
  projects: Project[]
}

export const ProjectList: React.FC<ProjectListProps> = ({ projects }) => {
  const { user } = useAuthStore()

  const canEditProject = user?.role_name === 'manager'

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  }

  const getDefectStats = (project: Project) => {
    const defects = project.defects || []
    const total = defects.length
    const closed = defects.filter((d) => d.status === 'closed').length
    const critical = defects.filter((d) => d.priority === 'critical').length

    return { total, closed, critical }
  }

  return (
    <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'>
      {projects.map((project) => {
        const stats = getDefectStats(project)

        return (
          <div
            key={project.id}
            className='bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow'
          >
            <div className='p-6'>
              <div className='flex justify-between items-start mb-4'>
                <div>
                  <Link
                    to={`/projects/${project.id}`}
                    className='text-lg font-semibold text-gray-900 hover:text-blue-600'
                  >
                    {project.name}
                  </Link>
                  <p className='text-sm text-gray-500 mt-1'>
                    Менеджер: {project.manager.full_name}
                  </p>
                </div>
                {canEditProject && (
                  <Link
                    to={`/projects/${project.id}/edit`}
                    className='text-gray-400 hover:text-gray-600'
                  >
                    ✏️
                  </Link>
                )}
              </div>

              <p className='text-gray-600 text-sm mb-4 line-clamp-2'>
                {project.description || 'Описание отсутствует'}
              </p>

              <div className='grid grid-cols-3 gap-2 mb-4 text-center'>
                <div className='bg-gray-50 rounded p-2'>
                  <div className='text-sm font-medium text-gray-900'>
                    {stats.total}
                  </div>
                  <div className='text-xs text-gray-500'>Всего</div>
                </div>
                <div className='bg-green-50 rounded p-2'>
                  <div className='text-sm font-medium text-green-900'>
                    {stats.closed}
                  </div>
                  <div className='text-xs text-green-700'>Закрыто</div>
                </div>
                <div className='bg-red-50 rounded p-2'>
                  <div className='text-sm font-medium text-red-900'>
                    {stats.critical}
                  </div>
                  <div className='text-xs text-red-700'>Критич.</div>
                </div>
              </div>

              <div className='flex justify-between items-center text-xs text-gray-500'>
                <span>Создан: {formatDate(project.created_at)}</span>
                <Link
                  to={`/projects/${project.id}`}
                  className='text-blue-600 hover:text-blue-800 font-medium'
                >
                  Подробнее →
                </Link>
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}
