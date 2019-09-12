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
}

// NewConfig create a new configuration object
func NewConfig() *Config {
	url := getEnv("DATABASE_URL", "")
	httpPort := getEnv("PORT", ":8080")
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

// DSN is the connexion key to the database
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s sslmode=disable dbname=%s",
		c.Hostname(), c.Port(), c.Username(), c.Path())
}

// DB return the connexion to the dabase
func (c *Config) DB() *gorm.DB {
	db, err := gorm.Open("postgres", c.DSN())
	if err != nil {
		panic(err)
	}

	return db
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
