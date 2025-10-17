import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import type { UserRole } from '../../types'
import { ROUTES } from '../../constants'

interface SidebarProps {
  role: UserRole | undefined
  ref: React.Ref<HTMLDivElement>
  isMenuOpen: boolean
}

export const Sidebar: React.FC<SidebarProps> = ({ role, ref, isMenuOpen }) => {
  const location = useLocation()

  const isActive = (path: string) => {
    return location.pathname === path
  }

  const navItems = [
    {
      path: ROUTES.HOME,
      label: 'Главная',
      icon: '🏠',
      roles: ['engineer', 'manager', 'observer'],
    },
    {
      path: ROUTES.DEFECTS,
      label: 'Все дефекты',
      icon: '🐛',
      roles: ['engineer', 'manager', 'observer'],
    },
    {
      path: ROUTES.DEFECTS_CREATE,
      label: 'Создать дефект',
      icon: '➕',
      roles: ['engineer', 'manager'],
    },
    {
      path: ROUTES.PROJECTS,
      label: 'Проекты',
      icon: '🏢',
      roles: ['manager', 'observer'],
    },
    {
      path: ROUTES.REPORTS,
      label: 'Отчеты',
      icon: '📊',
      roles: ['manager', 'observer'],
    },
    {
      path: ROUTES.PROFILE,
      label: 'Мой профиль',
      icon: '👤',
      roles: ['engineer', 'manager', 'observer'],
    },
  ]

  const filteredNavItems = navItems.filter((item) =>
    item.roles.includes(role || 'observer')
  )

  return (
    <aside
      ref={ref}
      className={`fixed top-0 left-0 shadow-xl z-50 transform transition-transform duration-300 ease-in-out h-full w-64 bg-white border-r border-gray-200 ${
        isMenuOpen ? 'translate-x-0' : '-translate-x-full'
      }`}
    >
      <nav className='p-4 space-y-2'>
        {filteredNavItems.map((item) => (
          <Link
            key={item.path}
            to={item.path}
            className={`
              flex items-center px-3 py-2 text-sm font-medium rounded-lg transition-colors
              ${
                isActive(item.path)
                  ? 'bg-blue-100 text-blue-700 border-r-2 border-blue-600'
                  : 'text-gray-700 hover:bg-gray-100'
              }
            `}
          >
            <span className='mr-3 text-lg'>{item.icon}</span>
            {item.label}
          </Link>
        ))}
      </nav>
    </aside>
  )
}
