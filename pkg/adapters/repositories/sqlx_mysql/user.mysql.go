package sqlx_mysql

import (
	"chi_boilerplate/pkg/adapters/db"
	"chi_boilerplate/pkg/domain/repositories"
	"chi_boilerplate/pkg/domain/requests"
	"chi_boilerplate/pkg/domain/responses"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// UserMysqlRepository is an implementation of the UserRepository interface
type UserMysqlRepository struct {
	db *sqlx.DB
}

// NewUserMysqlRepository creates a new UserMysqlRepository
func NewUserMysqlRepository(db *db.SqlxMySQL) *UserMysqlRepository {
	return &UserMysqlRepository{db: db.DB}
}

// Create creates a new user
func (u *UserMysqlRepository) Create(user requests.UserCreationRepository) error {
	_, err := u.db.Exec(`
		INSERT INTO users (id, email, password, lastname, firstname, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.ID,
		user.Email,
		user.Password,
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

// GetByID returns a user by ID
func (u *UserMysqlRepository) GetByID(req requests.UserByID) (responses.UserByIdRepository, error) {
	var user responses.UserByIdRepository
	row := u.db.QueryRowx(`
		SELECT id, email, lastname, firstname, created_at, updated_at
		FROM users 
		WHERE id = ?
			AND deleted_at IS NULL
		LIMIT 1`,
		req.ID,
	)
	if err := row.StructScan(&user); err != nil {
		return user, repositories.ErrUserNotFound
	}

	return user, nil
}

// GetByID returns a user by Email
func (u *UserMysqlRepository) GetByEmail(req requests.GetByEmail) (responses.GetByEmail, error) {
	var user responses.GetByEmailRepository
	row := u.db.QueryRowx(`
		SELECT id, password
		FROM users 
		WHERE email = ?
			AND deleted_at IS NULL
		LIMIT 1`,
		req.Email,
	)
	if err := row.StructScan(&user); err != nil {
		return responses.GetByEmail{}, repositories.ErrUserNotFound
	}

	response, err := user.ToGetByEmail()
	if err != nil {
		return responses.GetByEmail{}, repositories.ErrUserNotFound // TODO: change error type
	}

	return response, nil
}

// Delete deletes a user
func (u *UserMysqlRepository) Delete(req requests.UserDelete) error {
	result, err := u.db.Exec(`
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = ?
			AND deleted_at IS NULL`,
		req.ID,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return repositories.ErrUserNotFound
	}

	return err
}

func (u *UserMysqlRepository) CountAll() (int64, error) {
	var count int64
	row := u.db.QueryRowx(`
		SELECT COUNT(id)
		FROM users 
		WHERE deleted_at IS NULL
	`)
	if err := row.Scan(&count); err != nil {
		return count, err
	}

	return count, nil
}

func (u *UserMysqlRepository) GetAll(req requests.UsersList) ([]responses.UsersListRepository, error) {
	offset, limit := db.PaginateValues(req.Page, req.Limit)
	query_sort := db.OrderValues(req.Sorts)

	query := `
		SELECT id, email, lastname, firstname, created_at, updated_at
		FROM users 
		WHERE deleted_at IS NULL`

	if len(query_sort) > 0 {
		query += query_sort
	}
	query += " LIMIT ? OFFSET ?"

	rows, err := u.db.Queryx(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []responses.UsersListRepository
	for rows.Next() {
		var user responses.UsersListRepository
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserMysqlRepository) Update(req requests.UserUpdateRepository) error {
	result, err := u.db.Exec(`
		UPDATE users
		SET lastname = ?, firstname = ?, email = ?, password = ?, updated_at = ?
		WHERE id = ?
			AND deleted_at IS NULL`,
		req.Lastname,
		req.Firstname,
		req.Email,
		req.Password,
		req.UpdatedAt,
		req.ID,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return repositories.ErrUserNotFound
	}

	return err
}
