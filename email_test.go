package auth0

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestEmailResponse_Pending(t *testing.T) {
	tests := []struct {
		tag    string
		status string
		Type   string
		expRes error
	}{
		{
			tag:    `valid`,
			status: EmailStatusPending,
			Type:   EmailType,
			expRes: nil,
		},
		{
			tag:    `wrong type`,
			status: EmailStatusCompleted,
			Type:   `wrong type`,
			expRes: ErrVerifyEmail,
		},
		{
			tag:    `wrong status`,
			status: `wrong status`,
			Type:   EmailType,
			expRes: ErrVerifyEmail,
		},
		{
			tag:    `wrong status and type`,
			status: `wrong status`,
			Type:   `wrong type`,
			expRes: ErrVerifyEmail,
		},
	}

	for _, tc := range tests {
		t.Run(tc.tag, func(t *testing.T) {
			res := EmailResponse{
				Status: tc.status,
				Type:   tc.Type,
			}

			assert.Equal(t, tc.expRes, res.Pending())
		})
	}
}
