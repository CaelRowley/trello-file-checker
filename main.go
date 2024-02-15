package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Card struct {
	Name        string
	Attachments []Attachment
}

type Attachment struct {
	Name string
	Date string
}

var (
	apiKey  string
	token   string
	boardID string
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	apiKey = getEnv("API_KEY")
	token = getEnv("TOKEN")
	boardID = getEnv("BOARD_ID")

	cards := getCards()
	fmt.Println(cards)
}

func getCards() []Card {
	requestUrl := fmt.Sprintf("https://api.trello.com/1/boards/%s/cards?key=%s&token=%s&fields=all&attachments=true",
		boardID, apiKey, token)
	response, err := http.Get(requestUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response body: " + err.Error())
	}

	var cards []Card
	if err := json.Unmarshal(body, &cards); err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	return cards
}

func getEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal(key + " is missing from .env!")
	}
	return val
}
