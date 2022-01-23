/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Provide rates from Binance",
	Long:  "",
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
	serverCmd.Flags().StringP("host", "H", "localhost", "Set listeting host")
	serverCmd.Flags().IntP("port", "p", 3001, "Set listeting port")
}

type App struct {
	Host string
	Port int
}

type Rate struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type PostJSON struct {
	Pairs []string `json:"pairs"`
}

type RateSlice []Rate

const API_URL = "https://api.binance.com/api/v3/ticker/price"

func GetCurrentRates() (RateSlice, error) {
	response, err := http.Get(API_URL)
	if err != nil {
		return RateSlice{}, err
	}
	defer response.Body.Close()

	var curretRates RateSlice
	err = json.NewDecoder(response.Body).Decode(&curretRates)
	if err != nil {
		return RateSlice{}, err
	}

	return curretRates, nil
}

func FindRates(pairs []string) (RateSlice, error) {
	var result RateSlice

	curretRates, err := GetCurrentRates()
	if err != nil {
		return RateSlice{}, err
	}

	for _, pair := range pairs {
		for _, currentRate := range curretRates {
			tmp := strings.Replace(pair, "-", "", -1)
			if tmp == currentRate.Symbol {
				result = append(result, Rate{pair, currentRate.Price})
			}
		}
	}

	return result, nil
}

func RatesToJSON(rates RateSlice) ([]byte, error) {
	rateMap := make(map[string]float64)
	for _, rate := range rates {
		tmp, err := strconv.ParseFloat(rate.Price, 64)
		if err != nil {
			return []byte{}, err
		}
		rateMap[rate.Symbol] = tmp
	}

	result, err := json.Marshal(rateMap)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

func IntermediateHandler(rw http.ResponseWriter, r *http.Request, pairs []string) {
	rates, err := FindRates(pairs)
	if err != nil {
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

	result, err := RatesToJSON(rates)
	if err != nil {
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(result)
}

func RootHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pairsStr := r.URL.Query().Get("pairs")
		if pairsStr == "" {
			http.Error(rw, "Bad request", http.StatusBadRequest)
			return
		}
		pairs := strings.Split(pairsStr, ",")
		IntermediateHandler(rw, r, pairs)
	case http.MethodPost:
		receivedJSON, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}
		var postJSON PostJSON
		err = json.Unmarshal(receivedJSON, &postJSON)
		if err != nil {
			http.Error(rw, "Bad request", http.StatusBadRequest)
			return
		}
		IntermediateHandler(rw, r, postJSON.Pairs)
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func startServer(a App) {
	http.HandleFunc("/api/v1/rates", RootHandler)
	http.ListenAndServe(fmt.Sprintf("%s:%d", a.Host, a.Port), nil)
}
