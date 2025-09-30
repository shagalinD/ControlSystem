package models

import "time"

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



type DefectCreateRequest struct {
    Title       string         `json:"title" validate:"required"`
    Description string         `json:"description"`
    Priority    DefectPriority `json:"priority" validate:"required,oneof=low medium high critical"`
    Deadline    *time.Time     `json:"deadline,omitempty"`
    ProjectID   uint           `json:"project_id" validate:"required"`
    AssigneeID  *uint          `json:"assignee_id,omitempty"`
}

type DefectUpdateRequest struct {
    Title       *string         `json:"title,omitempty"`
    Description *string         `json:"description,omitempty"`
    Status      *DefectStatus   `json:"status,omitempty" validate:"omitempty,oneof=new in_progress on_review closed cancelled"`
    Priority    *DefectPriority `json:"priority,omitempty" validate:"omitempty,oneof=low medium high critical"`
    Deadline    *time.Time      `json:"deadline,omitempty"`
    AssigneeID  *uint           `json:"assignee_id,omitempty"`
}

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
    Deadline    *time.Time     `json:"deadline,omitempty"`
    
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