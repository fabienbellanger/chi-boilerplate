package db

import (
	"errors"
	"fmt"
	"time"
)

const (
	// MaxLimit represents the max number of items for pagination
	MaxLimit = 100
)

// Config represents the MySQL database configuration
type Config struct {
	Host            string
	Username        string
	Password        string
	Port            int
	Database        string
	Charset         string
	Collation       string
	Location        string
	MaxIdleConns    int           // Sets the maximum number of connections in the idle connection pool
	MaxOpenConns    int           // Sets the maximum number of open connections to the database
	ConnMaxLifetime time.Duration // Sets the maximum amount of time a connection may be reused
	ConnMaxIdleTime time.Duration // Sets the maximum amount of time a connection in the idle may be reused
}

// dsn returns the DSN if the configuration is OK or an error in other case
// TODO: Add test
func (c *Config) dsn() (dsn string, err error) {
	if c.Host == "" || c.Port == 0 || c.Username == "" || c.Password == "" {
		return dsn, errors.New("error in database configuration")
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database)
	if c.Charset != "" {
		dsn += fmt.Sprintf("&charset=%s", c.Charset)
	}
	if c.Collation != "" {
		dsn += fmt.Sprintf("&collation=%s", c.Collation)
	}
	if c.Location != "" {
		dsn += fmt.Sprintf("&loc=%s", c.Location)
	}
	return
}
