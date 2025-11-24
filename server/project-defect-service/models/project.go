package models

type Project struct {
	BaseModel
	Name        string   `gorm:"not null" json:"name"`
	Description string   `json:"description"`
	ManagerID   uint     `gorm:"not null" json:"manager_id"`
	Defects     []Defect `json:"defects,omitempty"`
}

type ProjectCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ManagerID   uint   `json:"manager_id" binding:"required"`
}