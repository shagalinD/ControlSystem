import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import type { LoginFormData } from '../types'
import { authService } from '../services/authService'
import { useAuthStore } from '../store/authStore'

export const LoginPage: React.FC = () => {
  const navigate = useNavigate()
  const { login } = useAuthStore()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string>('')
  const [formData, setFormData] = useState<LoginFormData>({
    email: '',
    password: '',
  })

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }))
    if (error) setError('')
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    console.log('Login form submitted - preventing default')

    if (!formData.email || !formData.password) {
      setError('Заполните все поля')
      return
    }

    setIsLoading(true)
    setError('')

    try {
      const { token, user } = await authService.login(formData)
      login(token, user)
      navigate('/')
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : 'Ошибка авторизации'
      setError(errorMessage)
      console.error('Login error:', errorMessage)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className='min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8'>
      <div className='max-w-md w-full space-y-8'>
        <div>
          <h2 className='mt-6 text-center text-3xl font-extrabold text-gray-900'>
            Вход в систему
          </h2>
          <p className='mt-2 text-center text-sm text-gray-600'>
            Система управления дефектами строительных объектов
          </p>
        </div>
        {/* Простая форма без лишних оберток */}
        <form onSubmit={handleSubmit} className='mt-8 space-y-6'>
          {error && (
            <div className='bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded'>
              {error}
            </div>
          )}

          <div className='rounded-md shadow-sm -space-y-px'>
            <div>
              <label htmlFor='email' className='sr-only'>
                Email
              </label>
              <input
                id='email'
                name='email'
                type='email'
                required
                value={formData.email}
                onChange={handleInputChange}
                className='appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm'
                placeholder='Email'
              />
            </div>
            <div>
              <label htmlFor='password' className='sr-only'>
                Пароль
              </label>
              <input
                id='password'
                name='password'
                type='password'
                required
                value={formData.password}
                onChange={handleInputChange}
                className='appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm'
                placeholder='Пароль'
              />
            </div>
          </div>

          <div>
            <button
              type='submit'
              disabled={isLoading}
              className='group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50'
            >
              {isLoading ? 'Вход...' : 'Войти'}
            </button>
          </div>
        </form>
        <div className='mt-6 text-center'>
          <p className='text-sm text-gray-600'>
            Нет аккаунта?{' '}
            <Link
              to='/register'
              className='font-medium text-blue-600 hover:text-blue-500'
            >
              Зарегистрироваться
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}
