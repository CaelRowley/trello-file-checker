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
)

type Board struct {
	Name string
}

type List struct {
	ID   string
	Name string
}

type Card struct {
	ID          string
	Name        string
	BoardName   string
	ListName    string
	Attachments []Attachment
}

type Attachment struct {
	Name string
	Date string
}

type Export struct {
	BoardName string
	ListName  string
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
	apiKey = getEnv("TRELLO_API_KEY")
	token = getEnv("TRELLO_TOKEN")
	boardIDs = strings.Split(getEnv("TRELLO_BOARD_IDS"), ",")

	var exportData = []Export{}

	for _, boardID := range boardIDs {
		board := getBoard(boardID)
		lists := getListsOnBoard(boardID)
		for _, list := range lists {
			var lastCardID string
			listHasCards := true
			var cards = []Card{}
			for listHasCards {
				newCards := getCardsOnList(list.ID, lastCardID)
				if len(newCards) > 0 {
					lastCardID = newCards[0].ID
					cards = append(cards, newCards...)
				} else {
					listHasCards = false
				}
			}
			for _, card := range cards {
				for _, attachment := range card.Attachments {
					export := Export{
						BoardName: board.Name,
						ListName:  list.Name,
						CardID:    card.ID,
						FileName:  attachment.Name,
						Date:      attachment.Date,
					}
					exportData = append(exportData, export)
				}
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

func getListsOnBoard(boardID string) []List {
	requestUrl := fmt.Sprintf("https://api.trello.com/1/boards/%s/lists?key=%s&token=%s", boardID, apiKey, token)
	response, err := http.Get(requestUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response body: " + err.Error())
	}

	var lists []List
	err = json.Unmarshal(body, &lists)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	return lists
}

func getCardsOnList(listID string, lastCardID string) []Card {
	requestUrl := fmt.Sprintf("https://api.trello.com/1/lists/%s/cards?key=%s&token=%s&attachments=true", listID, apiKey, token)
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

	err = writer.Write([]string{"Board", "List", "Card ID", "File", "Date"})
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range data {
		err = writer.Write([]string{row.BoardName, row.ListName, row.CardID, row.FileName, row.Date})
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
