package sqlx_mysql

import "github.com/jmoiron/sqlx"

type MysqlRepository struct {
	db *sqlx.DB
}

func NewMysqlRepository() (*MysqlRepository, error) {
	db, err := sqlx.Open("mysql", "root:root@tcp()")
	if err != nil {
		return nil, err
	}

	return &MysqlRepository{
		db: db,
	}, nil
}
