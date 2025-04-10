package gorm_mysql

import (
	"chi_boilerplate/pkg/adapters/db"

	"gorm.io/gorm"
)

// UserMysqlRepository is an implementation of the UserRepository interface
type UserMysqlRepository struct {
	db *gorm.DB
}

// NewUserMysqlRepository creates a new UserMysqlRepository
func NewUserMysqlRepository(db *db.GormMySQL) *UserMysqlRepository {
	return &UserMysqlRepository{db: db.DB}
}
