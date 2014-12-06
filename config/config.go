package config

type Config struct {
	ResponseType string `json:"response_type"` // webook, slackbot, direct
	ResponseURL  string `json:"response_url"`  // incoming webhook url or slackbot url
	SlackToken   string `json:"slack_token"`   // token used to identify requests coming from slack
}
