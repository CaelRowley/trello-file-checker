package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Board struct {
	Name string
}

type Card struct {
	ID          string
	Name        string
	BoardName   string
	Attachments []Attachment
}

type Attachment struct {
	Name string
	Date string
}

type Export struct {
	BoardName string
	CardID    string
	FileName  string
	Date      string
}

var (
	apiKey   string
	token    string
	boardIDs []string
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	apiKey = getEnv("API_KEY")
	token = getEnv("TOKEN")
	boardIDs = strings.Split(getEnv("BOARD_IDS"), ",")

	var exportData = []Export{}

	for _, boardID := range boardIDs {
		var lastCardID string
		boardHasCards := true
		var cards = []Card{}
		for boardHasCards {
			newCards := getCards(boardID, lastCardID)
			if len(newCards) > 0 {
				lastCardID = newCards[len(newCards)-1].ID
				cards = append(cards, newCards...)
			} else {
				boardHasCards = false
			}
		}
		for _, card := range cards {
			for _, attachment := range card.Attachments {
				export := Export{
					BoardName: card.BoardName,
					CardID:    card.ID,
					FileName:  attachment.Name,
					Date:      attachment.Date,
				}
				exportData = append(exportData, export)
			}
		}
	}

	exportFileData(exportData)
}

func getBoard(boardID string) Board {
	requestUrl := fmt.Sprintf("https://api.trello.com/1/boards/%s?key=%s&token=%s", boardID, apiKey, token)
	response, err := http.Get(requestUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response body: " + err.Error())
	}

	var board Board
	err = json.Unmarshal(body, &board)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	return board
}

func getCards(boardID string, lastCardID string) []Card {
	requestUrl := fmt.Sprintf("https://api.trello.com/1/boards/%s/cards?key=%s&token=%s&fields=all&attachments=true&limit=3",
		boardID, apiKey, token)
	if lastCardID != "" {
		requestUrl += "&before=" + lastCardID
	}
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
		log.Fatal("Error parsing JSON:", err)
	}

	board := getBoard(boardID)
	for i := range cards {
		cards[i].BoardName = board.Name
	}

	return cards
}

func exportFileData(data []Export) {
	file, err := os.Create(time.Now().Format("2006-01-02 15:04:05") + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Board", "Card ID", "File", "Date"})
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range data {
		err = writer.Write([]string{row.BoardName, row.CardID, row.FileName, row.Date})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal(key + " is missing from .env!")
	}
	return val
}
