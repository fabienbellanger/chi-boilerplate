package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// SqlxMySQL is a struct that contains the database connection
type SqlxMySQL struct {
	DB     *sqlx.DB
	config *Config
}

// NewSqlxMySQL creates a new MySQL database connection
func NewSqlxMySQL(config *Config) (*SqlxMySQL, error) {
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

	return &SqlxMySQL{
		DB:     db,
		config: config,
	}, nil
}

func (m *SqlxMySQL) DSN() (string, error) {
	return m.config.dsn()
}

func (m *SqlxMySQL) Database(d string) {
	m.config.Database = d
}
