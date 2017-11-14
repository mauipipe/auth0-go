package auth0

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

// UserCreateParams are the parameters for creating a user
type UserCreateParams struct {
	UserID        string            `json:"user_id,omitempty"`
	Connection    string            `json:"connection"` // required
	Email         string            `json:"email,omitempty"`
	Username      string            `json:"username,omitempty"`
	Password      string            `json:"password,omitempty"`
	PhoneNumber   string            `json:"phone_number,omitempty"`
	EmailVerified bool              `json:"email_verified,omitempty"`
	VerifyEmail   bool              `json:"verify_email,omitempty"`
	PhoneVerified bool              `json:"phone_verified,omitempty"`
	AppMetadata   map[string]string `json:"app_metadata,omitempty"`
	UserMetadata  map[string]string `json:"user_metadata,omitempty"`
}

// User is a created Auth0 user
type User struct {
	UserID        string            `json:"user_id,omitempty"`
	Email         string            `json:"email,omitempty"`
	Username      string            `json:"username,omitempty"`
	PhoneNumber   string            `json:"phone_number,omitempty"`
	EmailVerified bool              `json:"email_verified,omitempty"`
	VerifyEmail   bool              `json:"verify_email,omitempty"`
	PhoneVerified bool              `json:"phone_verified,omitempty"`
	AppMetadata   map[string]string `json:"app_metadata,omitempty"`
	UserMetadata  map[string]string `json:"user_metadata,omitempty"`
	CreatedAt     time.Time         `json:"created_at,omitempty"`
	UpdatedAt     time.Time         `json:"updated_at,omitempty"`
	Picture       string            `json:"picture,omitempty"`
	Name          string            `json:"name,omitempty"`
	Nickname      string            `json:"nickname,omitempty"`
	LastIP        string            `json:"last_ip,omitempty"`
	LastLogin     string            `json:"last_login,omitempty"`
	LoginsCount   int               `json:"logins_count,omitempty"`
	Blocked       bool              `json:"blocked,omitempty"`
	GivenName     string            `json:"given_name,omitempty"`
	FamilyName    string            `json:"family_name,omitempty"`
}

// UserCreate creates a user.
// see https://auth0.com/docs/api/management/v2#!/Users/post_users
func (c *Client) UserCreate(p *UserCreateParams) (*User, error) {
	var u User

	if err := c.POST("users", p, &u); err != nil {
		return nil, errors.Wrap(err, "Error creating user")
	}

	logrus.WithFields(logrus.Fields{
		"email": u.Email,
		"md":    u.UserMetadata,
		"id":    u.UserID,
	}).Info("User Created")

	return &u, nil
}
