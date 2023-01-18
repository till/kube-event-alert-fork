package notifier

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

const contentType = "application/json"

// WebhookNotifier sends notifications to webhook
type WebhookNotifier struct {
	webhookURL string
	tplFS      embed.FS
}

// NewWebhookNotifier creates new webhook notifier
func NewWebhookNotifier(webhookURL string, tplFs embed.FS) WebhookNotifier {
	return WebhookNotifier{
		webhookURL: webhookURL,
		tplFS:      tplFs,
	}
}

// Notify convers payload to a readable message
// and sends it to webhook webhook url
func (sn WebhookNotifier) Notify(payload Payload) error {
	reader, err := sn.toReader(payload)

	if err != nil {
		return err
	}

	resp, err := http.Post(sn.webhookURL, contentType, reader)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to send notification to webhook, status %d", resp.StatusCode)
	}

	return nil
}

func (sn WebhookNotifier) toReader(payload Payload) (io.Reader, error) {
	tmpl, err := template.ParseFS(sn.tplFS, "resources/notifier.tpl")
	if err != nil {
		return nil, err
	}

	var text bytes.Buffer
	err = tmpl.Execute(&text, payload)
	if err != nil {
		return nil, err
	}

	var rawBody bytes.Buffer
	err = json.NewEncoder(&rawBody).Encode(map[string]string{
		"text":     text.String(),
		"username": "kube-event-alert",
	})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(rawBody.Bytes()), nil
}
