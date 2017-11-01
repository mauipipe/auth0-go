package auth0

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Client runs Auth0 Management API methods whilst maintaining a fresh token
type Client struct {
	ClientID     string
	ClientSecret string
	Audience     string
	token        *Token
	valid        bool
}

// NewClient returns a Client usable for executing authenticated Auth0 Management API methods
func NewClient(clientID, clientSecret, audience string) (*Client, error) {
	token, err := GetToken(clientID, clientSecret, audience)
	if err != nil {
		return nil, errors.Wrap(err, "Error intializing new Auth0 client")
	}

	c := &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     audience,
		token:        token,
		valid:        true,
	}
	c.startTokenRefresher()

	return c, nil
}

// startTokenRefresher uses the reported token expiry time to automatically refresh the
// token 5 seconds before it expires.
//
// If the GetToken calls produces an error, we set the Client valid flag to false so the
// next call can try to refresh the token or error.  That's about all we do for retry logic
// at this time ¯\_(ツ)_/¯
func (c *Client) startTokenRefresher() {
	refresher := time.NewTicker(time.Second * time.Duration(c.token.ExpiresIn-5))
	go func() {
		for _ = range refresher.C {
			token, err := GetToken(c.ClientID, c.ClientSecret, c.Audience)
			if err != nil {
				fmt.Printf("Error refreshing Auth0 token: %s", err.Error())

				// set the client status to valid such that the next time someone
				// tries to use it, they will get an error.
				c.valid = false
				return
			}

			c.token = token
		}
	}()
}

// POST sends a POST request to the specified Auth0 Management API using the client token.
//
// The `input` and `output` interfaces are JSON marshalled/unmarshalled.
func (c *Client) POST(endpoint string, input interface{}, output interface{}) error {
	// TODO: fill me in!
	return nil
}

// GET sends a GET request to the specified Auth0 Management API using the client token.
//
// The `params` map will be added to the query string and the `output` interface is JSON umarshalled.
func (c *Client) GET(endpoint string, params map[string]string, output interface{}) error {
	// TODO: fill me in!
	return nil
}
