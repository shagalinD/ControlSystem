import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import type { LoginFormData } from '../types'
import { useAuthStore } from '../store/authStore'
import { LoginForm } from '../components/ui/LoginForm'
import { authService } from '../services/authService'

export const LoginPage: React.FC = () => {
  const navigate = useNavigate()
  const { login } = useAuthStore()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string>('')

  const handleLogin = async (formData: LoginFormData) => {
    setIsLoading(true)
    setError('')

    try {
      const { token, user } = await authService.login(formData)
      login(token, user)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка авторизации')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className='min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8'>
      <div className='sm:mx-auto sm:w-full sm:max-w-md'>
        <h2 className='mt-6 text-center text-3xl font-extrabold text-gray-900'>
          Вход в систему
        </h2>
        <p className='mt-2 text-center text-sm text-gray-600'>
          Система управления дефектами строительных объектов
        </p>
      </div>

      <div className='mt-8 sm:mx-auto sm:w-full sm:max-w-md'>
        <div className='bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10'>
          <LoginForm
            onSubmit={handleLogin}
            isLoading={isLoading}
            error={error}
          />

          <div className='mt-6'>
            <div className='relative'>
              <div className='absolute inset-0 flex items-center'>
                <div className='w-full border-t border-gray-300' />
              </div>
              <div className='relative flex justify-center text-sm'>
                <span className='px-2 bg-white text-gray-500'>
                  Тестовые данные
                </span>
              </div>
            </div>

            <div className='mt-4 text-xs text-gray-500 space-y-1'>
              <div>Менеджер: manager@company.com / password123</div>
              <div>Инженер: engineer@company.com / password123</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
