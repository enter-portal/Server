package database

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"portal/internal/server/models"
	"portal/internal/server/utils"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
)

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	sslmode    = os.Getenv("DB_SSL_MODE")
	dbInstance *gorm.DB
)

func New() *gorm.DB {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// Check the environment variable for the database type
	dbType := os.Getenv("DB")

	var db *gorm.DB
	var err error

	if dbType == "postgres" {

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, username, password, database, port, sslmode)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		dsn, ok := os.LookupEnv("DB_PATH")
		if !ok {
			appName, ok := os.LookupEnv("APP_NAME")
			if !ok {
				appName = "portal"
			}
			dsn = fmt.Sprintf("%s.db", appName)
		}
		log.Printf("Falling back on SQLite: %s\n", dsn)
		db, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s", dsn)), &gorm.Config{})
	}

	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})
	dbInstance = db
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func Health() (int, utils.JSON) {
	stats := utils.JSON{}

	if dbInstance == nil {
		stats["message"] = "Failed to disconnected: dbInstance is not nil"
		return http.StatusBadRequest, stats
	}

	// Get the underlying DB object
	sqlDB, err := dbInstance.DB()
	if err != nil {
		stats["message"] = err.Error()
		return http.StatusBadRequest, stats
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Ping the database
	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("%s", fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return http.StatusBadRequest, stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return http.StatusOK, stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func Close() error {
	if dbInstance == nil {
		return fmt.Errorf("Failed to close: dbInstance is nil")
	}

	// Get the underlying DB object
	sqlDB, err := dbInstance.DB()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Disconnected from database: %s", database)
	// Close
	return sqlDB.Close()
}
