import React from 'react'
import { useAuthStore } from '../store/authStore'
import { Button } from '../components/ui/Button'
import { DEFECT_STATUSES, DEFECT_PRIORITIES } from '../constants'

export const HomePage: React.FC = () => {
  const { user, logout } = useAuthStore()

  return (
    <div className='min-h-screen bg-gray-50'>
      <div className='max-w-7xl mx-auto py-6 sm:px-6 lg:px-8'>
        <div className='px-4 py-6 sm:px-0'>
          <div className='border-4 border-dashed border-gray-200 rounded-lg p-8'>
            <div className='flex justify-between items-center mb-8'>
              <div>
                <h1 className='text-3xl font-bold text-gray-900'>
                  Добро пожаловать, {user?.full_name}!
                </h1>
                <p className='text-gray-600 mt-2'>
                  Роль: {user?.role_name} • Email: {user?.email}
                </p>
              </div>
              <Button variant='secondary' onClick={logout}>
                Выйти
              </Button>
            </div>

            <div className='grid grid-cols-1 md:grid-cols-2 gap-6 mb-8'>
              <div className='bg-white p-6 rounded-lg shadow'>
                <h2 className='text-xl font-semibold mb-4'>Статусы дефектов</h2>
                <div className='space-y-2'>
                  {Object.values(DEFECT_STATUSES).map((status) => (
                    <div key={status.value} className='flex items-center'>
                      <span
                        className={`inline-block w-3 h-3 rounded-full ${
                          status.color.split(' ')[0]
                        } mr-2`}
                      />
                      <span>{status.label}</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className='bg-white p-6 rounded-lg shadow'>
                <h2 className='text-xl font-semibold mb-4'>
                  Приоритеты дефектов
                </h2>
                <div className='space-y-2'>
                  {Object.values(DEFECT_PRIORITIES).map((priority) => (
                    <div key={priority.value} className='flex items-center'>
                      <span
                        className={`inline-block w-3 h-3 rounded-full ${
                          priority.color.split(' ')[0]
                        } mr-2`}
                      />
                      <span>{priority.label}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <div className='bg-white p-6 rounded-lg shadow'>
              <h2 className='text-xl font-semibold mb-4'>Быстрые действия</h2>
              <div className='flex flex-wrap gap-4'>
                <Button>Просмотреть дефекты</Button>
                <Button>Создать дефект</Button>
                {user?.role_name === 'manager' && (
                  <Button>Управление проектами</Button>
                )}
                <Button variant='secondary'>Отчеты</Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
