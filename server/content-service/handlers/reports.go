package handlers

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"
	"time"

	"content-service/models"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type ReportHandler struct {
    Handler
    Client *resty.Client
}

func NewReportHandler(db *gorm.DB, jwtSecret, authServiceURL, projectDefectServiceURL string) *ReportHandler {
    return &ReportHandler{
        Handler: *NewHandler(db, jwtSecret, authServiceURL, projectDefectServiceURL),
        Client:  resty.New(),
    }
}

// GetDefectsReport - аналитический отчет по дефектам
func (h *ReportHandler) GetDefectsReport(c *gin.Context) {
    // Получаем токен из заголовков для межсервисного вызова
    authHeader := c.GetHeader("Authorization")
    
    // Получаем дефекты из project-defect-service
    var defects []struct {
        ID          uint                  `json:"id"`
        Title       string                `json:"title"`
        Status      models.DefectStatus   `json:"status"`
        Priority    models.DefectPriority `json:"priority"`
        CreatedAt   time.Time             `json:"created_at"`
        UpdatedAt   time.Time             `json:"updated_at"`
        Deadline    *time.Time            `json:"deadline"`
        ProjectID   uint                  `json:"project_id"`
    }
    
    resp, err := h.Client.R().
        SetHeader("Authorization", authHeader).
        SetResult(&defects).
        Get(h.ProjectDefectServiceURL + "/api/defects")
    
    if err != nil || resp.StatusCode() != http.StatusOK {
        h.internalError(c, "Failed to fetch defects from project-defect-service")
        return
    }
    
    var report struct {
        TotalDefects      int64                         `json:"total_defects"`
        DefectsByStatus   []models.DefectsByStatus      `json:"defects_by_status"`
        DefectsByPriority []models.DefectsByPriority    `json:"defects_by_priority"`
        OverdueDefects    int64                         `json:"overdue_defects"`
        AvgResolutionTime float64                       `json:"avg_resolution_time"`
        StatusStats       map[string]int64              `json:"status_stats"`
        PriorityStats     map[string]int64              `json:"priority_stats"`
    }
    
    report.TotalDefects = int64(len(defects))
    report.StatusStats = make(map[string]int64)
    report.PriorityStats = make(map[string]int64)
    
    now := time.Now()
    var totalResolutionTime time.Duration
    var resolvedDefects int
    
    // Анализируем дефекты
    for _, defect := range defects {
        // Статистика по статусам
        report.StatusStats[string(defect.Status)]++
        
        // Статистика по приоритетам
        report.PriorityStats[string(defect.Priority)]++
        
        // Просроченные дефекты
        if defect.Deadline != nil && defect.Deadline.Before(now) && 
           defect.Status != models.StatusClosed && defect.Status != models.StatusCancelled {
            report.OverdueDefects++
        }
        
        // Время решения для закрытых дефектов
        if defect.Status == models.StatusClosed {
            resolutionTime := defect.UpdatedAt.Sub(defect.CreatedAt)
            totalResolutionTime += resolutionTime
            resolvedDefects++
        }
    }
    
    // Преобразуем мапы в слайсы для ответа
    for status, count := range report.StatusStats {
        report.DefectsByStatus = append(report.DefectsByStatus, models.DefectsByStatus{
            Status: status,
            Count:  count,
        })
    }
    
    for priority, count := range report.PriorityStats {
        report.DefectsByPriority = append(report.DefectsByPriority, models.DefectsByPriority{
            Priority: priority,
            Count:    count,
        })
    }
    
    // Среднее время решения
    if resolvedDefects > 0 {
        report.AvgResolutionTime = totalResolutionTime.Hours() / float64(resolvedDefects)
    }
    
    h.success(c, gin.H{
        "report": report,
    }, "Defects report generated successfully")
}

// GetProjectReport - отчет по конкретному проекту
func (h *ReportHandler) GetProjectReport(c *gin.Context) {
    projectIDStr := c.Param("project_id")
    
    authHeader := c.GetHeader("Authorization")
    
    // Получаем проект
    var project struct {
        ID          uint   `json:"id"`
        Name        string `json:"name"`
        Description string `json:"description"`
    }
    
    resp, err := h.Client.R().
        SetHeader("Authorization", authHeader).
        SetResult(&project).
        Get(h.ProjectDefectServiceURL + "/api/projects/" + projectIDStr)
    
    if err != nil || resp.StatusCode() != http.StatusOK {
        h.notFound(c, "Project not found")
        return
    }
    
    // Получаем дефекты проекта
    var defects []struct {
        ID          uint                  `json:"id"`
        Title       string                `json:"title"`
        Status      models.DefectStatus   `json:"status"`
        Priority    models.DefectPriority `json:"priority"`
        CreatedAt   time.Time             `json:"created_at"`
    }
    
    resp, err = h.Client.R().
        SetHeader("Authorization", authHeader).
        SetResult(&defects).
        Get(h.ProjectDefectServiceURL + "/api/defects?project_id=" + projectIDStr)
    
    if err != nil {
        h.internalError(c, "Failed to fetch project defects")
        return
    }
    
    var report struct {
        Project         interface{}              `json:"project"`
        TotalDefects    int64                    `json:"total_defects"`
        DefectsByStatus []models.DefectsByStatus `json:"defects_by_status"`
        RecentDefects   []interface{}            `json:"recent_defects"`
        StatusStats     map[string]int64         `json:"status_stats"`
    }
    
    report.Project = project
    report.TotalDefects = int64(len(defects))
    report.StatusStats = make(map[string]int64)
    
    // Статистика по статусам
    for _, defect := range defects {
        report.StatusStats[string(defect.Status)]++
    }
    
    // Преобразуем в слайс
    for status, count := range report.StatusStats {
        report.DefectsByStatus = append(report.DefectsByStatus, models.DefectsByStatus{
            Status: status,
            Count:  count,
        })
    }
    
    // Последние 10 дефектов
    recentCount := 10
    if len(defects) < recentCount {
        recentCount = len(defects)
    }
    
    for i := 0; i < recentCount; i++ {
        report.RecentDefects = append(report.RecentDefects, defects[i])
    }
    
    h.success(c, gin.H{
        "report": report,
    }, "Project report generated successfully")
}

// ExportDefectsCSV - экспорт дефектов в CSV
func (h *ReportHandler) ExportDefectsCSV(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    
    // Получаем дефекты с фильтрацией
    var defects []struct {
        ID          uint                  `json:"id"`
        Title       string                `json:"title"`
        Description string                `json:"description"`
        Status      models.DefectStatus   `json:"status"`
        Priority    models.DefectPriority `json:"priority"`
        CreatedAt   time.Time             `json:"created_at"`
        ProjectID   uint                  `json:"project_id"`
        AuthorID    uint                  `json:"author_id"`
        AssigneeID  *uint                 `json:"assignee_id,omitempty"`
        Deadline    *time.Time            `json:"deadline,omitempty"`
    }
    
    url := h.ProjectDefectServiceURL + "/api/defects"
    
    // Добавляем фильтры из query параметров
    if projectID := c.Query("project_id"); projectID != "" {
        url += "?project_id=" + projectID
    }
    if status := c.Query("status"); status != "" {
        if strings.Contains(url, "?") {
            url += "&status=" + status
        } else {
            url += "?status=" + status
        }
    }
    
    resp, err := h.Client.R().
        SetHeader("Authorization", authHeader).
        SetResult(&defects).
        Get(url)
    
    if err != nil || resp.StatusCode() != http.StatusOK {
        h.internalError(c, "Failed to fetch defects for export")
        return
    }
    
    // Создаем CSV
    c.Writer.Header().Set("Content-Type", "text/csv")
    c.Writer.Header().Set("Content-Disposition", "attachment;filename=defects_export.csv")
    
    writer := csv.NewWriter(c.Writer)
    defer writer.Flush()
    
    // Заголовки CSV
    headers := []string{
        "ID", "Title", "Description", "Status", "Priority", 
        "Project ID", "Author ID", "Assignee ID", "Deadline", "Created At",
    }
    writer.Write(headers)
    
    // Данные
    for _, defect := range defects {
        assignee := ""
        if defect.AssigneeID != nil {
            assignee = strconv.FormatUint(uint64(*defect.AssigneeID), 10)
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
            strconv.FormatUint(uint64(defect.ProjectID), 10),
            strconv.FormatUint(uint64(defect.AuthorID), 10),
            assignee,
            deadline,
            defect.CreatedAt.Format("2006-01-02 15:04:05"),
        }
        writer.Write(record)
    }
}

// GetUserActivityReport - отчет по активности пользователей
func (h *ReportHandler) GetUserActivityReport(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    
    // Получаем пользователей из auth-service
    var users []struct {
        ID       uint   `json:"id"`
        Email    string `json:"email"`
        FullName string `json:"full_name"`
    }
    
    resp, err := h.Client.R().
        SetHeader("Authorization", authHeader).
        SetResult(&users).
        Get(h.AuthServiceURL + "/api/users")
    
    if err != nil || resp.StatusCode() != http.StatusOK {
        h.internalError(c, "Failed to fetch users from auth-service")
        return
    }
    
    var report []models.UserActivity
    
    // Для каждого пользователя собираем статистику
    for _, user := range users {
        activity := models.UserActivity{
            UserID:   user.ID,
            UserName: user.FullName,
        }
        
        // Получаем дефекты созданные пользователем
        var createdDefects []interface{}
        resp, err = h.Client.R().
            SetHeader("Authorization", authHeader).
            SetResult(&createdDefects).
            Get(h.ProjectDefectServiceURL + "/api/defects?author_id=" + strconv.FormatUint(uint64(user.ID), 10))
        
        if err == nil && resp.StatusCode() == http.StatusOK {
            activity.DefectsCreated = int64(len(createdDefects))
        }
        
        // Получаем дефекты назначенные пользователю
        var assignedDefects []interface{}
        resp, err = h.Client.R().
            SetHeader("Authorization", authHeader).
            SetResult(&assignedDefects).
            Get(h.ProjectDefectServiceURL + "/api/defects?assignee_id=" + strconv.FormatUint(uint64(user.ID), 10))
        
        if err == nil && resp.StatusCode() == http.StatusOK {
            activity.DefectsAssigned = int64(len(assignedDefects))
        }
        
        // Комментарии пользователя (локальные данные)
        var commentsCount int64
        h.DB.Model(&models.Comment{}).
            Where("author_id = ?", user.ID).
            Count(&commentsCount)
        activity.CommentsCount = commentsCount
        
        report = append(report, activity)
    }
    
    h.success(c, gin.H{
        "report": report,
    }, "User activity report generated successfully")
}

// GetSystemStats - общая статистика системы
func (h *ReportHandler) GetSystemStats(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    
    var stats struct {
        TotalProjects    int64   `json:"total_projects"`
        TotalDefects     int64   `json:"total_defects"`
        TotalUsers       int64   `json:"total_users"`
        TotalComments    int64   `json:"total_comments"`
        ResolutionRate   float64 `json:"resolution_rate"`
        ActiveDefects    int64   `json:"active_defects"`
    }
    
    // Получаем проекты
    var projects []interface{}
    resp, err := h.Client.R().
        SetHeader("Authorization", authHeader).
        SetResult(&projects).
        Get(h.ProjectDefectServiceURL + "/api/projects")
    
    if err == nil && resp.StatusCode() == http.StatusOK {
        stats.TotalProjects = int64(len(projects))
    }
    
    // Получаем дефекты для анализа
    var defects []struct {
        Status models.DefectStatus `json:"status"`
    }
    
    resp, err = h.Client.R().
        SetHeader("Authorization", authHeader).
        SetResult(&defects).
        Get(h.ProjectDefectServiceURL + "/api/defects")
    
    if err == nil && resp.StatusCode() == http.StatusOK {
        stats.TotalDefects = int64(len(defects))
        
        // Считаем активные дефекты и rate решения
        var closedDefects int64
        for _, defect := range defects {
            if defect.Status == models.StatusClosed {
                closedDefects++
            }
            if defect.Status == models.StatusNew || defect.Status == models.StatusInProgress {
                stats.ActiveDefects++
            }
        }
        
        if stats.TotalDefects > 0 {
            stats.ResolutionRate = float64(closedDefects) / float64(stats.TotalDefects) * 100
        }
    }
    
    // Получаем пользователей
    var users []interface{}
    resp, err = h.Client.R().
        SetHeader("Authorization", authHeader).
        SetResult(&users).
        Get(h.AuthServiceURL + "/api/users")
    
    if err == nil && resp.StatusCode() == http.StatusOK {
        stats.TotalUsers = int64(len(users))
    }
    
    // Комментарии (локальные)
    h.DB.Model(&models.Comment{}).Count(&stats.TotalComments)
    
    h.success(c, gin.H{
        "stats": stats,
    }, "System statistics retrieved successfully")
}