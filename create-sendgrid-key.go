package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
)

// APIKeyRequest defines the structure for the API key creation request payload.
type APIKeyRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

// CreateAPIKeyResponse defines the structure to parse the successful API key creation response.
type CreateAPIKeyResponse struct {
	APIKey   string `json:"api_key"`
	Name     string `json:"name"`
	ApiKeyID string `json:"api_key_id"`
}

// definedScopes holds the complete list of permissions for the new API key.
var definedScopes = []string{
	"mail.send",
	"mail.batch.create",
	"mail.batch.delete",
	"mail.batch.read",
	"mail.batch.update",
	"mail_settings.address_whitelist.read",
	"mail_settings.address_whitelist.update",
	"mail_settings.bcc.read",
	"mail_settings.bcc.update",
	"mail_settings.bounce_purge.read",
	"mail_settings.bounce_purge.update",
	"mail_settings.footer.read",
	"mail_settings.footer.update",
	"mail_settings.forward_bounce.read",
	"mail_settings.forward_bounce.update",
	"mail_settings.forward_spam.read",
	"mail_settings.forward_spam.update",
	"mail_settings.plain_content.read",
	"mail_settings.plain_content.update",
	"mail_settings.read",
	"mail_settings.spam_check.read",
	"mail_settings.spam_check.update",
	"mail_settings.template.read",
	"mail_settings.template.update",
	"user.scheduled_sends.create",
	"user.scheduled_sends.delete",
	"user.scheduled_sends.read",
	"user.scheduled_sends.update",
	"user.webhooks.event.settings.create",
	"user.webhooks.event.settings.read",
	"user.webhooks.event.settings.update",
	"user.webhooks.event.settings.delete",
	"user.webhooks.event.test.create",
	"user.webhooks.event.test.read",
	"user.webhooks.event.test.update",
}

func main() {

	// 1. Get the new API key name from command-line arguments
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run your_script_name.go \"New API Key Name\"")
		os.Exit(1)
	}
	newApiKeyName := os.Args[1]

	// 2. Get the Admin API Key (Bearer Token) from environment variable
	adminApiKey := os.Getenv("SENDGRID_API_KEY")
	if adminApiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: SENDGRID_API_KEY environment variable not set.")
		os.Exit(1)
	}

	// 3. Prepare the request body using the single definedScopes list
	requestPayload := APIKeyRequest{
		Name:   newApiKeyName,
		Scopes: definedScopes,
	}

	requestBodyBytes, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling request body: %v\n", err)
		os.Exit(1)
	}

	// 4. Prepare and send the SendGrid API request
	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest(adminApiKey, "/v3/api_keys", host)
	request.Method = "POST"
	request.Body = requestBodyBytes

	response, err := sendgrid.API(request)
	if err != nil {
		// Error from the sendgrid.API call itself (e.g., network issue)
		fmt.Fprintf(os.Stderr, "Error calling SendGrid API: %v\n", err)
		if response != nil { // If response is not nil, it might contain some details
			fmt.Fprintf(os.Stderr, "Response Status Code (if available): %d\n", response.StatusCode)
			fmt.Fprintf(os.Stderr, "Response Body (if available): %s\n", response.Body)
		}
		os.Exit(1)
	}

	// 5. Process the response from SendGrid
	if response.StatusCode == 201 { // 201 Created indicates success
		var apiKeyResp CreateAPIKeyResponse
		// response.Body is a string, so we convert it to []byte for Unmarshal
		if err := json.Unmarshal([]byte(response.Body), &apiKeyResp); err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling successful response body: %v\n", err)
			fmt.Fprintf(os.Stderr, "Raw Response Body: %s\n", response.Body)
			os.Exit(1)
		}

		if apiKeyResp.APIKey == "" {
			fmt.Fprintln(os.Stderr, "Error: API key not found in successful SendGrid response.")
			fmt.Fprintf(os.Stderr, "Raw Response Body: %s\n", response.Body)
			os.Exit(1)
		}
		// Only print the API key to stdout on success
		fmt.Println(apiKeyResp.APIKey)
	} else {
		// Handle non-201 status codes (API errors from SendGrid)
		fmt.Fprintf(os.Stderr, "Error: SendGrid API responded with status code %d.\n", response.StatusCode)
		fmt.Fprintf(os.Stderr, "Response Body: %s\n", response.Body)
		// You might want to parse the response.Body for specific SendGrid error messages here if needed
		os.Exit(1)
	}
}
