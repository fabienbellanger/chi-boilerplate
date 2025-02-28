package pkg

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// ConfigServer represents the configuration of the HTTP server
type ConfigServer struct {
	// Address
	Addr string

	// Port
	Port int

	// Timeout
	Timeout int

	// Basic Auth username
	BasicAuthUsername string

	// Basic Auth password
	BasicAuthPassword string
}

// ConfigDatabase represents the configuration of the database
type ConfigDatabase struct {
	// Driver
	Driver string

	// Host
	Host string

	// Username
	Username string

	// Password
	Password string

	// Port
	Port int

	// Database
	Database string

	// Charset
	Charset string

	// Collation
	Collation string

	// Location (UTC | Local)
	Location string

	// Max idle connections
	MaxIdleConns int

	// Max open connections
	MaxOpenConns int

	// Connection max lifetime
	ConnMaxLifetime time.Duration
}

// ConfigLog represents the configuration of the logs
type ConfigLog struct {
	// Path
	Path string

	// Outputs (stdout | file)
	Outputs []string

	// Level (ebug | info | warn | error | fatal | panic)
	Level string

	// Enable access log
	EnableAccessLog bool
}

// ConfigJWT represents the configuration of the JWT
type ConfigJWT struct {
	// Algorithm (HS512 | ES384)
	Algorithm string

	// Lifetime (in hour)
	Lifetime time.Duration

	// Secret key
	SecretKey string

	// Private key path
	PrivateKeyPath string

	// Public key path
	PublicKeyPath string
}

// NewConfigJWT creates a new ConfigJWT instance
func NewConfigJWT() (*ConfigJWT, error) {
	algo := viper.GetString("JWT_ALGO")
	secret := viper.GetString("JWT_SECRET")
	privateKeyPath := viper.GetString("JWT_PRIVATE_KEY_PATH")
	publicKeyPath := viper.GetString("JWT_PUBLIC_KEY_PATH")

	if algo != "HS512" && algo != "ES384" {
		return nil, fmt.Errorf("invalid JWT algorithm")
	}

	if algo == "HS512" && secret == "" {
		return nil, fmt.Errorf("missing JWT secret")
	}

	if algo == "ES384" && (privateKeyPath == "" || publicKeyPath == "") {
		return nil, fmt.Errorf("missing JWT private or public key path")
	}

	return &ConfigJWT{
		Algorithm:      algo,
		Lifetime:       viper.GetDuration("JWT_LIFETIME") * time.Hour,
		SecretKey:      secret,
		PrivateKeyPath: privateKeyPath,
		PublicKeyPath:  publicKeyPath,
	}, nil
}

// ConfigCORS represents the configuration of the CORS
type ConfigCORS struct {
	// Allowed origins
	AllowedOrigins []string

	// Allowed methods
	AllowedMethods []string

	// Allowed headers
	AllowedHeaders []string

	// Exposed headers
	ExposedHeaders []string

	// Allow credentials
	AllowCredentials bool

	// Max age
	MaxAge int
}

// NewConfigCORS creates a new ConfigCORS instance
func NewConfigCORS() *ConfigCORS {
	return &ConfigCORS{
		AllowedOrigins:   viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
		AllowedMethods:   viper.GetStringSlice("CORS_ALLOWED_METHODS"),
		AllowedHeaders:   viper.GetStringSlice("CORS_ALLOWED_HEADERS"),
		ExposedHeaders:   viper.GetStringSlice("CORS_EXPOSED_HEADERS"),
		AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
		MaxAge:           viper.GetInt("CORS_MAX_AGE"),
	}
}

// ConfigPprof represents the configuration of the pprof
type ConfigPprof struct {
	// Enable pprof
	Enable bool

	// Basic Auth username
	BasicAuthUsername string

	// Basic Auth password
	BasicAuthPassword string
}

// NewConfigServer creates a new ConfigServer instance
func NewConfigPprof() *ConfigPprof {
	return &ConfigPprof{
		Enable:            viper.GetBool("PPROF_ENABLE"),
		BasicAuthUsername: viper.GetString("PPROF_BASICAUTH_USERNAME"),
		BasicAuthPassword: viper.GetString("PPROF_BASICAUTH_PASSWORD"),
	}
}

// ConfigAMQP represents the configuration of the RabbitMQ server
type ConfigAMQP struct {
	// URL
	URL string
}

// NewConfigServer creates a new ConfigServer instance
func NewConfigAMQP() *ConfigAMQP {
	return &ConfigAMQP{
		URL: viper.GetString("AMQP_URL"),
	}
}

// Config represents the configuration of the application from the .env file
type Config struct {
	// Application environment (development, production or test)
	AppEnv string

	// Application name
	AppName string

	// Server configuration
	Server ConfigServer

	// Database configuration
	Database ConfigDatabase

	// Log configuration
	Log ConfigLog

	// JWT configuration
	JWT ConfigJWT

	// CORS configuration
	CORS ConfigCORS

	// Pprof configuration
	Pprof ConfigPprof

	// AMQP configuration
	AMQP ConfigAMQP
}

// NewConfig creates a new Config instance
// TODO: Add checks
func NewConfig() (*Config, error) {
	// Read .env file
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	jwtConfig, err := NewConfigJWT()
	if err != nil {
		return nil, err
	}

	return &Config{
		AppEnv:  viper.GetString("APP_ENV"),
		AppName: viper.GetString("APP_NAME"),
		Server: ConfigServer{
			Addr:              viper.GetString("SERVER_ADDR"),
			Port:              viper.GetInt("SERVER_PORT"),
			Timeout:           viper.GetInt("SERVER_TIMEOUT"),
			BasicAuthUsername: viper.GetString("SERVER_BASICAUTH_USERNAME"),
			BasicAuthPassword: viper.GetString("SERVER_BASICAUTH_PASSWORD"),
		},
		Database: ConfigDatabase{
			Driver:          viper.GetString("DB_DRIVER"),
			Host:            viper.GetString("DB_HOST"),
			Username:        viper.GetString("DB_USERNAME"),
			Password:        viper.GetString("DB_PASSWORD"),
			Port:            viper.GetInt("DB_PORT"),
			Database:        viper.GetString("DB_DATABASE"),
			Charset:         viper.GetString("DB_CHARSET"),
			Collation:       viper.GetString("DB_COLLATION"),
			Location:        viper.GetString("DB_LOCATION"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME") * time.Hour,
		},
		Log: ConfigLog{
			Path:            viper.GetString("LOG_PATH"),
			Outputs:         viper.GetStringSlice("LOG_OUTPUTS"),
			Level:           viper.GetString("LOG_LEVEL"),
			EnableAccessLog: viper.GetBool("LOG_LOG_ACCESS_ENABLE"),
		},
		JWT:   *jwtConfig,
		CORS:  *NewConfigCORS(),
		Pprof: *NewConfigPprof(),
		AMQP:  *NewConfigAMQP(),
	}, nil
}
