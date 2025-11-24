package database

import (
	"fmt"
	"log"
	"time"

	"project-defect-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
    
    gormConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    }
    
    var db *gorm.DB
    var err error
    for i := 0; i < 5; i++ {
        db, err = gorm.Open(postgres.Open(dsn), gormConfig)
        if err == nil {
            break
        }
        log.Printf("Attempt %d: failed to connect to database, retrying...", i+1)
        time.Sleep(3 * time.Second)
    }
    
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database after 5 attempts: %w", err)
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(25)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)
    
    log.Println("Connected to PostgreSQL database")
    
    if err := autoMigrate(db); err != nil {
        return nil, err
    }
    
    return db, nil
}

func autoMigrate(db *gorm.DB) error {
    models := []interface{}{
        &models.Project{},
        &models.Defect{},
        &models.DefectHistory{},
    }
    
    for _, model := range models {
        if err := db.AutoMigrate(model); err != nil {
            return fmt.Errorf("failed to migrate %T: %w", model, err)
        }
    }
    
    log.Println("Database migration completed")
    return nil
}