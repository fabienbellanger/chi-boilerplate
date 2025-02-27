package cli

import (
	"chi_boilerplate/pkg"
	"chi_boilerplate/pkg/adapters/db"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
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
func initConfig() (*pkg.Config, error) {
	return pkg.NewConfig()
}

// initDatabase initializes database connection.
func initDatabase(config *pkg.Config) (*db.MySQL, error) {
	return db.NewMySQL(&db.Config{
		Host:            config.Database.Host,
		Username:        config.Database.Username,
		Password:        config.Database.Password,
		Port:            config.Database.Port,
		Database:        config.Database.Database,
		Charset:         config.Database.Charset,
		Collation:       config.Database.Collation,
		Location:        config.Database.Location,
		MaxIdleConns:    config.Database.MaxIdleConns,
		MaxOpenConns:    config.Database.MaxOpenConns,
		ConnMaxLifetime: config.Database.ConnMaxLifetime,
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
