package models

type Role struct {
	BaseModel
	RoleName string `gorm:"uniqueIndex;not null" json:"role_name"`
	Users    []User `json:"users,omitempty"`
}