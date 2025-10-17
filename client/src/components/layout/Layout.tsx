import React from 'react'
import { Outlet } from 'react-router-dom'
import { Header } from './Header'

export const Layout: React.FC = () => {
  return (
    <div className='overflow-x-auto custom-scrollbar min-h-screen bg-gray-50'>
      <Header />
      <div className='flex'>
        <main className='flex-1 p-6'>
          <Outlet />
        </main>
      </div>
    </div>
  )
}
