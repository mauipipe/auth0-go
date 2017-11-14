package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Client runs Auth0 Management API methods whilst maintaining a fresh token
type Client struct {
	ClientID     string
	ClientSecret string
	Audience     string
	Domain       string
	token        *Token
	valid        bool
	Debug        bool
}

// NewClient returns a Client usable for executing authenticated Auth0 Management API methods
func NewClient(clientID, clientSecret, audience, domain string) (*Client, error) {
	token, err := GetToken(clientID, clientSecret, audience, domain)
	if err != nil {
		return nil, errors.Wrap(err, "Error intializing new Auth0 client")
	}

	c := &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     audience,
		Domain:       domain,
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
			token, err := GetToken(c.ClientID, c.ClientSecret, c.Audience, c.Domain)
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
	b, err := c.request("POST", endpoint, nil, input)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &output)
}

// GET sends a GET request to the specified Auth0 Management API using the client token.
//
// The `params` map will be added to the query string and the `output` interface is JSON umarshalled.
func (c *Client) GET(endpoint string, params map[string]string, output interface{}) error {
	b, err := c.request("GET", endpoint, params, nil)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &output)
}

// the low-level request function sets the headers properly
func (c *Client) request(method, endpoint string, params map[string]string, body interface{}) ([]byte, error) {
	// Ensure we have a valid client
	if !c.valid {
		return nil, errors.New("Auth0 client is not valid, try creating a new one with NewClient method")
	}

	// Construct the URL
	url := fmt.Sprintf("%s%s", c.Audience, endpoint)

	if c.Debug {
		fmt.Printf("ENDPOINT: %s\n", url)
	}

	// Marshal the payload (if we have one)
	var payloadReader io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "Error marshalling request payload")
		}

		if c.Debug {
			fmt.Printf("REQUEST: %s\n", string(b))
		}

		payloadReader = bytes.NewReader(b)
	}

	// Build the request
	req, err := http.NewRequest(method, url, payloadReader)
	if err != nil {
		return nil, errors.Wrap(err, "Error building request")
	}

	// Add the authorization header and content header
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.token.AccessToken)

	// Add query params if there are any
	if len(params) > 0 {
		q := req.URL.Query()

		for k, v := range params {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error executing request")
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	if c.Debug {
		fmt.Printf("RESPONSE: %s\n", string(resBody))
	}

	logrus.WithFields(logrus.Fields{
		"status code": res.StatusCode,
		"status":      res.Status,
		"header":      res.Header,
	}).Info("Response")

	return resBody, nil
}
