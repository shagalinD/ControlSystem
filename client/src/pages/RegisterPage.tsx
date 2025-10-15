import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import type { RegisterFormData } from '../types'
import { authService } from '../services/authService'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'

export const RegisterPage: React.FC = () => {
  const navigate = useNavigate()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string>('')
  const [success, setSuccess] = useState<string>('')

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm<RegisterFormData>()

  const watchPassword = watch('password', '')

  const onSubmit = async (formData: RegisterFormData) => {
    setIsLoading(true)
    setError('')
    setSuccess('')

    try {
      await authService.register(formData)
      setSuccess('Регистрация успешна! Теперь вы можете войти в систему.')

      // Перенаправляем на страницу логина через 2 секунды
      setTimeout(() => {
        navigate('/login')
      }, 2000)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка регистрации')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className='min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8'>
      <div className='sm:mx-auto sm:w-full sm:max-w-md'>
        <h2 className='mt-6 text-center text-3xl font-extrabold text-gray-900'>
          Регистрация в системе
        </h2>
        <p className='mt-2 text-center text-sm text-gray-600'>
          Создайте учетную запись для доступа к системе управления дефектами
        </p>
      </div>

      <div className='mt-8 sm:mx-auto sm:w-full sm:max-w-md'>
        <div className='bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10'>
          {success && (
            <div className='mb-4 p-3 bg-green-100 border border-green-400 text-green-700 rounded-lg'>
              {success}
            </div>
          )}

          {error && (
            <div className='mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded-lg'>
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit(onSubmit)} className='space-y-4'>
            <Input
              label='Полное имя'
              {...register('full_name', {
                required: 'Полное имя обязательно',
                minLength: {
                  value: 2,
                  message: 'Имя должно быть не менее 2 символов',
                },
              })}
              error={errors.full_name?.message}
              placeholder='Иван Петров'
            />

            <Input
              label='Email'
              type='email'
              {...register('email', {
                required: 'Email обязателен',
                pattern: {
                  value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                  message: 'Введите корректный email',
                },
              })}
              error={errors.email?.message}
              placeholder='user@example.com'
            />

            <Input
              label='Пароль'
              type='password'
              {...register('password', {
                required: 'Пароль обязателен',
                minLength: {
                  value: 6,
                  message: 'Пароль должен быть не менее 6 символов',
                },
              })}
              error={errors.password?.message}
              placeholder='Не менее 6 символов'
            />

            <Input
              label='Подтверждение пароля'
              type='password'
              {...register('confirm_password', {
                required: 'Подтверждение пароля обязательно',
                validate: (value) =>
                  value === watchPassword || 'Пароли не совпадают',
              })}
              error={errors.confirm_password?.message}
              placeholder='Повторите пароль'
            />

            <div>
              <label className='block text-sm font-medium text-gray-700 mb-1'>
                Роль *
              </label>
              <select
                {...register('role_id', {
                  required: 'Выберите роль',
                  setValueAs: (value) => parseInt(value),
                })}
                className={`
                  w-full px-3 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
                  ${errors.role_id ? 'border-red-500' : 'border-gray-300'}
                `}
              >
                <option value=''>Выберите роль</option>
                <option value='1'>Инженер</option>
                <option value='2'>Менеджер</option>
                <option value='3'>Наблюдатель</option>
              </select>
              {errors.role_id && (
                <p className='mt-1 text-sm text-red-600'>
                  {errors.role_id.message}
                </p>
              )}
              <p className='mt-1 text-xs text-gray-500'>
                • Инженер: создание и отслеживание дефектов
                <br />
                • Менеджер: управление проектами и отчетами
                <br />• Наблюдатель: просмотр информации
              </p>
            </div>

            <Button type='submit' isLoading={isLoading} className='w-full'>
              Зарегистрироваться
            </Button>
          </form>

          <div className='mt-6'>
            <div className='relative'>
              <div className='absolute inset-0 flex items-center'>
                <div className='w-full border-t border-gray-300' />
              </div>
              <div className='relative flex justify-center text-sm'>
                <span className='px-2 bg-white text-gray-500'>
                  Уже есть аккаунт?
                </span>
              </div>
            </div>

            <div className='mt-6 text-center'>
              <Link
                to='/login'
                className='font-medium text-blue-600 hover:text-blue-500'
              >
                Войти в существующий аккаунт
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
