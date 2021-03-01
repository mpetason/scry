/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a card by name.",
	Long:  "Search for a card by name or other cards with the name in it.",
	Run: func(cmd *cobra.Command, args []string) {
		cardName, _ := cmd.Flags().GetString("name")

		if cardName != "" {
			getCard("exact", cardName)
		} else {
			getRandomCard()
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.PersistentFlags().String("name", "", "The exact name of a card.")
}

// Card information we care about
type Card struct {
	Name       string `json:"name"`
	CardType   string `json:"type_line"`
	ManaCost   string `json:"mana_cost"`
	OracleText string `json:"oracle_text"`
}

func getCard(searchTerm string, name string) {
	url := fmt.Sprintf("https://api.scryfall.com/cards/named?%s=%s", searchTerm, name)
	getCardResponse(url)
}

func getRandomCard() {
	url := "https://api.scryfall.com/cards/random"
	getCardResponse(url)

}

func getCardData(baseAPI string) []byte {
	request, err := http.NewRequest(
		http.MethodGet,
		baseAPI,
		nil,
	)
	if err != nil {
		log.Printf("Could not request a card. %v", err)
	}

	request.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Could not make request. %v", err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}

	return responseBytes
}

func printCardData(card Card) {
	fmt.Println(string(card.Name))
	fmt.Println(string(card.CardType))
	fmt.Println(string(card.ManaCost))
	fmt.Println(string(card.OracleText))
}

func getCardResponse(url string) {
	responseBytes := getCardData(url)
	card := Card{}

	if err := json.Unmarshal(responseBytes, &card); err != nil {
		fmt.Printf("Could not unmarshal responseBytes. %v", err)
	}

	printCardData(card)
}
