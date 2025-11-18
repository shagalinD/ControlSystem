import React, { useState } from 'react'
import { useForm } from 'react-hook-form'
import type { UpdateProfileData } from '../types'
import { userService } from '../services/userService'
import { useAuthStore } from '../store/authStore'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { LoadingSpinner } from '../components/ui/LoadingSpinner'

export const ProfilePage: React.FC = () => {
  const { user, setUser } = useAuthStore()
  const [isLoading, setIsLoading] = useState(false)
  const [message, setMessage] = useState<string>('')
  const [error, setError] = useState<string>('')

  const profileForm = useForm<UpdateProfileData>({
    defaultValues: {
      full_name: user?.full_name || '',
      email: user?.email || '',
    },
  })

  const handleProfileUpdate = async (data: UpdateProfileData) => {
    if (!user) return

    setIsLoading(true)
    setError('')
    setMessage('')

    try {
      const response = await userService.updateUserProfile(user.id, data)
      setUser(response.user)
      setMessage('Профиль успешно обновлен')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка обновления профиля')
    } finally {
      setIsLoading(false)
    }
  }

  if (!user) {
    return (
      <div className='flex justify-center items-center min-h-64'>
        <LoadingSpinner size='lg' />
      </div>
    )
  }

  return (
    <div className='max-w-2xl mx-auto'>
      <div className='mb-6'>
        <h1 className='text-2xl font-bold text-gray-900'>Мой профиль</h1>
        <p className='text-gray-600 mt-1'>
          Управление личной информацией и настройками аккаунта
        </p>
      </div>

      {message && (
        <div className='bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-6'>
          {message}
        </div>
      )}

      {error && (
        <div className='bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6'>
          {error}
        </div>
      )}

      {/* Информация о пользователе */}
      <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200 mb-6'>
        <h3 className='text-lg font-semibold text-gray-900 mb-4'>
          Информация об аккаунте
        </h3>
        <div className='grid grid-cols-1 md:grid-cols-2 gap-4 text-sm'>
          <div>
            <span className='text-gray-600'>Роль:</span>
            <p className='font-medium capitalize'>{user.role_name}</p>
          </div>
          <div>
            <span className='text-gray-600'>ID пользователя:</span>
            <p className='font-medium'>{user.id}</p>
          </div>
          <div>
            <span className='text-gray-600'>Дата регистрации:</span>
            <p className='font-medium'>
              {user.created_at
                ? new Date(user.created_at).toLocaleDateString('ru-RU')
                : 'Неизвестно'}
            </p>
          </div>
        </div>
      </div>

      {/* Форма профиля */}

      <form
        onSubmit={profileForm.handleSubmit(handleProfileUpdate)}
        className='space-y-6'
      >
        <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
          <h3 className='text-lg font-semibold text-gray-900 mb-4'>
            Личная информация
          </h3>
          <div className='space-y-4'>
            <Input
              label='Полное имя'
              {...profileForm.register('full_name', {
                required: 'Полное имя обязательно',
                minLength: {
                  value: 2,
                  message: 'Имя должно быть не менее 2 символов',
                },
              })}
              error={profileForm.formState.errors.full_name?.message}
            />

            <Input
              label='Email'
              type='email'
              {...profileForm.register('email', {
                required: 'Email обязателен',
                pattern: {
                  value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                  message: 'Введите корректный email',
                },
              })}
              error={profileForm.formState.errors.email?.message}
            />
          </div>
        </div>

        <div className='flex justify-end'>
          <Button
            type='submit'
            isLoading={isLoading}
            disabled={!profileForm.formState.isDirty}
          >
            Сохранить изменения
          </Button>
        </div>
      </form>
    </div>
  )
}
