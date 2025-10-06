import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import type { Defect, UpdateDefectFormData, Project, User } from '../types'
import { defectService } from '../services/defectService'
import { projectService } from '../services/projectService'
import { userService } from '../services/userService'
import { useAuthStore } from '../store/authStore'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { LoadingSpinner } from '../components/ui/LoadingSpinner'

export const EditDefectPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [defect, setDefect] = useState<Defect | null>(null)
  const [projects, setProjects] = useState<Project[]>([])
  const [engineers, setEngineers] = useState<User[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string>('')

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<UpdateDefectFormData>()

  useEffect(() => {
    const loadData = async () => {
      if (!id) return

      setIsLoading(true)
      setError('')

      try {
        const [defectResponse, projectsResponse, engineersResponse] =
          await Promise.all([
            defectService.getDefectById(parseInt(id)),
            projectService.getProjects(),
            userService.getEngineers(),
          ])

        const defectData = defectResponse.defect
        setDefect(defectData)
        setProjects(projectsResponse.data)
        setEngineers(engineersResponse.data)

        // Заполняем форму данными дефекта
        reset({
          title: defectData.title,
          description: defectData.description,
          status: defectData.status,
          priority: defectData.priority,
          deadline: defectData.deadline
            ? defectData.deadline.split('T')[0]
            : undefined,
          project_id: defectData.project_id,
          assignee_id: defectData.assignee_id || undefined,
        })
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Ошибка загрузки данных')
      } finally {
        setIsLoading(false)
      }
    }

    loadData()
  }, [id, reset])

  const onSubmit = async (data: UpdateDefectFormData) => {
    if (!defect) return

    setIsSubmitting(true)
    setError('')

    try {
      // Преобразуем данные для отправки
      const updateData: any = { ...data }

      // Преобразуем assignee_id в число, если есть значение
      if (updateData.assignee_id) {
        updateData.assignee_id = parseInt(updateData.assignee_id)
      } else {
        updateData.assignee_id = null
      }

      // Преобразуем project_id в число
      updateData.project_id = parseInt(updateData.project_id)

      await defectService.updateDefect(defect.id, updateData)
      navigate(`/defects/${defect.id}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка обновления дефекта')
    } finally {
      setIsSubmitting(false)
    }
  }

  const canEditDefect =
    user?.role_name === 'manager' ||
    (user?.role_name === 'engineer' && defect?.author_id === user.id)

  // Если пользователь не может редактировать этот дефект
  if (!isLoading && !canEditDefect) {
    return (
      <div className='text-center py-12'>
        <div className='text-gray-400 text-6xl mb-4'>🚫</div>
        <h3 className='text-lg font-medium text-gray-900 mb-2'>
          Доступ запрещен
        </h3>
        <p className='text-gray-600 mb-4'>
          У вас нет прав для редактирования этого дефекта
        </p>
        <Button onClick={() => navigate(`/defects/${id}`)}>
          Вернуться к дефекту
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

  if (!defect) {
    return (
      <div className='text-center py-12'>
        <div className='text-gray-400 text-6xl mb-4'>❌</div>
        <h3 className='text-lg font-medium text-gray-900 mb-2'>
          Дефект не найден
        </h3>
        <p className='text-gray-600 mb-4'>
          Запрошенный дефект не существует или был удален
        </p>
        <Button onClick={() => navigate('/defects')}>
          Вернуться к списку дефектов
        </Button>
      </div>
    )
  }

  return (
    <div className='max-w-2xl mx-auto'>
      <div className='mb-6'>
        <h1 className='text-2xl font-bold text-gray-900'>
          Редактирование дефекта
        </h1>
        <p className='text-gray-600 mt-1'>Обновите информацию о дефекте</p>
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
              label='Название дефекта *'
              {...register('title', {
                required: 'Название обязательно',
                minLength: {
                  value: 5,
                  message: 'Название должно быть не менее 5 символов',
                },
              })}
              error={errors.title?.message}
            />

            <div>
              <label className='block text-sm font-medium text-gray-700 mb-1'>
                Описание *
              </label>
              <textarea
                {...register('description', {
                  required: 'Описание обязательно',
                  minLength: {
                    value: 10,
                    message: 'Описание должно быть не менее 10 символов',
                  },
                })}
                rows={4}
                className={`
                  w-full px-3 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
                  ${errors.description ? 'border-red-500' : 'border-gray-300'}
                `}
              />
              {errors.description && (
                <p className='mt-1 text-sm text-red-600'>
                  {errors.description.message}
                </p>
              )}
            </div>

            <div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
              <div>
                <label className='block text-sm font-medium text-gray-700 mb-1'>
                  Проект *
                </label>
                <select
                  {...register('project_id', {
                    required: 'Выберите проект',
                  })}
                  className={`
                    w-full px-3 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
                    ${errors.project_id ? 'border-red-500' : 'border-gray-300'}
                  `}
                >
                  <option value=''>Выберите проект</option>
                  {projects.map((project) => (
                    <option key={project.id} value={project.id}>
                      {project.name}
                    </option>
                  ))}
                </select>
                {errors.project_id && (
                  <p className='mt-1 text-sm text-red-600'>
                    {errors.project_id.message}
                  </p>
                )}
              </div>

              <div>
                <label className='block text-sm font-medium text-gray-700 mb-1'>
                  Приоритет *
                </label>
                <select
                  {...register('priority', {
                    required: 'Выберите приоритет',
                  })}
                  className={`
                    w-full px-3 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
                    ${errors.priority ? 'border-red-500' : 'border-gray-300'}
                  `}
                >
                  <option value=''>Выберите приоритет</option>
                  <option value='low'>Низкий</option>
                  <option value='medium'>Средний</option>
                  <option value='high'>Высокий</option>
                  <option value='critical'>Критический</option>
                </select>
                {errors.priority && (
                  <p className='mt-1 text-sm text-red-600'>
                    {errors.priority.message}
                  </p>
                )}
              </div>
            </div>

            <div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
              <div>
                <label className='block text-sm font-medium text-gray-700 mb-1'>
                  Статус *
                </label>
                <select
                  {...register('status', {
                    required: 'Выберите статус',
                  })}
                  className={`
                    w-full px-3 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
                    ${errors.status ? 'border-red-500' : 'border-gray-300'}
                  `}
                >
                  <option value='new'>Новая</option>
                  <option value='in_progress'>В работе</option>
                  <option value='on_review'>На проверке</option>
                  <option value='closed'>Закрыта</option>
                  <option value='cancelled'>Отменена</option>
                </select>
                {errors.status && (
                  <p className='mt-1 text-sm text-red-600'>
                    {errors.status.message}
                  </p>
                )}
              </div>

              <div>
                <label className='block text-sm font-medium text-gray-700 mb-1'>
                  Исполнитель
                </label>
                <select
                  {...register('assignee_id')}
                  className='w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500'
                >
                  <option value=''>Не назначен</option>
                  {engineers.map((engineer) => (
                    <option key={engineer.id} value={engineer.id}>
                      {engineer.full_name}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <Input
              label='Срок выполнения'
              type='date'
              {...register('deadline')}
              min={new Date().toISOString().split('T')[0]}
            />
          </div>
        </div>

        <div className='flex justify-end space-x-4'>
          <Button
            type='button'
            variant='secondary'
            onClick={() => navigate(`/defects/${defect.id}`)}
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
