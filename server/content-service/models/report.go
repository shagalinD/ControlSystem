package models

type DefectStatus string
type DefectPriority string

const (
	StatusNew        DefectStatus = "new"
	StatusInProgress DefectStatus = "in_progress"
	StatusOnReview   DefectStatus = "on_review"
	StatusClosed     DefectStatus = "closed"
	StatusCancelled  DefectStatus = "cancelled"

	PriorityLow      DefectPriority = "low"
	PriorityMedium   DefectPriority = "medium"
	PriorityHigh     DefectPriority = "high"
	PriorityCritical DefectPriority = "critical"
)

// Структуры для отчетов
type DefectsByStatus struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type DefectsByPriority struct {
	Priority string `json:"priority"`
	Count    int64  `json:"count"`
}

type UserActivity struct {
	UserID          uint   `json:"user_id"`
	UserName        string `json:"user_name"`
	DefectsCreated  int64  `json:"defects_created"`
	DefectsAssigned int64  `json:"defects_assigned"`
	CommentsCount   int64  `json:"comments_count"`
}