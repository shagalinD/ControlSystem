import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import type { CreateProjectData, User } from '../types'
import { projectService } from '../services/projectService'
import { userService } from '../services/userService'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { LoadingSpinner } from '../components/ui/LoadingSpinner'

export const CreateProjectPage: React.FC = () => {
  const navigate = useNavigate()
  const [managers, setManagers] = useState<User[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [isDataLoading, setIsDataLoading] = useState(true)
  const [error, setError] = useState<string>('')

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateProjectData>()

  useEffect(() => {
    const loadManagers = async () => {
      try {
        const response = await userService.getManagers()
        setManagers(response.managers || []) // Исправляем здесь
      } catch (err) {
        setError('Ошибка загрузки списка менеджеров')
      } finally {
        setIsDataLoading(false)
      }
    }

    loadManagers()
  }, [])

  const onSubmit = async (data: CreateProjectData) => {
    setIsLoading(true)
    setError('')

    try {
      await projectService.createProject(data)
      navigate('/projects')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка создания проекта')
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
        <h1 className='text-2xl font-bold text-gray-900'>Создание проекта</h1>
        <p className='text-gray-600 mt-1'>
          Добавьте новый строительный объект в систему
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
              placeholder='Например: ЖК Северный'
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
                  <option key={manager.id} value={manager.id}>
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

        <div className='flex justify-end space-x-4'>
          <Button
            type='button'
            variant='secondary'
            onClick={() => navigate('/projects')}
          >
            Отмена
          </Button>
          <Button type='submit' isLoading={isLoading}>
            Создать проект
          </Button>
        </div>
      </form>
    </div>
  )
}
