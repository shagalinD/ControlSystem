package handlers

import (
	"encoding/csv"
	"kopatel_online/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportHandler struct {
    Handler
}

func NewReportHandler(db *gorm.DB) *ReportHandler {
    return &ReportHandler{
        Handler: *NewHandler(db), // Исправляем здесь
    }
}

// GetDefectsReport - аналитический отчет по дефектам
func (h *ReportHandler) GetDefectsReport(c *gin.Context) {
    var report struct {
        TotalDefects     int64 `json:"total_defects"`
        DefectsByStatus  []struct {
            Status string `json:"status"`
            Count  int64  `json:"count"`
        } `json:"defects_by_status"`
        DefectsByPriority []struct {
            Priority string `json:"priority"`
            Count    int64  `json:"count"`
        } `json:"defects_by_priority"`
        OverdueDefects int64 `json:"overdue_defects"`
        AvgResolutionTime float64 `json:"avg_resolution_time"`
    }

    // Общее количество дефектов
    h.DB.Model(&models.Defect{}).Count(&report.TotalDefects)

    // Дефекты по статусам
    h.DB.Model(&models.Defect{}).
        Select("status, count(*) as count").
        Group("status").
        Scan(&report.DefectsByStatus)

    // Дефекты по приоритету
    h.DB.Model(&models.Defect{}).
        Select("priority, count(*) as count").
        Group("priority").
        Scan(&report.DefectsByPriority)

    // Просроченные дефекты
    now := time.Now()
    h.DB.Model(&models.Defect{}).
        Where("deadline < ? AND status NOT IN (?)", now, []string{"closed", "cancelled"}).
        Count(&report.OverdueDefects)

    // Среднее время решения (для закрытых дефектов)
    var avgTime struct {
        AvgHours float64 `json:"avg_hours"`
    }
    h.DB.Model(&models.Defect{}).
        Select("AVG(EXTRACT(EPOCH FROM (updated_at - created_at))/3600) as avg_hours").
        Where("status = ?", "closed").
        Scan(&avgTime)
    report.AvgResolutionTime = avgTime.AvgHours

    h.success(c, gin.H{
        "report": report,
    }, "Defects report generated successfully")
}

// GetProjectReport - отчет по конкретному проекту
func (h *ReportHandler) GetProjectReport(c *gin.Context) {
    projectID := c.Param("project_id")

    var project models.Project
    if err := h.DB.First(&project, projectID).Error; err != nil {
        h.notFound(c, "Project not found")
        return
    }

    var report struct {
        Project          models.Project `json:"project"`
        TotalDefects     int64          `json:"total_defects"`
        DefectsByStatus  []struct {
            Status string `json:"status"`
            Count  int64  `json:"count"`
        } `json:"defects_by_status"`
        RecentDefects []models.Defect `json:"recent_defects"`
    }

    report.Project = project

    // Дефекты проекта
    h.DB.Model(&models.Defect{}).Where("project_id = ?", projectID).Count(&report.TotalDefects)

    // Дефекты по статусам
    h.DB.Model(&models.Defect{}).
        Select("status, count(*) as count").
        Where("project_id = ?", projectID).
        Group("status").
        Scan(&report.DefectsByStatus)

    // Последние дефекты
    h.DB.Preload("Author").Preload("Assignee").
        Where("project_id = ?", projectID).
        Order("created_at DESC").
        Limit(10).
        Find(&report.RecentDefects)

    h.success(c, gin.H{
        "report": report,
    }, "Project report generated successfully")
}

// ExportDefectsCSV - экспорт дефектов в CSV
func (h *ReportHandler) ExportDefectsCSV(c *gin.Context) {
    var defects []models.Defect
    
    query := h.DB.Preload("Project").Preload("Author").Preload("Assignee")
    
    // Фильтрация
    if projectID := c.Query("project_id"); projectID != "" {
        query = query.Where("project_id = ?", projectID)
    }
    if status := c.Query("status"); status != "" {
        query = query.Where("status = ?", status)
    }

    if err := query.Find(&defects).Error; err != nil {
        h.internalError(c, "Failed to fetch defects for export")
        return
    }

    // Создаем CSV
    c.Writer.Header().Set("Content-Type", "text/csv")
    c.Writer.Header().Set("Content-Disposition", "attachment;filename=defects.csv")
    
    writer := csv.NewWriter(c.Writer)
    defer writer.Flush()

    // Заголовки CSV
    headers := []string{
        "ID", "Title", "Description", "Status", "Priority", 
        "Project", "Author", "Assignee", "Deadline", "Created At",
    }
    writer.Write(headers)

    // Данные
    for _, defect := range defects {
        assignee := ""
        if defect.Assignee != nil {
            assignee = defect.Assignee.FullName
        }
        
        deadline := ""
        if defect.Deadline != nil {
            deadline = defect.Deadline.Format("2006-01-02")
        }

        record := []string{
            strconv.FormatUint(uint64(defect.ID), 10),
            defect.Title,
            defect.Description,
            string(defect.Status),
            string(defect.Priority),
            defect.Project.Name,
            defect.Author.FullName,
            assignee,
            deadline,
            defect.CreatedAt.Format("2006-01-02 15:04:05"),
        }
        writer.Write(record)
    }
}

// GetUserActivityReport - отчет по активности пользователей
func (h *ReportHandler) GetUserActivityReport(c *gin.Context) {
    var report []struct {
        UserID    uint   `json:"user_id"`
        UserName  string `json:"user_name"`
        DefectsCreated int64 `json:"defects_created"`
        DefectsAssigned int64 `json:"defects_assigned"`
        CommentsCount  int64 `json:"comments_count"`
    }

    h.DB.Table("users u").
        Select(`
            u.id as user_id,
            u.full_name as user_name,
            (SELECT COUNT(*) FROM defects d WHERE d.author_id = u.id) as defects_created,
            (SELECT COUNT(*) FROM defects d WHERE d.assignee_id = u.id) as defects_assigned,
            (SELECT COUNT(*) FROM comments c WHERE c.author_id = u.id) as comments_count
        `).
        Where("u.deleted_at IS NULL").
        Scan(&report)

    h.success(c, gin.H{
        "report": report,
    }, "User activity report generated successfully")
}