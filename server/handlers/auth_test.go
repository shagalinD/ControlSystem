package handlers

import (
	"bytes"
	"encoding/json"
	"kopatel_online/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("failed to connect test database")
    }
    
    // Миграции
    db.AutoMigrate(&models.Role{}, &models.User{}, &models.Project{}, &models.Defect{}, &models.Comment{})
    
    // Тестовые данные
    db.Create(&models.Role{RoleName: "engineer"})
    db.Create(&models.Role{RoleName: "manager"})
    db.Create(&models.Role{RoleName: "observer"})
    
    return db
}

func TestAuthHandler_Register(t *testing.T) {
    db := setupTestDB()
    handler := NewAuthHandler(db, "test-secret")
    
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.POST("/register", handler.Register)

    // Тест успешной регистрации
    registerData := map[string]interface{}{
        "email":     "test@example.com",
        "password":  "password123",
        "full_name": "Test User",
        "role_id":   1,
    }
    jsonData, _ := json.Marshal(registerData)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
    
    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
    
    if !response.Success {
        t.Errorf("Expected success true, got false. Error: %s", response.Error)
    }
}

func TestAuthHandler_Login(t *testing.T) {
    db := setupTestDB()
    handler := NewAuthHandler(db, "test-secret")
    
    // Сначала создаем пользователя
    user := models.User{
        Email:    "test@example.com",
        FullName: "Test User",
        RoleID:   1,
    }
    user.SetPassword("password123")
    db.Create(&user)

    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.POST("/login", handler.Login)

    // Тест успешного логина
    loginData := map[string]string{
        "email":    "test@example.com",
        "password": "password123",
    }
    jsonData, _ := json.Marshal(loginData)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
    
    var response Response
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
    
    if !response.Success {
        t.Errorf("Expected success true, got false. Error: %s", response.Error)
    }
    
    // Проверяем наличие токена
    data, ok := response.Data.(map[string]interface{})
    if !ok {
        t.Fatal("Invalid response data format")
    }
    
    if data["token"] == nil {
        t.Error("Token not found in response")
    }
}

func TestAuthHandler_LoginInvalidCredentials(t *testing.T) {
    db := setupTestDB()
    handler := NewAuthHandler(db, "test-secret")

    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.POST("/login", handler.Login)

    loginData := map[string]string{
        "email":    "nonexistent@example.com",
        "password": "wrongpassword",
    }
    jsonData, _ := json.Marshal(loginData)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)

    if w.Code != http.StatusUnauthorized {
        t.Errorf("Expected status 401, got %d", w.Code)
    }
}