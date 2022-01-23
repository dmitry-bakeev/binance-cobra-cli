/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// rateCmd represents the rate command
var rateCmd = &cobra.Command{
	Use:   "rate",
	Short: "A brief description of your command",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		pair, err := cmd.Flags().GetString("pair")
		if err != nil {
			return err
		}
		if pair == "" {
			return fmt.Errorf("flag --pair is required")
		}
		GetRate(pair)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rateCmd)
	rateCmd.PersistentFlags().StringP("pair", "p", "", "Pair to rate. Example: ETH-USDT")
}

type AppClient struct {
	Host string
	Port int
}

func GetRate(pair string) error {
	url := fmt.Sprintf("http://%s:%d/api/v1/rates?pairs=%s", "localhost", 3001, pair)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var result map[string]interface{}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}

	rate, ok := result[pair]
	if !ok {
		fmt.Println("Did not receive a response")
	} else {
		fmt.Println(rate)
	}
	return nil
}
