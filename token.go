package auth0

// This does things yadda yada yadya dalskdjkasd
// 	func (t *Token) Thing() {
//		log.Println(lb)
//  }

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
func GetToken(clientID, clientSecret, audience string) (*Token, error) {
	url, err := authEndpointFromAudience(audience)

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

// e.g. https://mock.auth0.com/oauth/token from
// https://mock.auth0.com/api/v2/
func authEndpointFromAudience(audience string) (string, error) {
	parts := strings.Split(audience, "/")

	if len(parts) != 6 {
		return "", fmt.Errorf("Bad Audience URL '%s', should look like 'https://mock.auth0.com/api/v2/'", audience)
	}
	// TODO: more validations?

	return fmt.Sprintf("%s/oauth/token", strings.Join(parts[0:3], "/")), nil
}
