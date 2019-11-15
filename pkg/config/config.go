package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xo/dburl"
)

// Config represent the main configuration
type Config struct {
	DatabaseURL string
	parsedURL   *dburl.URL
	httpPort    string
	db          *gorm.DB
}

// NewConfig create a new configuration object
func NewConfig() *Config {
	url := GetEnv("DATABASE_URL", "")
	httpPort := GetEnv("PORT", "8080")
	parsedURL, err := dburl.Parse(url)

	if err != nil {
		panic(err)
	}

	return &Config{
		DatabaseURL: url,
		parsedURL:   parsedURL,
		httpPort:    httpPort,
	}
}

func NewConfigWithDB(db *gorm.DB) *Config {
	return &Config{
		DatabaseURL: "",
		parsedURL:   nil,
		httpPort:    GetEnv("PORT", "8080"),
		db:          db,
	}
}

// DSN is the connexion key to the database
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Hostname(), c.Port(), c.Username(), c.Password(), c.Path())
}

// DB return the connexion to the dabase
func (c *Config) DB() *gorm.DB {
	if c.db != nil {
		return c.db
	}

	db, err := gorm.Open("postgres", c.DSN())
	if err != nil {
		panic(err)
	}

	c.db = db

	return c.db
}

// Hostname of the configured database
func (c *Config) Hostname() string {
	return hostname(c.parsedURL.Host)
}

func (c *Config) HttpPort() string {
	return c.httpPort
}

// Port of the configured database
func (c *Config) Port() string {
	return hostport(c.parsedURL.Host)
}

// Path of the configured database
func (c *Config) Path() string {
	return strings.TrimPrefix(c.parsedURL.Path, "/")
}

// Username of the configured database
func (c *Config) Username() string {
	return c.parsedURL.User.Username()
}

// Password of the configured database
func (c *Config) Password() string {
	p, _ := c.parsedURL.User.Password()

	return p
}

// Env return the current run level
func (c *Config) Env() string {
	env := os.Getenv("ENV")
	if env == "" {
		return "development"
	}
	return env
}

// GetEnv return the current `key` value or `fallback`.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func hostname(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

func hostport(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return ""
	}
	if i := strings.Index(hostport, "]:"); i != -1 {
		return hostport[i+len("]:"):]
	}
	if strings.Contains(hostport, "]") {
		return ""
	}
	return hostport[colon+len(":"):]
}
