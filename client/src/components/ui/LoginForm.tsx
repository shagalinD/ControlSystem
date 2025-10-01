import React from 'react'
import { useForm } from 'react-hook-form'
import type { LoginFormData } from '../../types'
import { Button } from '../ui/Button'
import { Input } from '../ui/Input'

interface LoginFormProps {
  onSubmit: (data: LoginFormData) => void
  isLoading?: boolean
  error?: string
}

export const LoginForm: React.FC<LoginFormProps> = ({
  onSubmit,
  isLoading = false,
  error,
}) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>()

  return (
    <form onSubmit={handleSubmit(onSubmit)} className='space-y-4'>
      {error && (
        <div className='p-3 bg-red-100 border border-red-400 text-red-700 rounded-lg'>
          {error}
        </div>
      )}

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
        placeholder='Введите пароль'
      />

      <Button type='submit' isLoading={isLoading} className='w-full'>
        Войти
      </Button>
    </form>
  )
}
