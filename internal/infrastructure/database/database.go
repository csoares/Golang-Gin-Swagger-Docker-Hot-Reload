package database

import (
	"io/ioutil"
	"strings"

	postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB holds the database connection
type DB struct {
	*gorm.DB
}

// Config holds database configuration
type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

// New creates a new database connection
func New(config Config) (*DB, error) {
	dsn := "host=" + config.Host +
		" user=" + config.Username +
		" password=" + config.Password +
		" dbname=" + config.DBName +
		" port=" + config.Port +
		" sslmode=disable TimeZone=Europe/Lisbon"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// ReadConfig reads database configuration from file
func ReadConfig(filepath string) (Config, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	lines := strings.Split(string(content), "\n")
	config := Config{}

	if len(lines) >= 6 {
		config.Username = strings.TrimSpace(lines[1])
		config.Password = strings.TrimSpace(lines[2])
		config.Host = strings.TrimSpace(lines[3])
		config.Port = strings.TrimSpace(lines[4])
		config.DBName = strings.TrimSpace(lines[5])
	}

	return config, nil
}
