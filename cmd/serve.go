package cmd

import (
	"log"
	"net/http"
	"os"

	"github.com/dhaifley/dapi/server"
	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the application server",
	Long:  "The serve command starts the application server.",
	Run: func(cmd *cobra.Command, args []string) {
		s := server.Server{Log: logrus.New()}
		s.Log.(*logrus.Logger).Out = os.Stdout
		s.Log.(*logrus.Logger).Formatter = new(logrus.JSONFormatter)
		var opts []grpc.DialOption
		creds, err := dlib.GetGRPCClientCredentials(viper.GetString("cert"))
		if err != nil {
			log.Fatalf("Failed to create client TLS credentials: %v", err)
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
		conn, err := grpc.Dial(viper.GetString("auth_url"), opts...)
		if err != nil {
			s.Log.Fatal(err)
		}

		defer conn.Close()
		s.Auth = ptypes.NewAuthClient(conn)
		s.InitRouter()
		s.Log.Fatal(http.ListenAndServe(":3611", s.Router))
	},
}
