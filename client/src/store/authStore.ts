import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User, AuthState } from '../types'
import { STORAGE_KEYS } from '../constants'

interface AuthStore extends AuthState {
  login: (token: string, user: User) => void
  logout: () => void
  setUser: (user: User) => void
  setLoading: (isLoading: boolean) => void
  initializeAuth: () => void
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,

      login: (token: string, user: User) => {
        localStorage.setItem(STORAGE_KEYS.AUTH_TOKEN, token)
        set({
          user,
          token,
          isAuthenticated: true,
          isLoading: false,
        })
      },

      logout: () => {
        localStorage.removeItem(STORAGE_KEYS.AUTH_TOKEN)
        set({
          user: null,
          token: null,
          isAuthenticated: false,
          isLoading: false,
        })
      },

      setUser: (user: User) => {
        set({ user })
      },

      setLoading: (isLoading: boolean) => {
        set({ isLoading })
      },

      initializeAuth: () => {
        const token = localStorage.getItem(STORAGE_KEYS.AUTH_TOKEN)
        if (token) {
          set({
            token,
            isAuthenticated: true,
            isLoading: false,
          })
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
