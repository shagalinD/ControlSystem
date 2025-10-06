import React, { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import type { Defect, Comment } from '../types'
import { defectService } from '../services/defectService'
import { commentService } from '../services/commentService'
import { useAuthStore } from '../store/authStore'
import { Button } from '../components/ui/Button'
import { LoadingSpinner } from '../components/ui/LoadingSpinner'
import { DEFECT_STATUSES, DEFECT_PRIORITIES } from '../constants'

export const DefectDetailsPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const { user } = useAuthStore()
  const [defect, setDefect] = useState<Defect | null>(null)
  const [comments, setComments] = useState<Comment[]>([])
  const [newComment, setNewComment] = useState('')
  const [isLoading, setIsLoading] = useState(true)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string>('')

  const loadDefectData = async () => {
    if (!id) return

    setIsLoading(true)
    setError('')

    try {
      const [defectResponse, commentsResponse] = await Promise.all([
        defectService.getDefectById(parseInt(id)),
        commentService.getDefectComments(parseInt(id)),
      ])

      setDefect(defectResponse.defect)
      setComments(commentsResponse.comments)
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Ошибка загрузки данных дефекта'
      )
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    loadDefectData()
  }, [id])

  const handleAddComment = async () => {
    if (!newComment.trim() || !defect) return

    setIsSubmitting(true)
    try {
      const response = await commentService.createComment(defect.id, newComment)
      setComments((prev) => [response.comment, ...prev])
      setNewComment('')
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Ошибка добавления комментария'
      )
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleStatusChange = async (newStatus: string) => {
    if (!defect) return

    try {
      await defectService.updateDefectStatus(defect.id, newStatus as any)
      await loadDefectData() // Перезагружаем данные
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка обновления статуса')
    }
  }

  const canEditDefect =
    user?.role_name === 'manager' ||
    (user?.role_name === 'engineer' && defect?.author_id === user.id)

  const getNextStatus = (currentStatus: string): string | null => {
    const statusFlow: Record<string, string> = {
      new: 'in_progress',
      in_progress: 'on_review',
      on_review: 'closed',
    }
    return statusFlow[currentStatus] || null
  }

  if (isLoading) {
    return (
      <div className='flex justify-center items-center min-h-64'>
        <LoadingSpinner size='lg' />
      </div>
    )
  }

  if (!defect) {
    return (
      <div className='text-center py-12'>
        <div className='text-gray-400 text-6xl mb-4'>❌</div>
        <h3 className='text-lg font-medium text-gray-900 mb-2'>
          Дефект не найден
        </h3>
        <p className='text-gray-600 mb-4'>
          Запрошенный дефект не существует или был удален
        </p>
        <Link to='/defects'>
          <Button>Вернуться к списку дефектов</Button>
        </Link>
      </div>
    )
  }

  const statusInfo = DEFECT_STATUSES[defect.status]
  const priorityInfo = DEFECT_PRIORITIES[defect.priority]
  const nextStatus = getNextStatus(defect.status)

  return (
    <div className='max-w-6xl mx-auto space-y-6'>
      <div className='flex justify-between items-start'>
        <div>
          <h1 className='text-2xl font-bold text-gray-900'>{defect.title}</h1>
          <div className='flex items-center mt-2 space-x-4 text-sm text-gray-600'>
            <span>Проект: {defect.project.name}</span>
            <span>•</span>
            <span>
              Создан: {new Date(defect.created_at).toLocaleDateString('ru-RU')}
            </span>
            <span>•</span>
            <span>Автор: {defect.author.full_name}</span>
          </div>
        </div>

        <div className='flex space-x-3'>
          {canEditDefect && (
            <Link to={`/defects/${defect.id}/edit`}>
              <Button variant='secondary'>Редактировать</Button>
            </Link>
          )}
          <Link to='/defects'>
            <Button variant='secondary'>Назад к списку</Button>
          </Link>
        </div>
      </div>

      {error && (
        <div className='bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded'>
          {error}
        </div>
      )}

      <div className='grid grid-cols-1 lg:grid-cols-3 gap-6'>
        {/* Основная информация */}
        <div className='lg:col-span-2 space-y-6'>
          <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
            <h2 className='text-lg font-semibold text-gray-900 mb-4'>
              Описание дефекта
            </h2>
            <p className='text-gray-700 whitespace-pre-wrap'>
              {defect.description}
            </p>
          </div>

          {/* Комментарии */}
          <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
            <h2 className='text-lg font-semibold text-gray-900 mb-4'>
              Комментарии ({comments.length})
            </h2>

            {/* Форма добавления комментария */}
            {(user?.role_name === 'engineer' ||
              user?.role_name === 'manager') && (
              <div className='mb-6'>
                <textarea
                  value={newComment}
                  onChange={(e) => setNewComment(e.target.value)}
                  placeholder='Добавьте комментарий...'
                  rows={3}
                  className='w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500'
                />
                <div className='flex justify-end mt-2'>
                  <Button
                    onClick={handleAddComment}
                    isLoading={isSubmitting}
                    disabled={!newComment.trim()}
                  >
                    Добавить комментарий
                  </Button>
                </div>
              </div>
            )}

            {/* Список комментариев */}
            <div className='space-y-4'>
              {comments.length === 0 ? (
                <p className='text-gray-500 text-center py-4'>
                  Пока нет комментариев
                </p>
              ) : (
                comments.map((comment) => (
                  <div
                    key={comment.id}
                    className='border-b border-gray-200 pb-4 last:border-0'
                  >
                    <div className='flex justify-between items-start mb-2'>
                      <div className='font-medium text-gray-900'>
                        {comment.author.full_name}
                      </div>
                      <div className='text-sm text-gray-500'>
                        {new Date(comment.created_at).toLocaleString('ru-RU')}
                      </div>
                    </div>
                    <p className='text-gray-700 whitespace-pre-wrap'>
                      {comment.text}
                    </p>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>

        {/* Боковая панель */}
        <div className='space-y-6'>
          {/* Статус и приоритет */}
          <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
            <h3 className='text-lg font-semibold text-gray-900 mb-4'>
              Информация
            </h3>
            <div className='space-y-4'>
              <div>
                <label className='block text-sm font-medium text-gray-700 mb-1'>
                  Статус
                </label>
                <div className='flex items-center justify-between'>
                  <span
                    className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${statusInfo.color}`}
                  >
                    {statusInfo.label}
                  </span>
                  {nextStatus &&
                    (user?.role_name === 'manager' ||
                      defect.assignee_id === user?.id) && (
                      <Button
                        size='sm'
                        onClick={() => handleStatusChange(nextStatus)}
                      >
                        {
                          DEFECT_STATUSES[
                            nextStatus as keyof typeof DEFECT_STATUSES
                          ].label
                        }
                      </Button>
                    )}
                </div>
              </div>

              <div>
                <label className='block text-sm font-medium text-gray-700 mb-1'>
                  Приоритет
                </label>
                <span
                  className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${priorityInfo.color}`}
                >
                  {priorityInfo.label}
                </span>
              </div>

              <div>
                <label className='block text-sm font-medium text-gray-700 mb-1'>
                  Исполнитель
                </label>
                <p className='text-gray-900'>
                  {defect.assignee?.full_name || 'Не назначен'}
                </p>
              </div>

              <div>
                <label className='block text-sm font-medium text-gray-700 mb-1'>
                  Срок выполнения
                </label>
                <p className='text-gray-900'>
                  {defect.deadline
                    ? new Date(defect.deadline).toLocaleDateString('ru-RU')
                    : 'Не установлен'}
                </p>
              </div>

              {defect.deadline &&
                new Date(defect.deadline) < new Date() &&
                defect.status !== 'closed' && (
                  <div className='bg-red-50 border border-red-200 rounded-lg p-3'>
                    <p className='text-red-800 text-sm font-medium'>
                      ⚠️ Просрочено
                    </p>
                  </div>
                )}
            </div>
          </div>

          {/* Действия */}
          <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
            <h3 className='text-lg font-semibold text-gray-900 mb-4'>
              Действия
            </h3>
            <div className='space-y-2'>
              {canEditDefect && (
                <Link
                  to={`/defects/${defect.id}/edit`}
                  className='block w-full'
                >
                  <Button variant='secondary' className='w-full justify-center'>
                    Редактировать дефект
                  </Button>
                </Link>
              )}

              {user?.role_name === 'manager' && (
                <Button variant='secondary' className='w-full justify-center'>
                  Назначить исполнителя
                </Button>
              )}

              <Link
                to={`/projects/${defect.project_id}`}
                className='block w-full'
              >
                <Button variant='secondary' className='w-full justify-center'>
                  Перейти к проекту
                </Button>
              </Link>
            </div>
          </div>

          {/* История изменений (можно добавить позже) */}
          <div className='bg-white p-6 rounded-lg shadow-sm border border-gray-200'>
            <h3 className='text-lg font-semibold text-gray-900 mb-4'>
              История
            </h3>
            <p className='text-gray-500 text-sm'>
              История изменений будет доступна в следующих версиях
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
