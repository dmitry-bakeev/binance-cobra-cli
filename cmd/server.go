/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			return err
		}

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}

		a := App{host, port}

		log.Printf("Listening http://%s:%d", a.Host, a.Port)
		startServer(a)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serverCmd.Flags().StringP("host", "H", "localhost", "Set listeting host")
	serverCmd.Flags().IntP("port", "p", 3001, "Set listeting port")
}

type App struct {
	Host string
	Port int
}

func RootHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, "OK!\n")
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func startServer(a App) {
	http.HandleFunc("/", RootHandler)
	http.ListenAndServe(fmt.Sprintf("%s:%d", a.Host, a.Port), nil)
}
