package db

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// MaxLimit represents the max number of items for pagination
	MaxLimit = 100

	// DefaultSlowThreshold represents the default slow threshold value
	DefaultSlowThreshold time.Duration = 200 * time.Millisecond
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
	SlowThreshold   time.Duration // Slow SQL threshold (Default: 200ms)
}

// dsn returns the DSN if the configuration is OK or an error in other case
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

// PaginateValues transforms page and limit into offset and limit.
func PaginateValues(p, l string) (offset int, limit int) {
	page, err := strconv.Atoi(p)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err = strconv.Atoi(l)
	if err != nil || limit > MaxLimit || limit < 1 {
		limit = MaxLimit
	}

	offset = (page - 1) * limit

	return
}

// orderValues transforms list of fields to sort into a map.
func orderValues(list string, prefixes ...string) []string {
	r := make([]string, 0)

	if len(list) <= 0 {
		return r
	}

	prefix := ""
	if len(prefixes) == 1 {
		prefix = prefixes[0] + "."
	}

	sorts := strings.Split(list, ",")
	for _, s := range sorts {
		if len(s) > 0 {
			key := fmt.Sprintf("%s%s", prefix, s[1:])
			if strings.HasPrefix(s, "+") && len(s[1:]) > 1 {
				r = append(r, fmt.Sprintf("%s ASC", key))
			} else if strings.HasPrefix(s, "-") && len(s[1:]) > 1 {
				r = append(r, fmt.Sprintf("%s DESC", key))
			}
		}
	}

	return r
}

// OrderValues returns the ORDER BY clause for a list of fields to sort.
func OrderValues(list string, prefixes ...string) (s string) {
	values := orderValues(list, prefixes...)
	s = strings.Join(values, ", ")

	if len(s) > 0 {
		s = " ORDER BY " + s
	}

	return
}
