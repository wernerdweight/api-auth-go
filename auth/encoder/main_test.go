package encoder

import (
	"github.com/stretchr/testify/assert"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

type mockApiUser struct {
	password string
}

func (m mockApiUser) AddApiToken(apiToken contract.ApiUserTokenInterface) {}

func (m mockApiUser) GetCurrentToken() contract.ApiUserTokenInterface {
	return nil
}

func (m mockApiUser) GetUserScope() *contract.AccessScope {
	return nil
}

func (m mockApiUser) GetLastLoginAt() *time.Time {
	return nil
}

func (m mockApiUser) SetLastLoginAt(lastLoginAt *time.Time) {}

func (m mockApiUser) GetPassword() string {
	return m.password
}

func (m mockApiUser) SetPassword(password string) {
	m.password = password
}

func (m mockApiUser) GetLogin() string {
	return ""
}

func (m mockApiUser) SetLogin(login string) {}

func (m mockApiUser) SetConfirmationToken(confirmationToken *string) {}

func (m mockApiUser) GetConfirmationRequestedAt() *time.Time {
	return nil
}

func (m mockApiUser) SetConfirmationRequestedAt(confirmationRequestedAt *time.Time) {}

func (m mockApiUser) IsActive() bool {
	return false
}

func (m mockApiUser) SetActive(active bool) {}

func (m mockApiUser) GetResetRequestedAt() *time.Time {
	return nil
}

func (m mockApiUser) SetResetRequestedAt(resetRequestedAt *time.Time) {}

func (m mockApiUser) GetResetToken() *string {
	return nil
}

func (m mockApiUser) SetResetToken(resetToken *string) {}

func TestEncoder_ComparePassword(t *testing.T) {
	assertion := assert.New(t)
	type args struct {
		apiUser  mockApiUser
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid password (bcrypt)",
			args: args{apiUser: mockApiUser{password: "$2a$10$jDL769dJ7XlZ7YiYGoDUFO4jCZmUjEYS7wCdgOvT/hV4JDwvcy7.G"}, password: "password"},
			want: true,
		},
		{
			name: "Valid password (argon)",
			args: args{apiUser: mockApiUser{password: "$argon2id$v=19$m=65536,t=1,p=2$c+l5u1BO+xsi9i8eXDPpCw$OlrkxToUI3ruy0smOlP6bBbJy/WbgZnPHW/5wx8aW8E"}, password: "password"},
			want: true,
		},
		{
			name: "Invalid password",
			args: args{apiUser: mockApiUser{password: "password"}, password: "invalid"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComparePassword(tt.args.apiUser, tt.args.password)
			if tt.want {
				assertion.Nil(got)
				return
			}
			assertion.NotNil(got)
		})
	}
}

func TestEncoder_EncryptPassword(t *testing.T) {
	assertion := assert.New(t)
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Password 1",
			args: args{password: "password"},
		},
		{
			name: "Password 2",
			args: args{password: "password"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptPassword(tt.args.password)
			assertion.Nil(err)
			assertion.Nil(bcrypt.CompareHashAndPassword([]byte(got), []byte(tt.args.password)))
		})
	}
	p1encrypted, err := EncryptPassword("password")
	assertion.Nil(err)
	p2encrypted, err := EncryptPassword("password")
	assertion.Nil(err)
	assertion.NotEqual(p1encrypted, p2encrypted)
}
