package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Token is the access token returned from the oauth endpoint
type Token struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// GetToken authenticates with the oauth endpoint and returns a token
//
// The ClientID and Audience can be found on your Auth0 APIs page, located
// at https://manage.auth0.com/#/apis/.  (Audience is labeled as Identifier).
// The Domain can be found on your Auth0 Clients page, located at
// https://manage.auth0.com/#/clients/ under your Client's "Settings"
func GetToken(clientID, clientSecret, audience, domain string) (*Token, error) {
	url, err := authEndpointFromDomain(domain)

	if err != nil {
		return nil, errors.Wrap(err, "Error getting token request URL")
	}

	payload := struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Audience     string `json:"audience"`
		GrantType    string `json:"grant_type"`
	}{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     audience,
		GrantType:    "client_credentials",
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "Error marshalling JSON token request payload")
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, "Error building token request")
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error sending token request")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	var t Token
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling token response")
	}

	if t.Error != "" {
		return &t, fmt.Errorf("Error: %s - %s", t.Error, t.ErrorDescription)
	}

	return &t, nil
}

// e.g. https://mock.auth0.com/oauth/token from mock.auth0.com
func authEndpointFromDomain(domain string) (string, error) {
	parts := strings.Split(domain, ".")

	lastIdx := len(parts) - 1
	end := strings.Join(parts[lastIdx-1:], ".")

	if end != "auth0.com" {
		return "", fmt.Errorf("Bad Domain URL '%s', should look like 'mock.auth0.com'", domain)
	}

	// TODO: more validations?

	return fmt.Sprintf("https://%s/oauth/token", domain), nil
}
