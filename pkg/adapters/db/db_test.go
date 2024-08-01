package db

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDsn(t *testing.T) {
	type result struct {
		dsn string
		err error
	}

	tests := []struct {
		name   string
		args   Config
		wanted result
	}{
		{
			name: "Simple valid DSN",
			args: Config{
				Username: "root",
				Password: "root",
				Database: "test",
				Host:     "localhost",
				Port:     3306,
			},
			wanted: result{
				dsn: "root:root@tcp(localhost:3306)/test?parseTime=True",
				err: nil,
			},
		},
		{
			name: "Complet valid DSN",
			args: Config{
				Username:  "root",
				Password:  "root",
				Database:  "test",
				Host:      "localhost",
				Port:      3306,
				Charset:   "utf8mb4",
				Collation: "utf8mb4_general_ci",
				Location:  "Local",
			},
			wanted: result{
				dsn: "root:root@tcp(localhost:3306)/test?parseTime=True&charset=utf8mb4&collation=utf8mb4_general_ci&loc=Local",
				err: nil,
			},
		},
		{
			name: "Invalid DSN (no username)",
			args: Config{
				Password: "root",
				Database: "test",
				Port:     3306,
				Host:     "localhost",
			},
			wanted: result{
				dsn: "",
				err: errors.New("error in database configuration"),
			},
		},
		{
			name: "Invalid DSN (no password)",
			args: Config{
				Username: "root",
				Database: "test",
				Port:     3306,
				Host:     "localhost",
			},
			wanted: result{
				dsn: "",
				err: errors.New("error in database configuration"),
			},
		},
		{
			name: "Invalid DSN (no port)",
			args: Config{
				Username: "root",
				Password: "root",
				Database: "test",
				Host:     "localhost",
			},
			wanted: result{
				dsn: "",
				err: errors.New("error in database configuration"),
			},
		},
		{
			name: "Invalid DSN (no host)",
			args: Config{
				Username: "root",
				Password: "root",
				Database: "test",
				Port:     3306,
			},
			wanted: result{
				dsn: "",
				err: errors.New("error in database configuration"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn, err := tt.args.dsn()
			got := result{dsn, err}

			if got.err != nil {
				assert.Equal(t, got.dsn, tt.wanted.dsn)
			}
			assert.Equal(t, got.err, tt.wanted.err)
		})
	}
}
