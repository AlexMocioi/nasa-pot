package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"nasa-pot/src/config"
	"net/http"
	"time"
)

// Slack sends message to a Slack hook for quick alerts.
func Slack(parts ...interface{}) error {
	marshalledJSON, err := json.Marshal(struct {
		Text string `json:"text"`
	}{
		Text: fmt.Sprintf("[%s]: %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprint(parts...)),
	})
	if err != nil {
		return fmt.Errorf("json.Marshal(): %s", err)
	}
	resp, err := http.Post(config.Configuration.SlackHookURL, "application/json", bytes.NewBuffer(marshalledJSON))
	if err != nil {
		return fmt.Errorf("http.Post(): %s", err)
	}
	defer resp.Body.Close()
	return nil
}
