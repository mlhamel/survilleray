package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	"github.com/xo/dburl"
)

// Config represent the main configuration
type Config struct {
	DatabaseURL string
	parsedURL   *dburl.URL
	httpPort    string
	database    *gorm.DB
	log         *zerolog.Logger
	statsd      *statsd.Client
}

// NewConfig create a new configuration object
func NewConfig() *Config {
	url := GetEnv("DATABASE_URL", "")
	httpPort := GetEnv("PORT", "8080")
	parsedURL, err := dburl.Parse(url)

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	log := zerolog.New(output).With().Timestamp().Logger()

	if err != nil {
		panic(err)
	}

	statsd, err := statsd.New("127.0.0.1:8125")

	if err != nil {
		panic(err)
	}

	return &Config{
		DatabaseURL: url,
		parsedURL:   parsedURL,
		httpPort:    httpPort,
		log:         &log,
		statsd:      statsd,
	}
}

func NewConfigWithDatabase(database *gorm.DB) *Config {
	cfg := NewConfig()
	cfg.database = database
	return cfg
}

func (c *Config) Database() *gorm.DB {
	if c.database != nil {
		return c.database
	}

	database, err := gorm.Open("postgres", c.DSN())

	if err != nil {
		panic(err)
	}

	c.database = database

	return c.database
}

// DSN is the connexion key to the database
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Hostname(), c.Port(), c.Username(), c.Password(), c.Path())
}

// Env return the current run level
func (c *Config) Env() string {
	env := os.Getenv("ENV")
	if env == "" {
		return "development"
	}
	return env
}

// Hostname of the configured database
func (c *Config) Hostname() string {
	return hostname(c.parsedURL.Host)
}

func (c *Config) HttpPort() string {
	return c.httpPort
}

func (c *Config) Logger() *zerolog.Logger {
	return c.log
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

func (c *Config) OpenSkyURL() string {
	return "https://opensky-network.org/api/states/all?lamin=%d&lamax=%d&lomin=%d&lomax=%d"
}

func (c *Config) Statsd() *statsd.Client {
	c.statsd.Namespace = "survilleray"
	return c.statsd
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
