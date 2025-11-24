package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BaseModel struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type Role struct {
    BaseModel
    RoleName string `gorm:"uniqueIndex;not null" json:"role_name"`
    Users    []User `json:"users,omitempty"`
}

type User struct {
    BaseModel
    Email        string `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash string `gorm:"not null" json:"-"`
    FullName     string `gorm:"not null" json:"full_name"`
    RoleID       uint   `gorm:"not null" json:"role_id"`
    Role         Role   `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

type UserCreateRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    FullName string `json:"full_name" binding:"required"`
    RoleID   uint   `json:"role_id" binding:"required"`
}

type UserUpdateRequest struct {
    Email    string `json:"email" binding:"required,email"`
    FullName string `json:"full_name" binding:"required"`
}

type UserLoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
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