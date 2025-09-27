package models

type Project struct {
	BaseModel
	Name        string   `gorm:"not null" json:"name"`
	Description string   `json:"description"`
	ManagerID   uint     `gorm:"not null" json:"manager_id"`
	Manager     User     `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	Defects     []Defect `json:"defects,omitempty"`
}

type ProjectCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	ManagerID   uint   `json:"manager_id" validate:"required"`
}