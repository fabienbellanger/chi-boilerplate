package cli

import (
	"chi_boilerplate/pkg/adapters/db"
	"time"

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
