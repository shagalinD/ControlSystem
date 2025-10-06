import React, { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import type { Project, Defect } from '../types'
import { projectService } from '../services/projectService'
import { defectService } from '../services/defectService'
import { useAuthStore } from '../store/authStore'
import { Button } from '../components/ui/Button'
import { LoadingSpinner } from '../components/ui/LoadingSpinner'
import { DefectList } from '../components/ui/DefectList'

export const ProjectDetailsPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const { user } = useAuthStore()
  const [project, setProject] = useState<Project | null>(null)
  const [defects, setDefects] = useState<Defect[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string>('')

  const loadProjectData = async () => {
    if (!id) return

    setIsLoading(true)
    setError('')

    try {
      const [projectResponse, defectsResponse] = await Promise.all([
        projectService.getProjectById(parseInt(id)),
        defectService.getDefects({ project_id: parseInt(id) }),
      ])

      setProject(projectResponse.project)
      setDefects(defectsResponse.data)
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Ошибка загрузки данных проекта'
      )
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    loadProjectData()
  }, [id])

  const handleStatusChange = async (defectId: number, newStatus: string) => {
    try {
      await defectService.updateDefectStatus(defectId, newStatus as any)
      await loadProjectData() // Перезагружаем данные
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

  if (!project) {
    return (
      <div className='text-center py-12'>
        <div className='text-gray-400 text-6xl mb-4'>❌</div>
        <h3 className='text-lg font-medium text-gray-900 mb-2'>
          Проект не найден
        </h3>
        <p className='text-gray-600 mb-4'>
          Запрошенный проект не существует или был удален
        </p>
        <Link to='/projects'>
          <Button>Вернуться к списку проектов</Button>
        </Link>
      </div>
    )
  }

  const canCreateDefect =
    user?.role_name === 'engineer' || user?.role_name === 'manager'

  return (
    <div className='space-y-6'>
      <div className='flex justify-between items-start'>
        <div>
          <h1 className='text-2xl font-bold text-gray-900'>{project.name}</h1>
          <p className='text-gray-600 mt-1'>{project.description}</p>
          <div className='flex items-center mt-2 text-sm text-gray-500'>
            <span>Менеджер: {project.manager.full_name}</span>
            <span className='mx-2'>•</span>
            <span>
              Создан: {new Date(project.created_at).toLocaleDateString('ru-RU')}
            </span>
          </div>
        </div>

        {canCreateDefect && (
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

      <div className='grid grid-cols-1 lg:grid-cols-3 gap-6'>
        <div className='lg:col-span-2'>
          <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
            <div className='flex justify-between items-center mb-4'>
              <h2 className='text-lg font-semibold text-gray-900'>
                Дефекты проекта ({defects.length})
              </h2>
            </div>

            {defects.length > 0 ? (
              <DefectList
                defects={defects}
                onStatusChange={handleStatusChange}
                currentUser={user}
              />
            ) : (
              <div className='text-center py-8'>
                <div className='text-gray-400 text-4xl mb-3'>🐛</div>
                <p className='text-gray-600'>
                  В этом проекте пока нет дефектов
                </p>
                {canCreateDefect && (
                  <Link to='/defects/create' className='inline-block mt-3'>
                    <Button size='sm'>Создать первый дефект</Button>
                  </Link>
                )}
              </div>
            )}
          </div>
        </div>

        <div className='space-y-6'>
          <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
            <h3 className='text-lg font-semibold text-gray-900 mb-4'>
              Статистика
            </h3>
            <div className='space-y-3'>
              <div className='flex justify-between'>
                <span className='text-gray-600'>Всего дефектов:</span>
                <span className='font-medium'>{defects.length}</span>
              </div>
              <div className='flex justify-between'>
                <span className='text-gray-600'>В работе:</span>
                <span className='font-medium text-blue-600'>
                  {defects.filter((d) => d.status === 'in_progress').length}
                </span>
              </div>
              <div className='flex justify-between'>
                <span className='text-gray-600'>На проверке:</span>
                <span className='font-medium text-yellow-600'>
                  {defects.filter((d) => d.status === 'on_review').length}
                </span>
              </div>
              <div className='flex justify-between'>
                <span className='text-gray-600'>Закрыто:</span>
                <span className='font-medium text-green-600'>
                  {defects.filter((d) => d.status === 'closed').length}
                </span>
              </div>
              <div className='flex justify-between'>
                <span className='text-gray-600'>Критических:</span>
                <span className='font-medium text-red-600'>
                  {defects.filter((d) => d.priority === 'critical').length}
                </span>
              </div>
            </div>
          </div>

          <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
            <h3 className='text-lg font-semibold text-gray-900 mb-4'>
              Информация
            </h3>
            <div className='space-y-2 text-sm'>
              <div>
                <span className='text-gray-600'>Менеджер:</span>
                <p className='font-medium'>{project.manager.full_name}</p>
                <p className='text-gray-500'>{project.manager.email}</p>
              </div>
              <div>
                <span className='text-gray-600'>Дата создания:</span>
                <p className='font-medium'>
                  {new Date(project.created_at).toLocaleDateString('ru-RU')}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
