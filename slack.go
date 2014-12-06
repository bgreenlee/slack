package main

/**
 * Bare-bones slash command handler for Slack
 */

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/bgreenlee/slack/config"
	"github.com/bgreenlee/slack/handlers"
)

func main() {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	config := config.Config{}
	if err := decoder.Decode(&config); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		handlers.Weather(w, r, config)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
