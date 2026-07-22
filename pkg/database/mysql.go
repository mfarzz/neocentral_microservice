package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

// EnsureDatabaseExists connects to MySQL without a specific database
// and creates the target database if it does not exist.
// This allows each microservice to be self-contained and responsible
// for provisioning its own database at startup.
func EnsureDatabaseExists(cfg Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", cfg.DBName))
	if err != nil {
		return fmt.Errorf("failed to create database %s: %w", cfg.DBName, err)
	}

	log.Printf("✅ Database '%s' is ready", cfg.DBName)
	return nil
}

// EnsureDatabaseFromDSN parses a GORM-style DSN string, extracts the database name,
// and creates it if it does not exist. This is for services that use a DSN string
// directly instead of the Config struct.
// DSN format: user:password@tcp(host:port)/dbname?params
func EnsureDatabaseFromDSN(dsn string) error {
	// Extract the database name from the DSN
	// DSN format: user:password@tcp(host:port)/dbname?params
	slashIdx := -1
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' {
			slashIdx = i
			break
		}
	}
	if slashIdx == -1 {
		return fmt.Errorf("invalid DSN format: no '/' found")
	}

	dbNamePart := dsn[slashIdx+1:]
	dbName := dbNamePart
	if qIdx := indexOf(dbNamePart, '?'); qIdx >= 0 {
		dbName = dbNamePart[:qIdx]
	}

	// Build a DSN without the database name
	rootDSN := dsn[:slashIdx+1] + "?" + "charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", rootDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	if err != nil {
		return fmt.Errorf("failed to create database %s: %w", dbName, err)
	}

	log.Printf("✅ Database '%s' is ready", dbName)
	return nil
}

func indexOf(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
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

