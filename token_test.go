package auth0

import (
	"testing"
)

func Test_authEndpointFromAudience(t *testing.T) {
	type args struct {
		audience string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "good audience",
			args: args{
				audience: "https://mock.auth0.com/api/v2/",
			},
			want:    "https://mock.auth0.com/oauth/token",
			wantErr: false,
		},
		{
			name: "bad audience",
			args: args{
				audience: "https://mock.auth0.com",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authEndpointFromAudience(tt.args.audience)
			if (err != nil) != tt.wantErr {
				t.Errorf("authEndpointFromAudience() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("authEndpointFromAudience() = %v, want %v", got, tt.want)
			}
		})
	}
}
