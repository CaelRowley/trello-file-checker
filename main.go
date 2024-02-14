package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(".env is missing! " + err.Error())
	}
	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		panic("API_KEY is missing!")
	}
	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		panic("TOKEN is missing!")
	}
	boardId, ok := os.LookupEnv("BOARD_ID")
	if !ok {
		panic("BOARD_ID is missing!")
	}

	requestUrl := fmt.Sprintf("https://api.trello.com/1/boards/%s/cards?key=%s&token=%s&fields=all&attachments=true", boardId, apiKey, token)
	response, err := http.Get(requestUrl)
	if err != nil {
		panic("Error: " + err.Error())
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic("Error reading response body: " + err.Error())
	}

	fmt.Println("Response Status:", response.Status)
	fmt.Println("Response Body:", string(body))
}
