package models

type Comment struct {
	BaseModel
	Text     string `gorm:"not null" json:"text"`
	DefectID uint   `gorm:"not null" json:"defect_id"`
	Defect   Defect `gorm:"foreignKey:DefectID" json:"defect,omitempty"`
	AuthorID uint   `gorm:"not null" json:"author_id"`
	Author   User   `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

type CommentCreateRequest struct {
	Text     string `json:"text" validate:"required"`
	DefectID uint   `json:"defect_id" validate:"required"`
}