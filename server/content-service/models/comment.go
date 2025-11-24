package models

type Comment struct {
	BaseModel
	Text     string `gorm:"not null" json:"text"`
	DefectID uint   `gorm:"not null" json:"defect_id"`
	AuthorID uint   `gorm:"not null" json:"author_id"`
}

type CommentCreateRequest struct {
	Text     string `json:"text" binding:"required"`
	DefectID uint   `json:"defect_id" binding:"required"`
}