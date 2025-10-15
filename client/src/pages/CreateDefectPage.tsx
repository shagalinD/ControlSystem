import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import type { CreateDefectFormData, Project, User } from '../types'
import { defectService } from '../services/defectService'
import { projectService } from '../services/projectService'
import { userService } from '../services/userService'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { LoadingSpinner } from '../components/ui/LoadingSpinner'

export const CreateDefectPage: React.FC = () => {
  const navigate = useNavigate()
  const [projects, setProjects] = useState<Project[]>([])
  const [engineers, setEngineers] = useState<User[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [isDataLoading, setIsDataLoading] = useState(true)
  const [error, setError] = useState<string>('')

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateDefectFormData>()

  useEffect(() => {
    const loadInitialData = async () => {
      try {
        const [projectsResponse, engineersResponse] = await Promise.all([
          projectService.getProjects(),
          userService.getEngineers(),
        ])

        setProjects(projectsResponse.projects || []) // Исправляем здесь
        setEngineers(engineersResponse.data || [])
      } catch (err) {
        setError('Ошибка загрузки данных')
      } finally {
        setIsDataLoading(false)
      }
    }

    loadInitialData()
  }, [])

  const onSubmit = async (data: CreateDefectFormData) => {
    setIsLoading(true)
    setError('')

    try {
      await defectService.createDefect(data)
      navigate('/defects')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка создания дефекта')
    } finally {
      setIsLoading(false)
    }
  }

  if (isDataLoading) {
    return (
      <div className='flex justify-center items-center min-h-64'>
        <LoadingSpinner size='lg' />
      </div>
    )
  }

  return (
    <div className='max-w-2xl mx-auto'>
      <div className='mb-6'>
        <h1 className='text-2xl font-bold text-gray-900'>Создание дефекта</h1>
        <p className='text-gray-600 mt-1'>
          Заполните информацию о новом дефекте на строительном объекте
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
              label='Название дефекта *'
              {...register('title', {
                required: 'Название обязательно',
                minLength: {
                  value: 5,
                  message: 'Название должно быть не менее 5 символов',
                },
              })}
              error={errors.title?.message}
              placeholder='Например: Трещина в несущей стене'
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
                placeholder='Подробное описание дефекта, местоположение, возможные причины...'
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
                  Исполнитель
                </label>
                <select
                  {...register('assignee_id', {
                    setValueAs: (value) =>
                      value ? parseInt(value) : undefined,
                  })}
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

              <Input
                label='Срок выполнения'
                type='date'
                {...register('deadline')}
                min={new Date().toISOString().split('T')[0]}
              />
            </div>
          </div>
        </div>

        <div className='flex justify-end space-x-4'>
          <Button
            type='button'
            variant='secondary'
            onClick={() => navigate('/defects')}
          >
            Отмена
          </Button>
          <Button type='submit' isLoading={isLoading}>
            Создать дефект
          </Button>
        </div>
      </form>
    </div>
  )
}
