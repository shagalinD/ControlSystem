import React from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../../store/authStore'
import { Button } from '../ui/Button'

export const Header: React.FC = () => {
  const { user, logout } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <header className='bg-white shadow-sm border-b border-gray-200'>
      <div className='max-w-7xl mx-auto px-4 sm:px-6 lg:px-8'>
        <div className='flex justify-between items-center h-16'>
          <div className='flex items-center'>
            <Link to='/' className='flex-shrink-0'>
              <h1 className='text-xl font-bold text-gray-900'>
                СистемаКонтроля
              </h1>
            </Link>

            <nav className='hidden md:ml-8 md:flex space-x-4'>
              <Link
                to='/defects'
                className='text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium'
              >
                Дефекты
              </Link>
              <Link
                to='/projects'
                className='text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium'
              >
                Проекты
              </Link>
              {user?.role_name === 'manager' && (
                <Link
                  to='/reports'
                  className='text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium'
                >
                  Отчеты
                </Link>
              )}
            </nav>
          </div>

          <div className='flex items-center space-x-4'>
            <div className='hidden sm:flex flex-col text-right'>
              <span className='text-sm font-medium text-gray-900'>
                {user?.full_name}
              </span>
              <span className='text-xs text-gray-500 capitalize'>
                {user?.role_name}
              </span>
            </div>

            <Button variant='secondary' size='sm' onClick={handleLogout}>
              Выйти
            </Button>
          </div>
        </div>
      </div>
    </header>
  )
}
