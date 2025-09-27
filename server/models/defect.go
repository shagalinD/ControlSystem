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
    
    Comments    []Comment `json:"comments,omitempty"`
}

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