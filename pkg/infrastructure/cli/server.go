package cli

import (
	"chi_boilerplate/pkg/infrastructure/chi_router"
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
	// Configuration initialization
	// ----------------------------
	err := initConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// Start server
	// ------------
	server := chi_router.NewChiServer(viper.GetString("SERVER_ADDR"), viper.GetString("SERVER_PORT"))
	err = server.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
