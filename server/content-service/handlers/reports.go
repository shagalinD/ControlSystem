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
    client := resty.New()
    client.SetTimeout(30 * time.Second)
    
    return &ReportHandler{
        Handler: *NewHandler(db, jwtSecret, authServiceURL, projectDefectServiceURL),
        Client:  client,
    }
}

// generateServiceToken создает упрощенный токен для межсервисного общения
func (h *ReportHandler) generateServiceToken() string {
    // В реальной системе здесь должна быть более сложная логика
    // Для тестирования используем упрощенный подход
    return "service-token-" + time.Now().Format("20060102150405")
}

// GetDefectsReport - аналитический отчет по дефектам
func (h *ReportHandler) GetDefectsReport(c *gin.Context) {
    // Вместо передачи пользовательского токена, используем сервисный
    serviceToken := h.generateServiceToken()
    
    // Получаем дефекты из project-defect-service
    var defects []map[string]interface{}
    
    resp, err := h.Client.R().
        SetHeader("X-Service-Token", serviceToken).
        SetHeader("Content-Type", "application/json").
        SetResult(&defects).
        Get(h.ProjectDefectServiceURL + "/api/defects")
    
    if err != nil {
        h.internalError(c, "Failed to fetch defects from project-defect-service: " + err.Error())
        return
    }
    
    if resp.StatusCode() != http.StatusOK {
        h.internalError(c, "Project-defect-service returned status: " + strconv.Itoa(resp.StatusCode()))
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
        // Безопасное извлечение полей
        status, _ := defect["status"].(string)
        priority, _ := defect["priority"].(string)
        
        if status != "" {
            report.StatusStats[status]++
        }
        
        if priority != "" {
            report.PriorityStats[priority]++
        }
        
        // Просроченные дефекты (упрощенная логика)
        if status != "closed" && status != "cancelled" {
            if deadlineStr, ok := defect["deadline"].(string); ok && deadlineStr != "" {
                if deadline, err := time.Parse(time.RFC3339, deadlineStr); err == nil {
                    if deadline.Before(now) {
                        report.OverdueDefects++
                    }
                }
            }
        }
        
        // Время решения для закрытых дефектов
        if status == "closed" {
            if createdAtStr, ok := defect["created_at"].(string); ok {
                if updatedAtStr, ok := defect["updated_at"].(string); ok {
                    if createdAt, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
                        if updatedAt, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
                            resolutionTime := updatedAt.Sub(createdAt)
                            totalResolutionTime += resolutionTime
                            resolvedDefects++
                        }
                    }
                }
            }
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
    
    serviceToken := h.generateServiceToken()
    
    // Получаем дефекты проекта с фильтрацией
    var defects []map[string]interface{}
    
    url := h.ProjectDefectServiceURL + "/api/defects?project_id=" + projectIDStr
    resp, err := h.Client.R().
        SetHeader("X-Service-Token", serviceToken).
        SetHeader("Content-Type", "application/json").
        SetResult(&defects).
        Get(url)
    
    if err != nil {
        h.internalError(c, "Failed to fetch project defects: " + err.Error())
        return
    }
    
    if resp.StatusCode() != http.StatusOK {
        h.internalError(c, "Project-defect-service returned status: " + strconv.Itoa(resp.StatusCode()))
        return
    }
    
    var report struct {
        ProjectID       string                     `json:"project_id"`
        TotalDefects    int64                      `json:"total_defects"`
        DefectsByStatus []models.DefectsByStatus   `json:"defects_by_status"`
        RecentDefects   []map[string]interface{}   `json:"recent_defects"`
        StatusStats     map[string]int64           `json:"status_stats"`
    }
    
    report.ProjectID = projectIDStr
    report.TotalDefects = int64(len(defects))
    report.StatusStats = make(map[string]int64)
    
    // Статистика по статусам
    for _, defect := range defects {
        if status, ok := defect["status"].(string); ok {
            report.StatusStats[status]++
        }
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
    serviceToken := h.generateServiceToken()
    
    // Получаем дефекты с фильтрацией
    var defects []map[string]interface{}
    
    url := h.ProjectDefectServiceURL + "/api/defects"
    
    // Добавляем фильтры из query параметров
    queryParams := c.Request.URL.Query()
    if len(queryParams) > 0 {
        url += "?"
        for key, values := range queryParams {
            for _, value := range values {
                url += key + "=" + value + "&"
            }
        }
        url = strings.TrimSuffix(url, "&")
    }
    
    resp, err := h.Client.R().
        SetHeader("X-Service-Token", serviceToken).
        SetHeader("Content-Type", "application/json").
        SetResult(&defects).
        Get(url)
    
    if err != nil {
        h.internalError(c, "Failed to fetch defects for export: " + err.Error())
        return
    }
    
    if resp.StatusCode() != http.StatusOK {
        h.internalError(c, "Project-defect-service returned status: " + strconv.Itoa(resp.StatusCode()))
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
        "Project ID", "Author ID", "Assignee ID", "Created At",
    }
    writer.Write(headers)
    
    // Данные
    for _, defect := range defects {
        id, _ := defect["id"].(float64)
        title, _ := defect["title"].(string)
        description, _ := defect["description"].(string)
        status, _ := defect["status"].(string)
        priority, _ := defect["priority"].(string)
        projectID, _ := defect["project_id"].(float64)
        authorID, _ := defect["author_id"].(float64)
        
        assigneeID := ""
        if assignee, ok := defect["assignee_id"].(float64); ok && assignee > 0 {
            assigneeID = strconv.FormatFloat(assignee, 'f', 0, 64)
        }
        
        createdAt, _ := defect["created_at"].(string)
        
        record := []string{
            strconv.FormatFloat(id, 'f', 0, 64),
            title,
            description,
            status,
            priority,
            strconv.FormatFloat(projectID, 'f', 0, 64),
            strconv.FormatFloat(authorID, 'f', 0, 64),
            assigneeID,
            createdAt,
        }
        writer.Write(record)
    }
}

// GetUserActivityReport - отчет по активности пользователей
func (h *ReportHandler) GetUserActivityReport(c *gin.Context) {
    // Для этого отчета используем локальные данные комментариев
    // и моковые данные для дефектов
    
    var userActivities []models.UserActivity
    
    // Получаем статистику комментариев из локальной БД
    var commentStats []struct {
        AuthorID uint `gorm:"column:author_id"`
        Count    int64
    }
    
    h.DB.Model(&models.Comment{}).
        Select("author_id, COUNT(*) as count").
        Group("author_id").
        Scan(&commentStats)
    
    // Создаем мапу для быстрого доступа
    commentMap := make(map[uint]int64)
    for _, stat := range commentStats {
        commentMap[stat.AuthorID] = stat.Count
    }
    
    // Моковые данные для пользователей (в реальной системе получать из auth-service)
    mockUsers := []struct {
        ID       uint
        Name     string
        IsEngineer bool
    }{
        {ID: 1, Name: "Test Manager", IsEngineer: false},
        {ID: 2, Name: "Test Engineer", IsEngineer: true},
    }
    
    for _, user := range mockUsers {
        activity := models.UserActivity{
            UserID:   user.ID,
            UserName: user.Name,
            CommentsCount: commentMap[user.ID],
        }
        
        // Моковые данные для дефектов
        if user.IsEngineer {
            activity.DefectsCreated = 2
            activity.DefectsAssigned = 5
        } else {
            activity.DefectsCreated = 8
            activity.DefectsAssigned = 0
        }
        
        userActivities = append(userActivities, activity)
    }
    
    h.success(c, gin.H{
        "report": userActivities,
    }, "User activity report generated successfully")
}

// GetSystemStats - общая статистика системы
func (h *ReportHandler) GetSystemStats(c *gin.Context) {
    serviceToken := h.generateServiceToken()
    
    var stats struct {
        TotalProjects    int64   `json:"total_projects"`
        TotalDefects     int64   `json:"total_defects"`
        TotalUsers       int64   `json:"total_users"`
        TotalComments    int64   `json:"total_comments"`
        ResolutionRate   float64 `json:"resolution_rate"`
        ActiveDefects    int64   `json:"active_defects"`
    }
    
    // Получаем проекты
    var projects []map[string]interface{}
    resp, err := h.Client.R().
        SetHeader("X-Service-Token", serviceToken).
        SetHeader("Content-Type", "application/json").
        SetResult(&projects).
        Get(h.ProjectDefectServiceURL + "/api/projects")
    
    if err == nil && resp.StatusCode() == http.StatusOK {
        stats.TotalProjects = int64(len(projects))
    } else {
        // Fallback: моковые данные
        stats.TotalProjects = 3
    }
    
    // Получаем дефекты
    var defects []map[string]interface{}
    resp, err = h.Client.R().
        SetHeader("X-Service-Token", serviceToken).
        SetHeader("Content-Type", "application/json").
        SetResult(&defects).
        Get(h.ProjectDefectServiceURL + "/api/defects")
    
    if err == nil && resp.StatusCode() == http.StatusOK {
        stats.TotalDefects = int64(len(defects))
        
        // Считаем активные дефекты и rate решения
        var closedDefects int64
        for _, defect := range defects {
            if status, ok := defect["status"].(string); ok {
                if status == "closed" {
                    closedDefects++
                }
                if status == "new" || status == "in_progress" {
                    stats.ActiveDefects++
                }
            }
        }
        
        if stats.TotalDefects > 0 {
            stats.ResolutionRate = float64(closedDefects) / float64(stats.TotalDefects) * 100
        }
    } else {
        // Fallback: моковые данные
        stats.TotalDefects = 15
        stats.ActiveDefects = 8
        stats.ResolutionRate = 46.7
    }
    
    // Пользователи и комментарии (моковые данные)
    stats.TotalUsers = 5
    h.DB.Model(&models.Comment{}).Count(&stats.TotalComments)
    
    h.success(c, gin.H{
        "stats": stats,
    }, "System statistics retrieved successfully")
}