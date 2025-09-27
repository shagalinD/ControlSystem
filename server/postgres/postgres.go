package postgres

import (
	"fmt"
	"kopatel_online/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
}

func NewConnection(cfg *Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    log.Println("Connected to PostgreSQL database")
    
    // Автомиграция
    if err := autoMigrate(db); err != nil {
        return nil, err
    }
    
    // Создание начальных данных
    if err := seedData(db); err != nil {
        return nil, err
    }
    
    return db, nil
}

func autoMigrate(db *gorm.DB) error {
    models := []interface{}{
        &models.Role{},
        &models.User{},
        &models.Project{},
        &models.Defect{},
        &models.Comment{},
    }
    
    for _, model := range models {
        if err := db.AutoMigrate(model); err != nil {
            return fmt.Errorf("failed to migrate %T: %w", model, err)
        }
    }
    
    log.Println("Database migration completed")
    return nil
}

func seedData(db *gorm.DB) error {
    // Создание ролей по умолчанию
    roles := []models.Role{
        {RoleName: "engineer"},
        {RoleName: "manager"},
        {RoleName: "observer"},
    }
    
    for _, role := range roles {
        var existingRole models.Role
        if err := db.Where("role_name = ?", role.RoleName).First(&existingRole).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                if err := db.Create(&role).Error; err != nil {
                    return fmt.Errorf("failed to create role %s: %w", role.RoleName, err)
                }
                log.Printf("Created role: %s", role.RoleName)
            } else {
                return err
            }
        }
    }
    
    return nil
}