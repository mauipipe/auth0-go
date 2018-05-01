package auth0

import (
	`time`
	"github.com/pkg/errors"
)

const (
	EmailStatusCompleted = `completed`
	EmailStatusPending   = `pending`

	EmailType = `verification_email`
)

var (
	// ErrVerifyEmail error on sending email
	ErrVerifyEmail = errors.New(`err sending verification email`)
)

// VerificationPostBody body to be sent to verify an user email address
type VerificationPostBody struct {
	UserID   string `json:"user_id"`
	ClientID string `json:"client_id"`
}

// EmailResponse response from verification call
// see https://auth0.com/docs/api/management/v2#!/Jobs/post_verification_email
type EmailResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// Pending checks if the response belong to an email_verification response
func (vr *EmailResponse) Pending() error {
	if vr.Status != EmailStatusPending || vr.Type != EmailType {
		return ErrVerifyEmail
	}

	return nil
}

// SendVerificationEmail send a verification email to the user
// see https://auth0.com/docs/api/management/v2#!/Jobs/post_verification_email
func (c *Client) SendVerificationEmail(userID string) error {
	vpb := VerificationPostBody{
		UserID:   userID,
		ClientID: c.ClientID,
	}

	var vr EmailResponse
	if err := c.POST("jobs/verification-email", vpb, &vr); err != nil {
		return errors.Wrap(err, `error on verification call`)
	}

	return nil
}
