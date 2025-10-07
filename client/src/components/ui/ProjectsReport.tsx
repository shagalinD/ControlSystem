import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import type { Project, ProjectReport } from '../../types'
import { projectService } from '../../services/projectService'
import { reportService } from '../../services/reportService'
import { LoadingSpinner } from '../ui/LoadingSpinner'

export const ProjectsReport: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([])
  const [projectReports, setProjectReports] = useState<
    Record<number, ProjectReport>
  >({})
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string>('')

  const loadProjectsReport = async () => {
    setIsLoading(true)
    setError('')

    try {
      const projectsResponse = await projectService.getProjects()
      setProjects(projectsResponse.data)

      // Загружаем отчеты для каждого проекта
      const reports: Record<number, ProjectReport> = {}
      for (const project of projectsResponse.data) {
        try {
          const reportResponse = await reportService.getProjectReport(
            project.id
          )
          reports[project.id] = reportResponse.report
        } catch (err) {
          console.error(`Error loading report for project ${project.id}:`, err)
        }
      }
      setProjectReports(reports)
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : 'Ошибка загрузки отчетов по проектам'
      )
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    loadProjectsReport()
  }, [])

  if (isLoading) {
    return (
      <div className='flex justify-center items-center min-h-64'>
        <LoadingSpinner size='lg' />
      </div>
    )
  }

  if (error) {
    return (
      <div className='bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded'>
        {error}
      </div>
    )
  }

  const getCompletionColor = (percentage: number) => {
    if (percentage >= 80) return 'text-green-600'
    if (percentage >= 60) return 'text-yellow-600'
    return 'text-red-600'
  }

  return (
    <div className='space-y-6'>
      <div className='grid grid-cols-1 md:grid-cols-3 gap-6 mb-6'>
        <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200 text-center'>
          <div className='text-2xl font-bold text-gray-900'>
            {projects.length}
          </div>
          <div className='text-gray-600 text-sm mt-1'>Всего проектов</div>
        </div>

        <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200 text-center'>
          <div className='text-2xl font-bold text-blue-600'>
            {
              projects.filter((p) => {
                const report = projectReports[p.id]
                return report && report.completion_percentage >= 80
              }).length
            }
          </div>
          <div className='text-gray-600 text-sm mt-1'>
            Проектов с высокой завершенностью
          </div>
        </div>

        <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200 text-center'>
          <div className='text-2xl font-bold text-red-600'>
            {
              projects.filter((p) => {
                const report = projectReports[p.id]
                return (
                  (report &&
                    report.defects_by_priority.find(
                      (d) => d.priority === 'critical'
                    )?.count) ||
                  0 > 0
                )
              }).length
            }
          </div>
          <div className='text-gray-600 text-sm mt-1'>
            Проектов с критическими дефектами
          </div>
        </div>
      </div>

      <div className='bg-white rounded-lg shadow-sm border border-gray-200'>
        <div className='px-6 py-4 border-b border-gray-200'>
          <h3 className='text-lg font-semibold text-gray-900'>
            Детальная статистика по проектам
          </h3>
        </div>

        <div className='overflow-x-auto'>
          <table className='min-w-full divide-y divide-gray-200'>
            <thead className='bg-gray-50'>
              <tr>
                <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                  Проект
                </th>
                <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                  Всего дефектов
                </th>
                <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                  Критические
                </th>
                <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                  В работе
                </th>
                <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                  Завершенность
                </th>
                <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                  Действия
                </th>
              </tr>
            </thead>
            <tbody className='bg-white divide-y divide-gray-200'>
              {projects.map((project) => {
                const report = projectReports[project.id]

                if (!report) {
                  return (
                    <tr key={project.id}>
                      <td className='px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900'>
                        {project.name}
                      </td>
                      <td
                        colSpan={5}
                        className='px-6 py-4 text-sm text-gray-500'
                      >
                        Загрузка данных...
                      </td>
                    </tr>
                  )
                }

                const criticalDefects =
                  report.defects_by_priority.find(
                    (d) => d.priority === 'critical'
                  )?.count || 0
                const inProgressDefects =
                  report.defects_by_status.find(
                    (d) => d.status === 'in_progress'
                  )?.count || 0

                return (
                  <tr key={project.id} className='hover:bg-gray-50'>
                    <td className='px-6 py-4 whitespace-nowrap'>
                      <Link
                        to={`/projects/${project.id}`}
                        className='text-sm font-medium text-blue-600 hover:text-blue-900'
                      >
                        {project.name}
                      </Link>
                      <p className='text-sm text-gray-500 truncate max-w-xs'>
                        {project.description}
                      </p>
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-900'>
                      {report.total_defects}
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap'>
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          criticalDefects > 0
                            ? 'bg-red-100 text-red-800'
                            : 'bg-green-100 text-green-800'
                        }`}
                      >
                        {criticalDefects}
                      </span>
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-900'>
                      {inProgressDefects}
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap'>
                      <div className='flex items-center space-x-2'>
                        <div className='w-16 bg-gray-200 rounded-full h-2'>
                          <div
                            className='bg-blue-600 h-2 rounded-full'
                            style={{
                              width: `${report.completion_percentage}%`,
                            }}
                          />
                        </div>
                        <span
                          className={`text-sm font-medium ${getCompletionColor(
                            report.completion_percentage
                          )}`}
                        >
                          {report.completion_percentage}%
                        </span>
                      </div>
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap text-sm font-medium'>
                      <Link
                        to={`/projects/${project.id}`}
                        className='text-blue-600 hover:text-blue-900 mr-4'
                      >
                        Подробнее
                      </Link>
                      <Link
                        to={`/reports/projects/${project.id}`}
                        className='text-gray-600 hover:text-gray-900'
                      >
                        Отчет
                      </Link>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      </div>

      {projects.length === 0 && (
        <div className='text-center py-12'>
          <div className='text-gray-400 text-6xl mb-4'>📊</div>
          <h3 className='text-lg font-medium text-gray-900 mb-2'>
            Нет данных для отчета
          </h3>
          <p className='text-gray-600'>
            Создайте проекты и дефекты для получения аналитики
          </p>
        </div>
      )}
    </div>
  )
}
