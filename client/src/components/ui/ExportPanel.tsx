import React, { useState } from 'react'
import type { DefectFilters } from '../../types'
import { reportService } from '../../services/reportService'
import { Button } from '../ui/Button'

interface ExportPanelProps {
  filters: DefectFilters
}

export const ExportPanel: React.FC<ExportPanelProps> = ({ filters }) => {
  const [isExporting, setIsExporting] = useState(false)

  const handleExportCSV = async () => {
    setIsExporting(true)

    try {
      const blob = await reportService.exportDefectsToCSV(filters)

      // Создаем ссылку для скачивания
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.style.display = 'none'
      a.href = url
      a.download = `defects-report-${
        new Date().toISOString().split('T')[0]
      }.csv`

      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
    } catch (error) {
      console.error('Export error:', error)
      alert('Ошибка при экспорте данных')
    } finally {
      setIsExporting(false)
    }
  }

  return (
    <div className='flex space-x-3'>
      <Button
        onClick={handleExportCSV}
        isLoading={isExporting}
        variant='secondary'
      >
        📊 Экспорт в CSV
      </Button>
    </div>
  )
}
