import React, { useEffect, useRef, useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useAuthStore } from '../../store/authStore'
import { Sidebar } from './Sidebar'
import { FiMenu } from 'react-icons/fi'

export const Header: React.FC = () => {
  const { user } = useAuthStore()
  const [isMenuOpen, setIsMenuOpen] = useState<boolean>(false)
  const asideRef = useRef<HTMLDivElement>(null)
  const buttonRef = useRef<HTMLDivElement>(null)
  const location = useLocation()

  useEffect(() => {
    setIsMenuOpen(false)
  }, [location.pathname])

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        asideRef.current &&
        !asideRef.current.contains(event.target as Node) &&
        buttonRef.current &&
        !buttonRef.current.contains(event.target as Node)
      ) {
        setIsMenuOpen(false)
      }
    }

    // Закрытие меню при нажатии Escape
    const handleEscapeKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setIsMenuOpen(false)
      }
    }

    if (isMenuOpen) {
      document.addEventListener('mousedown', handleClickOutside)
      document.addEventListener('keydown', handleEscapeKey)
      // Блокируем скролл на body когда меню открыто
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = 'unset'
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
      document.removeEventListener('keydown', handleEscapeKey)
      document.body.style.overflow = 'unset'
    }
  }, [isMenuOpen])

  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen)
  }

  return (
    <>
      <div>
        <Sidebar
          isMenuOpen={isMenuOpen}
          ref={asideRef}
          role={user?.role_name}
        />
      </div>

      <header className='bg-white shadow-sm border-b border-gray-200'>
        <div className='max-w-7xl mx-auto px-4 sm:px-6 lg:px-8'>
          <div className='flex justify-between items-center h-16 '>
            <div className='flex items-center align-middle'>
              <div
                className=' hidden max-md:block p-3 rounded-full'
                onClick={toggleMenu}
                ref={buttonRef}
              >
                <FiMenu className='flex-shrink-0 size-5' />
              </div>
              <Link to='/' className='flex-shrink-0'>
                <h1 className='text-xl font-bold text-gray-900'>
                  СистемаКонтроля
                </h1>
              </Link>

              <nav className='hidden md:ml-8 md:flex space-x-4 items-center mr-5'>
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
                <Link
                  to='/profile'
                  className='text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium'
                >
                  Мой профиль
                </Link>
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
            </div>
          </div>
        </div>
      </header>
    </>
  )
}
