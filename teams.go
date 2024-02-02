package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"k8s.io/klog/v2"
)

type MicrosoftTeams struct {
	WebhookURL string
	MuteSeconds    int
	ClusterName    string
	History map[string]time.Time
}

type MicrosoftTeamsMessage struct {
	Text       string `json:"text,omitempty"`
	Summary    string `json:"summary,omitempty"`
	ThemeColor string `json:"themeColor,omitempty"`
}

func NewMicrosoftTeams() MicrosoftTeams {
	var teamsWebhookURL string

	if teamsWebhookURL = os.Getenv("MICROSOFT_TEAMS_WEBHOOK_URL"); teamsWebhookURL == "" {
		klog.Exit("Environment variable MICROSOFT_TEAMS_WEBHOOK_URL is not set")
	}

	return MicrosoftTeams{
		WebhookURL: teamsWebhookURL,
	}
}

func (m MicrosoftTeams) TeamsSendMessage(msg MicrosoftTeamsMessage) error {
	client := &http.Client{}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", m.WebhookURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Example usage:
// func main() {
// 	teams := NewMicrosoftTeams()

// 	message := MicrosoftTeamsMessage{
// 		Text:       "This is the main text of the message",
// 		Summary:    "Summary of the message",
// 		ThemeColor: "#4599DF", // Customize as needed
// 	}

// 	err := teams.SendMessage(message)
// 	if err != nil {
// 		fmt.Printf("Error sending message to Microsoft Teams: %v\n", err)
// 	} else {
// 		fmt.Println("Message sent successfully to Microsoft Teams.")
// 	}
// }
