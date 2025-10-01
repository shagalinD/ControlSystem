import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import type { UserRole } from '../../types'
import { ROUTES } from '../../constants'

interface SidebarProps {
  role: UserRole | undefined
}

export const Sidebar: React.FC<SidebarProps> = ({ role }) => {
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
  ]

  const filteredNavItems = navItems.filter((item) =>
    item.roles.includes(role || 'observer')
  )

  return (
    <aside className='w-64 bg-white shadow-sm border-r border-gray-200'>
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
