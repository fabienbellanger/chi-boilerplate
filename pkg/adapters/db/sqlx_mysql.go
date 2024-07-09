package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQL is a struct that contains the database connection
type MySQL struct {
	DB *sqlx.DB
}

// NewMySQL creates a new MySQL database connection
func NewMySQL(config *Config) (*MySQL, error) {
	dsn, err := config.dsn()
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	return &MySQL{
		DB: db,
	}, nil
}
