package tests

import (
	"chi_boilerplate/pkg/adapters/db"
	"chi_boilerplate/pkg/infrastructure/chi_router"
	"chi_boilerplate/pkg/infrastructure/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	UserEmail    = "test@test.com"
	UserPassword = "00000000"
)

// Test defines a structure for specifying input and output data of a single test case.
type Test struct {
	Description string

	// Test input
	Route   string
	Method  string
	Body    io.Reader
	Headers []Header

	// Check
	CheckError bool
	CheckBody  bool
	CheckCode  bool

	// Expected output
	ExpectedError bool
	ExpectedCode  int
	ExpectedBody  string
}

// Header represents an header value.
type Header struct {
	Key   string
	Value string
}

// TestMysql is used to create and use a database for tests.
type TestMysql struct {
	name  string
	DB    *db.MySQL
	Token string
}

// newTestMysql returns a TestMysql instance.
func newTestMysql() (TestMysql, error) {
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
	// TODO !
	// dbt.MakeMigrations()

	// Create first user and get token
	token, err := createUserAndAuthenticate(dbt)
	if err != nil {
		return TestMysql{}, err
	}

	return TestMysql{DB: dbt, name: dbName, Token: token}, nil
}

// Drop database after the test.
func (tdb *TestMysql) Drop() error {
	_, err := tdb.DB.DB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`;", tdb.name))

	return err
}

// Init initializes configuration from .env path and returns TestMysql.
func Init(p string) TestMysql {
	viper.SetConfigFile(p)
	viper.ReadInConfig()

	viper.Set("APP_ENV", "test")
	viper.Set("JWT_ALGO", "HS512")
	viper.Set("JWT_SECRET", "mySecretForTest")
	viper.Set("SERVER_PPROF", false)
	viper.Set("ENABLE_ACCESS_LOG", false)

	tdb, err := newTestMysql()
	if err != nil {
		log.Panicf("%v\n", err)
	}
	return tdb
}

// Execute runs all tests.
func Execute(t *testing.T, db *db.MySQL, tests []Test, templatesPath string) {
	// Set up the app as it is done in the main function
	l, _ := logger.NewZapLogger()
	s := chi_router.NewChiServer("localhost", "7777", db, l)
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

// TODO
func createUserAndAuthenticate(db *db.MySQL) (token string, err error) {
	return
}

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.ServeHTTP(rr, req)

	return rr
}

// JsonToString converts a JSON to a string.
func JsonToString(d interface{}) string {
	b, err := json.Marshal(d)
	if err != nil {
		log.Panicf("%v\n", err)
	}
	return string(b)
}
