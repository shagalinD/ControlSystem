package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	BaseModel
	Email           string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash    string    `gorm:"not null" json:"-"`
	FullName        string    `gorm:"not null" json:"full_name"`
	RoleID          uint      `gorm:"not null" json:"role_id"`
	Role            Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Projects        []Project `gorm:"foreignKey:ManagerID" json:"projects,omitempty"`
	Defects         []Defect  `gorm:"foreignKey:AuthorID" json:"defects,omitempty"`
	AssignedDefects []Defect  `gorm:"foreignKey:AssigneeID" json:"assigned_defects,omitempty"`
	Comments        []Comment `gorm:"foreignKey:AuthorID" json:"comments,omitempty"`
}

type UserCreateRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	RoleID   uint   `json:"role_id" validate:"required"`
}

type UserUpdateRequest struct {
	Email string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u *User) SetPassword(password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedBytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type UserResponse struct {
    ID        uint   `json:"id"`
    Email     string `json:"email"`
    FullName  string `json:"full_name"`
    RoleID    uint   `json:"role_id"`
    RoleName  string `json:"role_name"`
    CreatedAt string `json:"created_at,omitempty"`
    UpdatedAt string `json:"updated_at,omitempty"`
}

// ToResponse - преобразует User в UserResponse
func (u *User) ToResponse() UserResponse {
    roleName := ""
    if u.Role.ID != 0 {
        roleName = u.Role.RoleName
    }
    
    return UserResponse{
        ID:        u.ID,
        Email:     u.Email,
        FullName:  u.FullName,
        RoleID:    u.RoleID,
        RoleName:  roleName,
        CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
        UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z"),
    }
}