package models

type Attachment struct {
	BaseModel
	DefectID   uint   `gorm:"not null" json:"defect_id"`
	Filename   string `gorm:"not null" json:"filename"`
	Filepath   string `gorm:"not null" json:"filepath"`
	FileSize   int64  `json:"file_size"`
	MimeType   string `json:"mime_type"`
	UploadedBy uint   `gorm:"not null" json:"uploaded_by"`
}