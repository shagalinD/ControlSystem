import React, { useEffect } from 'react'
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import { ProtectedRoute } from './components/ProtectedRoute'
import { Layout } from './components/layout/Layout'
import { LoginPage } from './pages/LoginPage'
import { HomePage } from './pages/HomePage'
import { DefectsPage } from './pages/DefectsPage'
import { LoadingSpinner } from './components/ui/LoadingSpinner'
import { CreateDefectPage } from './pages/CreateDefectPage'
import { ProjectsPage } from './pages/ProjectPage'
import { CreateProjectPage } from './pages/CreateProjectPage'
import { ProjectDetailsPage } from './pages/ProjectDetailsPage'
import { DefectDetailsPage } from './pages/DefectDetailsPage'
import { EditDefectPage } from './pages/EditDefectPage'
import { EditProjectPage } from './pages/EditProjectPage'
import { ReportsPage } from './pages/ReportsPage'
import { ProfilePage } from './pages/ProfilePage'
import { RegisterPage } from './pages/RegisterPage'

const App: React.FC = () => {
  const { isAuthenticated, initializeAuth, isLoading } = useAuthStore()

  useEffect(() => {
    initializeAuth()
  }, [initializeAuth])

  if (isLoading) {
    return (
      <div className='min-h-screen flex items-center justify-center'>
        <LoadingSpinner size='lg' />
      </div>
    )
  }

  return (
    <Router>
      <Routes>
        <Route
          path='/login'
          element={
            !isAuthenticated ? <LoginPage /> : <Navigate to='/' replace />
          }
        />
        <Route
          path='/register'
          element={
            !isAuthenticated ? <RegisterPage /> : <Navigate to='/' replace />
          }
        />

        <Route
          path='/'
          element={
            <ProtectedRoute>
              <Layout />
            </ProtectedRoute>
          }
        >
          <Route index element={<HomePage />} />
          <Route path='defects' element={<DefectsPage />} />
          <Route path='defects/create' element={<CreateDefectPage />} />
          <Route path='defects/:id' element={<DefectDetailsPage />} />
          <Route path='defects/:id/edit' element={<EditDefectPage />} />
          <Route path='projects' element={<ProjectsPage />} />
          <Route path='projects/create' element={<CreateProjectPage />} />
          <Route path='projects/:id' element={<ProjectDetailsPage />} />
          <Route path='projects/:id/edit' element={<EditProjectPage />} />
          <Route path='reports' element={<ReportsPage />} />
          <Route path='profile' element={<ProfilePage />} />
          <Route path='*' element={<Navigate to='/' replace />} />
        </Route>
      </Routes>
    </Router>
  )
}

export default App
