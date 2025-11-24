package handlers

import (
	"content-service/models"
	"encoding/csv"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportHandler struct {
    Handler
}

func NewReportHandler(db *gorm.DB, jwtSecret, authServiceURL, projectDefectServiceURL string) *ReportHandler {
    return &ReportHandler{
        Handler: *NewHandler(db, jwtSecret, authServiceURL, projectDefectServiceURL),
    }
}

func (h *ReportHandler) GetDefectsReport(c *gin.Context) {
    var report struct {
        TotalDefects     int64 `json:"total_defects"`
        DefectsByStatus  []models.DefectsByStatus `json:"defects_by_status"`
        DefectsByPriority []models.DefectsByPriority `json:"defects_by_priority"`
        OverdueDefects   int64 `json:"overdue_defects"`
        AvgResolutionTime float64 `json:"avg_resolution_time"`
    }

    // TODO: Получать данные через project-defect-service API
    // Временная реализация с локальными данными
    
    h.success(c, gin.H{
        "report": report,
    }, "Defects report generated successfully")
}

func (h *ReportHandler) GetProjectReport(c *gin.Context) {
    projectIDStr := c.Param("project_id")

		// Конвертируем project_id из string в uint
    projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
    if err != nil {
        h.badRequest(c, "Invalid project ID")
        return
    }

    // TODO: Получать данные проекта через project-defect-service
    // TODO: Получать дефекты проекта через project-defect-service

    var report struct {
        ProjectID       uint `json:"project_id"`
        TotalDefects    int64  `json:"total_defects"`
        DefectsByStatus []models.DefectsByStatus `json:"defects_by_status"`
    }

    report.ProjectID = uint(projectID)

    h.success(c, gin.H{
        "report": report,
    }, "Project report generated successfully")
}

func (h *ReportHandler) ExportDefectsCSV(c *gin.Context) {
    // TODO: Получать дефекты через project-defect-service API
    
    // Временная реализация
    c.Writer.Header().Set("Content-Type", "text/csv")
    c.Writer.Header().Set("Content-Disposition", "attachment;filename=defects.csv")
    
    writer := csv.NewWriter(c.Writer)
    defer writer.Flush()

    headers := []string{
        "ID", "Title", "Description", "Status", "Priority", 
        "Project ID", "Author ID", "Assignee ID", "Created At",
    }
    writer.Write(headers)

    // TODO: Добавить реальные данные из project-defect-service
}

func (h *ReportHandler) GetUserActivityReport(c *gin.Context) {
    var userActivities []models.UserActivity

    // Аналитика по комментариям (локальные данные)
    h.DB.Raw(`
        SELECT 
            author_id as user_id,
            COUNT(*) as comments_count
        FROM comments 
        WHERE deleted_at IS NULL
        GROUP BY author_id
    `).Scan(&userActivities)

    // TODO: Получить данные по дефектам через project-defect-service
    // и объединить с локальными данными по комментариям

    h.success(c, gin.H{
        "report": userActivities,
    }, "User activity report generated successfully")
}