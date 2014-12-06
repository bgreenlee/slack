package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

/**
 * postResponse uses an Incoming Webhook (https://hackarts.slack.com/services/new/incoming-webhook)
 * to post a response back to Slack. This allows us to have a public response, rather
 * than the private one we would get if we just output the response directly to the
 * http.ResponseWriter
 */
func postResponse(responseURL string, responseText string, username string, icon string) {
	type SlackResponse struct {
		Text     string `json:"text"`
		Username string `json:"username"`
		Icon     string `json:"icon_url"`
	}
	response := SlackResponse{
		Text:     responseText,
		Username: username,
		Icon:     icon,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = http.Post(responseURL, "application/json", bytes.NewReader(responseJSON))
	if err != nil {
		log.Println(err)
	}
}
