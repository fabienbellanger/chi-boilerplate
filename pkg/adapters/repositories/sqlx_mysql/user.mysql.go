package sqlx_mysql

import (
	"chi_boilerplate/pkg/adapters/db"
	"chi_boilerplate/pkg/domain/entities"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// UserMysqlRepository is an implementation of the UserRepository interface
type UserMysqlRepository struct {
	db *sqlx.DB
}

// NewUserMysqlRepository creates a new UserMysqlRepository
func NewUserMysqlRepository(db *db.MySQL) *UserMysqlRepository {
	return &UserMysqlRepository{db: db.DB}
}

// Login returns a user if a user is found
func (u UserMysqlRepository) Login(email, password string) (entities.User, error) {
	hashedPassword := entities.HashUserPassword(password)

	var user entities.User
	row := u.db.QueryRowx("SELECT * FROM users WHERE email = ? AND password = ? LIMIT 1", email, hashedPassword)
	if err := row.StructScan(&user); err != nil {
		return user, err
	}

	return user, nil
}

// Create creates a new user
func (u UserMysqlRepository) Create(user *entities.User) error {
	hashedPassword := entities.HashUserPassword(user.Password.Value)

	_, err := u.db.Exec(
		`
			INSERT INTO users (id, email, password, lastname, firstname, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
		user.ID,
		user.Email.Value,
		hashedPassword,
		user.Lastname,
		user.Firstname,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
