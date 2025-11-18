package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type DefectStatus string
type DefectPriority string

const (
    StatusNew       DefectStatus = "new"
    StatusInProgress DefectStatus = "in_progress"
    StatusOnReview  DefectStatus = "on_review"
    StatusClosed    DefectStatus = "closed"
    StatusCancelled DefectStatus = "cancelled"
    
    PriorityLow    DefectPriority = "low"
    PriorityMedium DefectPriority = "medium"
    PriorityHigh   DefectPriority = "high"
    PriorityCritical DefectPriority = "critical"
)



func (s DefectStatus) IsValid() bool {
    switch s {
    case StatusNew, StatusInProgress, StatusOnReview, StatusClosed, StatusCancelled:
        return true
    default:
        return false
    }
}

func (p DefectPriority) IsValid() bool {
    switch p {
    case PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
        return true
    default:
        return false
    }
}

// Добавим модель для вложений
type Attachment struct {
    BaseModel
    DefectID   uint   `gorm:"not null" json:"defect_id"`
    Defect     Defect `gorm:"foreignKey:DefectID" json:"defect,omitempty"`
    Filename   string `gorm:"not null" json:"filename"`
    Filepath   string `gorm:"not null" json:"filepath"`
    FileSize   int64  `json:"file_size"`
    MimeType   string `json:"mime_type"`
    UploadedBy uint   `gorm:"not null" json:"uploaded_by"`
    Uploader   User   `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
}

// Добавим историю изменений дефектов
type DefectHistory struct {
    BaseModel
    DefectID  uint      `gorm:"not null" json:"defect_id"`
    Defect    Defect    `gorm:"foreignKey:DefectID" json:"defect,omitempty"`
    Field     string    `gorm:"not null" json:"field"` // "status", "priority", "assignee", "title", etc.
    OldValue  string    `json:"old_value"`
    NewValue  string    `json:"new_value"`
    ChangedBy uint      `gorm:"not null" json:"changed_by"`
    Changer   User      `gorm:"foreignKey:ChangedBy" json:"changer,omitempty"`
}

// Обновим модель Defect - добавим связь с вложениями
type Defect struct {
    BaseModel
    Title       string         `gorm:"not null" json:"title"`
    Description string         `json:"description"`
    Status      DefectStatus   `gorm:"not null;default:'new'" json:"status"`
    Priority    DefectPriority `gorm:"not null;default:'medium'" json:"priority"`
    Deadline    *Date     `json:"deadline,omitempty"`
    
    ProjectID   uint    `gorm:"not null" json:"project_id"`
    Project     Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
    
    AuthorID    uint    `gorm:"not null" json:"author_id"`
    Author      User    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
    
    AssigneeID  *uint   `json:"assignee_id,omitempty"`
    Assignee    *User   `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
    
    Comments    []Comment     `json:"comments,omitempty"`
    Attachments []Attachment  `json:"attachments,omitempty"`
    History     []DefectHistory `json:"history,omitempty"`
}

type Date struct {
    time.Time
}

// Формат для JSON - УПРОЩЕННАЯ ВЕРСИЯ
func (d *Date) UnmarshalJSON(b []byte) error {
    str := string(b)
    
    // Обрабатываем null и пустые строки
    if str == "null" || str == `""` || str == "" {
        d.Time = time.Time{}
        return nil
    }
    
    // Убираем кавычки если есть
    str = strings.Trim(str, `"`)
    
    // Парсим дату в формате YYYY-MM-DD
    t, err := time.Parse("2006-01-02", str)
    if err != nil {
        return fmt.Errorf("failed to parse date %s: %v", str, err)
    }
    d.Time = t
    return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
    if d.IsZero() {
        return []byte("null"), nil
    }
    return []byte(fmt.Sprintf(`"%s"`, d.Format("2006-01-02"))), nil
}

// Формат для базы данных
func (d Date) Value() (driver.Value, error) {
    if d.IsZero() {
        return nil, nil
    }
    return d.Format("2006-01-02"), nil
}

func (d *Date) Scan(value interface{}) error {
    if value == nil {
        d.Time = time.Time{}
        return nil
    }
    
    switch v := value.(type) {
    case time.Time:
        d.Time = time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location())
    case string:
        t, err := time.Parse("2006-01-02", v)
        if err != nil {
            return err
        }
        d.Time = t
    case []byte:
        t, err := time.Parse("2006-01-02", string(v))
        if err != nil {
            return err
        }
        d.Time = t
    default:
        return fmt.Errorf("unsupported type for Date: %T", value)
    }
    return nil
}

// В структурах Request добавляем теги binding для Gin
type DefectCreateRequest struct {
    Title       string         `json:"title" binding:"required"`
    Description string         `json:"description"`
    Priority    DefectPriority `json:"priority" binding:"required,oneof=low medium high critical"`
    Deadline    *Date          `json:"deadline,omitempty"`
    ProjectID   uint           `json:"project_id" binding:"required"`
    AssigneeID  *uint          `json:"assignee_id,omitempty"`
}

type DefectUpdateRequest struct {
    Title       *string         `json:"title,omitempty"`
    Description *string         `json:"description,omitempty"`
    Status      *DefectStatus   `json:"status,omitempty" binding:"omitempty,oneof=new in_progress on_review closed cancelled"`
    Priority    *DefectPriority `json:"priority,omitempty" binding:"omitempty,oneof=low medium high critical"`
    Deadline    *Date           `json:"deadline,omitempty"`
    AssigneeID  *uint           `json:"assignee_id,omitempty"`
}
