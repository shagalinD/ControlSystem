import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import type { Project, ProjectFilters } from '../types'
import { projectService } from '../services/projectService'
import { useAuthStore } from '../store/authStore'
import { Button } from '../components/ui/Button'
import { LoadingSpinner } from '../components/ui/LoadingSpinner'
import { ProjectList } from '../components/ui/ProjectList'

export const ProjectsPage: React.FC = () => {
  const { user } = useAuthStore()
  const [projects, setProjects] = useState<Project[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string>('')
  const [filters, setFilters] = useState<ProjectFilters>({
    page: 1,
    page_size: 20,
  })

  const loadProjects = async () => {
    setIsLoading(true)
    setError('')

    try {
      const response = await projectService.getProjects(filters)
      // Исправляем здесь: сервер возвращает { projects: [], pagination: {} }
      setProjects(response.projects || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка загрузки проектов')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    loadProjects()
  }, [filters])

  const canCreateProject = user?.role_name === 'manager'

  if (isLoading) {
    return (
      <div className='flex justify-center items-center min-h-64'>
        <LoadingSpinner size='lg' />
      </div>
    )
  }

  return (
    <div className='space-y-6'>
      <div className='flex justify-between items-center'>
        <div>
          <h1 className='text-2xl font-bold text-gray-900'>Проекты</h1>
          <p className='text-gray-600 mt-1'>
            Список строительных объектов и проектов
          </p>
        </div>

        {canCreateProject && (
          <Link to='/projects/create'>
            <Button>Создать проект</Button>
          </Link>
        )}
      </div>

      {error && (
        <div className='bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded'>
          {error}
        </div>
      )}

      <ProjectList projects={projects} />

      {projects.length === 0 && !isLoading && (
        <div className='text-center py-12'>
          <div className='text-gray-400 text-6xl mb-4'>🏢</div>
          <h3 className='text-lg font-medium text-gray-900 mb-2'>
            Проекты не найдены
          </h3>
          <p className='text-gray-600 mb-4'>
            {canCreateProject
              ? 'Создайте первый проект для управления строительными объектами'
              : 'На данный момент нет доступных проектов'}
          </p>
          {canCreateProject && (
            <Link to='/projects/create'>
              <Button>Создать первый проект</Button>
            </Link>
          )}
        </div>
      )}
    </div>
  )
}
