package cli

import (
	"chi_boilerplate/pkg/adapters/repositories/sqlx_mysql"
	"chi_boilerplate/pkg/domain/requests"
	"chi_boilerplate/pkg/domain/usecases"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	userEmail     string
	userPassword  string
	userLastname  string
	userFirstname string
)

func init() {
	userCmd.Flags().StringVarP(&userLastname, "lastname", "l", "", "user lastname")
	userCmd.Flags().StringVarP(&userFirstname, "firstname", "f", "", "user firstname")
	userCmd.Flags().StringVarP(&userEmail, "email", "e", "", "user email")
	userCmd.Flags().StringVarP(&userPassword, "password", "p", "", "user password")

	userCmd.MarkFlagRequired("lastname")
	userCmd.MarkFlagRequired("firstname")
	userCmd.MarkFlagRequired("email")
	userCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(userCmd)
}

var userCmd = &cobra.Command{
	Use:   "register",
	Short: "User creation",
	Long:  `User creation`,
	Run: func(cmd *cobra.Command, args []string) {
		user := requests.UserCreation{
			Lastname:  strings.TrimSpace(userLastname),
			Firstname: strings.TrimSpace(userFirstname),
			Password:  strings.TrimSpace(userPassword),
			Email:     strings.TrimSpace(userEmail),
		}

		// Initialize configuration
		if err := initConfig(); err != nil {
			fmt.Printf("\nError: %v\n", err)
			return
		}

		// Initialize database
		db, err := initDatabase()
		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			return
		}

		// Call use case
		userRepo := sqlx_mysql.NewUserMysqlRepository(db)
		userUseCase := usecases.NewUser(userRepo)
		res, errRes := userUseCase.Create(user)
		if errRes != nil {
			fmt.Printf("\nError: %v (%v)\n", errRes.Message, errRes.Details)
			return
		}

		// Display result
		fmt.Printf(`
User successfully created:
    - ID:        %s
    - Lastname:  %s
    - Firstname: %s
    - Email:     %s
    - Password:  %s
`,
			res.ID,
			res.Lastname,
			res.Firstname,
			res.Email.Value,
			user.Password,
		)
	},
}
