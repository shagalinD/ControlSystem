package handlers

import (
	"kopatel_online/models"
	"testing"
)

func TestDefectHandler_CreateDefect(t *testing.T) {
    db := setupTestDB()
    handler := NewDefectHandler(db)
    
    // Создаем тестовые данные
    project := models.Project{Name: "Test Project", ManagerID: 1}
    db.Create(&project)
    
    user := models.User{Email: "author@test.com", FullName: "Author", RoleID: 1}
    user.SetPassword("password")
    db.Create(&user)

    // Простой тест - проверяем что обработчик инициализирован
    if handler.DB == nil {
        t.Error("Defect handler DB is nil")
    }
}

func TestProjectHandler_CreateProject(t *testing.T) {
    db := setupTestDB()
    handler := NewProjectHandler(db)
    
    if handler.DB == nil {
        t.Error("Project handler DB is nil")
    }
}