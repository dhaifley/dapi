package cmd

import (
	"errors"

	"github.com/dhaifley/dapi/client"
	"github.com/dhaifley/dlib/dauth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	password  string
	clientURL string
)

func init() {
	dauthLoginCmd.Flags().StringVarP(&clientURL, "url", "u", "https://dapi.io", "URL for client connections")
	dauthLoginCmd.Flags().StringVarP(&password, "password", "p", "", "password for login command")
	dauthLogoutCmd.Flags().StringVarP(&clientURL, "url", "u", "https://dapi.io", "URL for client connections")
	rootCmd.AddCommand(dauthCmd)
	dauthCmd.AddCommand(dauthLoginCmd)
	dauthCmd.AddCommand(dauthLogoutCmd)
}

var dauthCmd = &cobra.Command{
	Use:   "dauth",
	Short: "Sends commands to the dauth service",
	Long:  "The dauth command sends commands to the authentication service.",
}

var dauthLoginCmd = &cobra.Command{
	Use:   "login <user>",
	Short: "Authenticates and obtains an API token",
	Long:  "The login command authenticates a user and obtains a token for API access.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("no user specified for login")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		rc, err := client.NewRESTClient(
			clientURL+"/dauth",
			clientURL+"/dauth",
			viper.GetString("cert"))
		if err != nil {
			log.Error(err)
			return
		}

		ch := rc.Login(&dauth.User{User: args[0], Pass: password})
		for res := range ch {
			if res.Err != nil {
				log.Error(res.Err)
				return
			}

			viper.Set("token", rc.Token.Token)
			if err := viper.WriteConfig(); err != nil {
				log.Error(err)
				return
			}

			log.Info("Login successful")
			log.Infoln(rc.Token)
		}
	},
}

var dauthLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Destroys an API token",
	Long:  "The login command destroys a token for API access.",
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.GetString("token")
		if token != "" {
			rc, err := client.NewRESTClient(
				clientURL+"/dauth",
				clientURL+"/dauth",
				viper.GetString("cert"))
			if err != nil {
				log.Error(err)
				return
			}

			rc.Token = &dauth.Token{Token: token}
			ch := rc.Logout()
			for res := range ch {
				if res.Err != nil {
					log.Error(res.Err)
					return
				}

				viper.Set("token", "none")
				if err := viper.WriteConfig(); err != nil {
					log.Error(err)
					return
				}
			}
		}

		log.Info("Logout successful")
	},
}
