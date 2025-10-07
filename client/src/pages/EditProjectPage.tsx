import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import type { Project, UpdateProjectData, User } from '../types'
import { projectService } from '../services/projectService'
import { userService } from '../services/userService'
import { useAuthStore } from '../store/authStore'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { LoadingSpinner } from '../components/ui/LoadingSpinner'

export const EditProjectPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [project, setProject] = useState<Project | null>(null)
  const [managers, setManagers] = useState<User[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string>('')

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<UpdateProjectData>()

  useEffect(() => {
    const loadData = async () => {
      if (!id) return

      setIsLoading(true)
      setError('')

      try {
        const [projectResponse, managersResponse] = await Promise.all([
          projectService.getProjectById(parseInt(id)),
          userService.getManagers(),
        ])

        const projectData = projectResponse.project
        setProject(projectData)
        setManagers(managersResponse.data)

        // Заполняем форму данными проекта
        reset({
          name: projectData.name,
          description: projectData.description,
          manager_id: projectData.manager_id,
        })
      } catch (err) {
        setError(
          err instanceof Error ? err.message : 'Ошибка загрузки данных проекта'
        )
      } finally {
        setIsLoading(false)
      }
    }

    loadData()
  }, [id, reset])

  const onSubmit = async (data: UpdateProjectData) => {
    if (!project) return

    setIsSubmitting(true)
    setError('')

    try {
      // Преобразуем manager_id в число
      const updateData: any = { ...data }
      if (updateData.manager_id) {
        updateData.manager_id = parseInt(updateData.manager_id)
      }

      await projectService.updateProject(project.id, updateData)
      navigate(`/projects/${project.id}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка обновления проекта')
    } finally {
      setIsSubmitting(false)
    }
  }

  // Проверяем права доступа
  if (!isLoading && user?.role_name !== 'manager') {
    return (
      <div className='text-center py-12'>
        <div className='text-gray-400 text-6xl mb-4'>🚫</div>
        <h3 className='text-lg font-medium text-gray-900 mb-2'>
          Доступ запрещен
        </h3>
        <p className='text-gray-600 mb-4'>
          Только менеджеры могут редактировать проекты
        </p>
        <Button onClick={() => navigate(`/projects/${id}`)}>
          Вернуться к проекту
        </Button>
      </div>
    )
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
        <Button onClick={() => navigate('/projects')}>
          Вернуться к списку проектов
        </Button>
      </div>
    )
  }

  return (
    <div className='max-w-2xl mx-auto'>
      <div className='mb-6'>
        <h1 className='text-2xl font-bold text-gray-900'>
          Редактирование проекта
        </h1>
        <p className='text-gray-600 mt-1'>
          Обновите информацию о проекте "{project.name}"
        </p>
      </div>

      {error && (
        <div className='bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6'>
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className='space-y-6'>
        <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
          <div className='space-y-4'>
            <Input
              label='Название проекта *'
              {...register('name', {
                required: 'Название проекта обязательно',
                minLength: {
                  value: 3,
                  message: 'Название должно быть не менее 3 символов',
                },
              })}
              error={errors.name?.message}
            />

            <div>
              <label className='block text-sm font-medium text-gray-700 mb-1'>
                Описание проекта
              </label>
              <textarea
                {...register('description')}
                rows={4}
                className='w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500'
                placeholder='Описание строительного объекта, местоположение, особенности...'
              />
            </div>

            <div>
              <label className='block text-sm font-medium text-gray-700 mb-1'>
                Ответственный менеджер *
              </label>
              <select
                {...register('manager_id', {
                  required: 'Выберите ответственного менеджера',
                  setValueAs: (value) => parseInt(value),
                })}
                className={`
                  w-full px-3 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
                  ${errors.manager_id ? 'border-red-500' : 'border-gray-300'}
                `}
              >
                <option value=''>Выберите менеджера</option>
                {managers.map((manager) => (
                  <option
                    key={manager.id}
                    value={manager.id}
                    selected={manager.id === project.manager_id}
                  >
                    {manager.full_name} ({manager.email})
                  </option>
                ))}
              </select>
              {errors.manager_id && (
                <p className='mt-1 text-sm text-red-600'>
                  {errors.manager_id.message}
                </p>
              )}
            </div>
          </div>
        </div>

        {/* Информация о проекте */}
        <div className='bg-gray-50 p-6 rounded-lg border border-gray-200'>
          <h3 className='text-lg font-semibold text-gray-900 mb-4'>
            Информация о проекте
          </h3>
          <div className='grid grid-cols-1 md:grid-cols-2 gap-4 text-sm'>
            <div>
              <span className='text-gray-600'>Дата создания:</span>
              <p className='font-medium'>
                {new Date(project.created_at).toLocaleDateString('ru-RU')}
              </p>
            </div>
            <div>
              <span className='text-gray-600'>Текущий менеджер:</span>
              <p className='font-medium'>{project.manager.full_name}</p>
            </div>
            <div>
              <span className='text-gray-600'>Количество дефектов:</span>
              <p className='font-medium'>{project.defects?.length || 0}</p>
            </div>
            <div>
              <span className='text-gray-600'>ID проекта:</span>
              <p className='font-medium'>{project.id}</p>
            </div>
          </div>
        </div>

        <div className='flex justify-end space-x-4'>
          <Button
            type='button'
            variant='secondary'
            onClick={() => navigate(`/projects/${project.id}`)}
          >
            Отмена
          </Button>
          <Button type='submit' isLoading={isSubmitting}>
            Сохранить изменения
          </Button>
        </div>
      </form>
    </div>
  )
}
