package notifier_test

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ronenlib/kube-event-alert/pkg/notifier"
)

//go:embed resources
var tplTestFs embed.FS

type webhookRequest struct {
	Text string `json:"text"`
}

var body webhookRequest

func setWebhookBody(req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		panic("failed to parse json")
	}
}

func fakeOkHandler(w http.ResponseWriter, req *http.Request) {
	setWebhookBody(req)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")
}

func fakeBadMessageHandler(w http.ResponseWriter, req *http.Request) {
	setWebhookBody(req)
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "error")
}

func TestNotify(t *testing.T) {
	cases := []struct {
		handler         http.Handler
		payload         notifier.Payload
		expectedMessage string
		expectError     bool
		name            string
	}{
		{
			handler: http.HandlerFunc(fakeOkHandler),
			payload: notifier.Payload{
				Kind:      "kind",
				Namespace: "namespace",
				Name:      "name",
				Error:     "error",
			},
			expectedMessage: "name, kind: error",
			expectError:     false,
			name:            "NotifySuccessful",
		},
		{
			handler:         http.HandlerFunc(fakeBadMessageHandler),
			payload:         notifier.Payload{},
			expectedMessage: "",
			expectError:     true,
			name:            "NotifyFailure",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			sn := notifier.NewWebhookNotifier(server.URL, tplTestFs)

			body = webhookRequest{} // reset variable
			err := sn.Notify(tc.payload)

			if tc.expectError {
				receivedErr := err != nil
				if err == nil {
					t.Errorf("Expected error to be %v but received error was %v", tc.expectError, receivedErr)
				}
			} else {
				// strip \n
				if tc.expectedMessage != strings.TrimSpace(body.Text) {
					t.Errorf("Notifier sent the wrong message: %s (expected: %s)", body.Text, tc.expectedMessage)
				}
			}
		})
	}
}
