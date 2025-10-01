import React from 'react'
import { Outlet } from 'react-router-dom'
import { useAuthStore } from '../../store/authStore'
import { Header } from './Header'
import { Sidebar } from './Sidebar'

export const Layout: React.FC = () => {
  const { user } = useAuthStore()

  return (
    <div className='min-h-screen bg-gray-50'>
      <Header />
      <div className='flex'>
        <Sidebar role={user?.role_name} />
        <main className='flex-1 p-6'>
          <Outlet />
        </main>
      </div>
    </div>
  )
}
