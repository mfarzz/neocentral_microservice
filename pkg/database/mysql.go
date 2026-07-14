package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds MySQL connection parameters.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	MaxOpen  int
	MaxIdle  int
	MaxLife  time.Duration
}

// NewMySQLConnection creates a new GORM DB instance connected to MySQL.
func NewMySQLConnection(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true, // better performance for read-heavy workloads
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if cfg.MaxOpen <= 0 {
		cfg.MaxOpen = 25
	}
	if cfg.MaxIdle <= 0 {
		cfg.MaxIdle = 10
	}
	if cfg.MaxLife <= 0 {
		cfg.MaxLife = 5 * time.Minute
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetConnMaxLifetime(cfg.MaxLife)

	log.Printf("✅ Connected to MySQL database: %s", cfg.DBName)
	return db, nil
}
