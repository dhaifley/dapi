// Package cmd provides command line processing functions.
package cmd

import (
	"fmt"
	"os"

	"github.com/dhaifley/dapi/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   lib.ServiceInfo.Name,
	Short: lib.ServiceInfo.Short,
	Long:  lib.ServiceInfo.Long,
}

func init() {
	viper.SetConfigFile("dapi_config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	viper.SetEnvPrefix(lib.ServiceInfo.Name)
	viper.SetDefault("auth_url", "localhost:3612")
	if err := viper.BindEnv("auth_url"); err != nil {
		fmt.Println(err)
	}

	viper.SetDefault("cert", "")
	if err := viper.BindEnv("cert"); err != nil {
		fmt.Println(err)
	}

	viper.SetDefault("token", "")
	if err := viper.BindEnv("token"); err != nil {
		fmt.Println(err)
	}
}

// Execute starts the command processor.
func Execute() {
	fmt.Println(lib.ServiceInfo.Short)
	fmt.Println("Version:", lib.ServiceInfo.Version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
