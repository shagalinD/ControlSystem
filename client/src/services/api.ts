import axios from 'axios'
import type { ApiResponse } from '../types'
import { API_BASE_URL, STORAGE_KEYS } from '../constants'

export const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor для добавления токена
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem(STORAGE_KEYS.AUTH_TOKEN)
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor для обработки ошибок
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Токен истек или невалидный
      localStorage.removeItem(STORAGE_KEYS.AUTH_TOKEN)
      localStorage.removeItem(STORAGE_KEYS.USER_DATA)
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Вспомогательные функции для обработки ответов
export const handleApiResponse = <T>(response: { data: ApiResponse<T> }): T => {
  if (response.data.success) {
    return response.data.data
  }
  throw new Error(response.data.message || 'API request failed')
}

export const handleApiError = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    return error.response?.data?.error || error.message || 'Network error'
  }
  return error instanceof Error ? error.message : 'Unknown error occurred'
}
