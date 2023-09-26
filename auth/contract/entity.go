package contract

import (
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"strings"
	"time"
)

type AccessScope map[string]any

func (s AccessScope) GetAccessibility(key string) constants.ScopeAccessibility {
	currentScope := s
	pathSegments := strings.Split(key, ".")
	index := 0
	for _, segment := range pathSegments {
		if value, ok := currentScope[segment]; ok {
			if typedValue, ok := value.(AccessScope); ok {
				currentScope = typedValue
				index++
				continue
			}
			if typedValue, ok := value.(string); ok {
				if index == len(pathSegments)-1 {
					if typedValue == string(constants.ScopeAccessibilityAccessible) {
						return constants.ScopeAccessibilityAccessible
					}
					if typedValue == string(constants.ScopeAccessibilityOnBehalf) {
						return constants.ScopeAccessibilityOnBehalf
					}
				}
				return constants.ScopeAccessibilityForbidden
			}
			if typedValue, ok := value.(bool); ok {
				if index == len(pathSegments)-1 {
					if typedValue {
						return constants.ScopeAccessibilityAccessible
					}
				}
				return constants.ScopeAccessibilityForbidden
			}
		}
	}
	return constants.ScopeAccessibilityForbidden
}

type ApiClientInterface interface {
	GetClientId() string
	GetClientSecret() string
	GetApiKey() string
	GetClientScope() *AccessScope
}
type ApiUserInterface interface {
	AddApiToken(apiToken ApiUserTokenInterface)
	GetCurrentToken() ApiUserTokenInterface
	GetUserScope() *AccessScope
	GetLastLoginAt() *time.Time
	SetLastLoginAt(lastLoginAt *time.Time)
	GetPassword() string
	SetPassword(password string)
	GetLogin() string
	SetLogin(login string)
	SetConfirmationToken(confirmationToken *string)
	GetConfirmationRequestedAt() *time.Time
	SetConfirmationRequestedAt(confirmationRequestedAt *time.Time)
	IsActive() bool
	SetActive(active bool)
	GetResetRequestedAt() *time.Time
	SetResetRequestedAt(resetRequestedAt *time.Time)
	GetResetToken() *string
	SetResetToken(resetToken *string)
}
type ApiUserTokenInterface interface {
	SetToken(token string)
	GetToken() string
	SetExpirationDate(expirationDate time.Time)
	GetExpirationDate() time.Time
	SetApiUser(apiUser ApiUserInterface)
	GetApiUser() ApiUserInterface
}
