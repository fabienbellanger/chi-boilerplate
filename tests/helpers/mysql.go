package helpers

import (
	"chi_boilerplate/pkg/adapters/db"
	"chi_boilerplate/pkg/adapters/repositories/sqlx_mysql"
	"chi_boilerplate/pkg/domain/entities"
	"chi_boilerplate/pkg/domain/requests"
	vo "chi_boilerplate/pkg/domain/value_objects"
	"chi_boilerplate/pkg/infrastructure/chi_router"
	"chi_boilerplate/pkg/infrastructure/logger"
	"chi_boilerplate/utils"

	"net/http"
	"testing"

	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestMysql is used to create and use a database for tests.
type TestMysql struct {
	name  string
	DB    *db.MySQL
	Token string
}

// newTestMysql returns a TestMysql instance.
func newTestMysql(m string) (TestMysql, error) {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	dbName := viper.GetString("DB_DATABASE") + "__" + fmt.Sprintf("%08d", rand.Int63n(1e8))
	config := db.Config{
		Host:            viper.GetString("DB_HOST"),
		Username:        viper.GetString("DB_USERNAME"),
		Password:        viper.GetString("DB_PASSWORD"),
		Port:            viper.GetInt("DB_PORT"),
		Database:        "",
		Charset:         viper.GetString("DB_CHARSET"),
		Collation:       viper.GetString("DB_COLLATION"),
		Location:        viper.GetString("DB_LOCATION"),
		MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
		MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
		ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME") * time.Hour,
	}

	dbt, err := db.NewMySQL(&config)
	if err != nil {
		return TestMysql{}, err
	}

	// Create database for test, use it and run migrations
	dbt.DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", dbName))
	dbt.DB.Exec(fmt.Sprintf("USE `%s`;", dbName))
	dbt.Database(dbName)
	err = runMySQLMigrations(m, dbt)
	if err != nil {
		return TestMysql{}, err
	}

	// Create first user and get token
	token, err := createMySQLUserAndAuthenticate(dbt)
	if err != nil {
		return TestMysql{}, err
	}

	return TestMysql{DB: dbt, name: dbName, Token: token}, nil
}

// Init initializes configuration from .env path and returns TestMysql.
func InitMySQL(p, m string) TestMysql {
	viper.SetConfigFile(p)
	viper.ReadInConfig()

	viper.Set("APP_ENV", "test")
	viper.Set("JWT_ALGO", "HS512")
	viper.Set("JWT_SECRET", "mySecretForTest")
	viper.Set("SERVER_PPROF", false)
	viper.Set("ENABLE_ACCESS_LOG", false)

	tdb, err := newTestMysql(m)
	if err != nil {
		log.Panicf("%v\n", err)
	}
	return tdb
}

// Execute runs all tests.
func (tdb *TestMysql) Execute(t *testing.T, tests []Test, templatesPath string) {
	// Set up the app as it is done in the main function
	l, _ := logger.NewZapLogger()
	s := chi_router.NewChiServer("", "", tdb.DB, l)
	app, _ := s.Setup()

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route from the test case
		req, _ := http.NewRequest(test.Method, test.Route, test.Body)
		for _, h := range test.Headers {
			req.Header.Add(h.Key, h.Value)
		}

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res := executeRequest(req, app)

		// Verify if the status code is as expected
		if test.CheckCode {
			assert.Equalf(t, test.ExpectedCode, res.Code, test.Description)
		}

		// Verify if the body is as expected
		if test.CheckBody {
			// Read the response body
			body, err := io.ReadAll(res.Body)

			// Reading the response body should work everytime, such that
			// the err variable should be nil
			assert.Nilf(t, err, test.Description)

			// Verify, that the response body equals the expected body
			assert.Equalf(t, test.ExpectedBody, string(body), test.Description)
		}
	}
}

// Drop database after the test.
func (tdb *TestMysql) Drop() error {
	_, err := tdb.DB.DB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`;", tdb.name))

	return err
}

func runMySQLMigrations(m string, db *db.MySQL) error {
	newDSN, err := db.DSN()
	if err != nil {
		return err
	}

	mg, err := migrate.New(
		fmt.Sprintf("file://%s", m),
		fmt.Sprintf("mysql://%s", newDSN))
	if err != nil {
		return err
	}

	err = mg.Up()
	if err != nil {
		return err
	}

	return nil
}

func createMySQLUserAndAuthenticate(db *db.MySQL) (token string, err error) {
	// Create first user
	created_at, err := time.Parse(time.RFC3339, UserCreatedAt)
	if err != nil {
		return
	}
	updated_at, err := time.Parse(time.RFC3339, UserUpdatedAt)
	if err != nil {
		return
	}
	userID, err := vo.NewIDFrom(UserID)
	if err != nil {
		return "", err
	}
	userRepo := sqlx_mysql.NewUserMysqlRepository(db)
	password, err := entities.HashUserPassword(UserPassword)
	if err != nil {
		return "", err
	}
	err = userRepo.Create(requests.UserCreationRepository{
		ID:        userID.String(),
		Lastname:  "Test",
		Firstname: "Test",
		Password:  password,
		Email:     UserEmail,
		CreatedAt: created_at.Format(utils.SqlDateTimeFormat),
		UpdatedAt: updated_at.Format(utils.SqlDateTimeFormat),
	})
	if err != nil {
		return
	}

	// Generate JWT token
	token, _, err = entities.GenerateJWT(userID, viper.GetDuration("JWT_LIFETIME"), viper.GetString("JWT_ALGO"), viper.GetString("JWT_SECRET"))
	if err != nil {
		return
	}

	return
}
