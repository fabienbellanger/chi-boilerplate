package cli

import (
	"chi_boilerplate/pkg/infrastructure/chi_router"
	"chi_boilerplate/pkg/infrastructure/logger"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "run",
	Short: "Start server",
	Long:  `Start server`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func startServer() {
	config, err := initConfig()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := initDatabase(config)
	if err != nil {
		log.Fatalln(err)
	}

	l, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalln(err)
	}

	server := chi_router.NewChiServer(viper.GetString("SERVER_ADDR"), viper.GetString("SERVER_PORT"), db, l)
	if err = server.Start(); err != nil {
		log.Fatalln(err)
	}
}
