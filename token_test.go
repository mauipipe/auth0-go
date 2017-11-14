package auth0

import (
	"testing"
)

func Test_authEndpointFromDomain(t *testing.T) {
	type args struct {
		domain string
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
				domain: "mock.auth0.com",
			},
			want:    "https://mock.auth0.com/oauth/token",
			wantErr: false,
		},
		{
			name: "bad audience",
			args: args{
				domain: "mock.com",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authEndpointFromDomain(tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("authEndpointFromDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("authEndpointFromDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
