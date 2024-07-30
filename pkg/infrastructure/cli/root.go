package cli

import (
	"chi_boilerplate/pkg/adapters/db"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "Chi Boilerplate",
	Short:   "A Chi boilerplate",
	Long:    "A Chi boilerplate",
	Version: version,
}

// Execute starts CLI
func Execute() error {
	return rootCmd.Execute()
}

// initConfig initializes configuration from config file.
func initConfig() error {
	viper.SetConfigFile(".env")
	return viper.ReadInConfig()
}

// initDatabase initializes database connection.
func initDatabase() (*db.MySQL, error) {
	return db.NewMySQL(&db.Config{
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
	})
}

func displayLogLevel(l string) aurora.Value {
	switch l {
	case "DEBUG":
		return aurora.Cyan(l)
	case "INFO":
		return aurora.Green(l)
	case "WARN":
		return aurora.Yellow(l)
	default:
		return aurora.Red(l)
	}
}

func displayLogMethod(m string) aurora.Value {
	switch m {
	case "GET":
		return aurora.Cyan(m)
	case "POST":
		return aurora.Blue(m)
	case "PUT":
		return aurora.Yellow(m)
	case "PATCH":
		return aurora.Magenta(m)
	case "DELETE":
		return aurora.Red(m)
	default:
		return aurora.Gray(12, m)
	}
}

func displayLogStatusCode(c uint) aurora.Value {
	if c < 200 {
		return aurora.Cyan(c)
	} else if c < 300 {
		return aurora.Green(c)
	} else if c < 400 {
		return aurora.Magenta(c)
	} else if c < 500 {
		return aurora.Yellow(c)
	} else {
		return aurora.Red(c)
	}
}
