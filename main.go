package main

import (
	"concurrent-checkout/src/api"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	finalizeCheckout()
}

func finalizeCheckout() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	token := os.Getenv("API_TOKEN")
	baseUrl := os.Getenv("API_BASE_URL")

	authedReq := api.NewAuthorizedApiRequest(token, baseUrl)
	getSelection := api.NewGetSelectionRequest(authedReq)

	api.SendApiRequest(&getSelection)
}
