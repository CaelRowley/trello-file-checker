package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
	exportFileData(cards)
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

func exportFileData(data []Card) {
	file, err := os.Create(time.Now().Format("2006-01-02 15:04:05") + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Board", "File", "Date"})
	if err != nil {
		log.Fatal(err)
	}

	for _, card := range data {
		for _, attachment := range card.Attachments {
			err = writer.Write([]string{boardID, attachment.Name, attachment.Date})
		}
	}

	fmt.Println(file.Name())
}

func getEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal(key + " is missing from .env!")
	}
	return val
}
