package contract

import (
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"regexp"
	"strings"
	"time"
)

type AccessScope map[string]any

func (s AccessScope) getStringAccessibility(index int, pathSegments []string, typedValue string) constants.ScopeAccessibility {
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

func (s AccessScope) getBoolAccessibility(index int, pathSegments []string, typedValue bool) constants.ScopeAccessibility {
	if index == len(pathSegments)-1 {
		if typedValue {
			return constants.ScopeAccessibilityAccessible
		}
	}
	return constants.ScopeAccessibilityForbidden
}

func (s AccessScope) GetAccessibility(key string, hierarchySeparator string) constants.ScopeAccessibility {
	currentScope := s
	if hierarchySeparator == "" {
		hierarchySeparator = "|"
	}
	pathSegments := strings.Split(key, hierarchySeparator)
	index := 0
	for _, segment := range pathSegments {
		if value, ok := currentScope[segment]; ok {
			if typedValue, ok := value.(AccessScope); ok {
				currentScope = typedValue
				index++
				continue
			}
			if typedValue, ok := value.(string); ok {
				return s.getStringAccessibility(index, pathSegments, typedValue)
			}
			if typedValue, ok := value.(bool); ok {
				return s.getBoolAccessibility(index, pathSegments, typedValue)
			}
		}
		for scopeEntry, value := range currentScope {
			// regex-enabled scope keys must start with `r#`
			if strings.Index(scopeEntry, "r#") != 0 {
				continue
			}
			scopeEntryRegex, err := regexp.Compile(scopeEntry[2:])
			if nil == err {
				if scopeEntryRegex.MatchString(segment) {
					if typedValue, ok := value.(AccessScope); ok {
						currentScope = typedValue
						index++
						continue
					}
					if typedValue, ok := value.(string); ok {
						return s.getStringAccessibility(index, pathSegments, typedValue)
					}
					if typedValue, ok := value.(bool); ok {
						return s.getBoolAccessibility(index, pathSegments, typedValue)
					}
				}
			}
		}
	}
	return constants.ScopeAccessibilityForbidden
}

type FUPScope map[string]any

func (s FUPScope) getIntLimit(index int, pathSegments []string, typedValue int) *int {
	if index == len(pathSegments)-1 {
		return &typedValue
	}
	return nil
}

func (s FUPScope) getFloat64Limit(index int, pathSegments []string, typedValue float64) *int {
	if index == len(pathSegments)-1 {
		intValue := int(typedValue)
		return &intValue
	}
	return nil
}

func (s FUPScope) getFloat32Limit(index int, pathSegments []string, typedValue float32) *int {
	if index == len(pathSegments)-1 {
		intValue := int(typedValue)
		return &intValue
	}
	return nil
}

func (s FUPScope) GetLimit(key string) *int {
	currentScope := s
	pathSegments := strings.Split(key, ".")
	index := 0
	for _, segment := range pathSegments {
		if value, ok := currentScope[segment]; ok {
			if typedValue, ok := value.(map[string]any); ok {
				currentScope = typedValue
				index++
				continue
			}
			if typedValue, ok := value.(int); ok {
				return s.getIntLimit(index, pathSegments, typedValue)
			}
			if typedValue, ok := value.(float64); ok {
				return s.getFloat64Limit(index, pathSegments, typedValue)
			}
			if typedValue, ok := value.(float32); ok {
				return s.getFloat32Limit(index, pathSegments, typedValue)
			}
		}
		for scopeEntry, value := range currentScope {
			// regex-enabled scope keys must start with `r#`
			if strings.Index(scopeEntry, "r#") != 0 {
				continue
			}
			scopeEntryRegex, err := regexp.Compile(scopeEntry[2:])
			if nil == err {
				if scopeEntryRegex.MatchString(segment) {
					if typedValue, ok := value.(map[string]any); ok {
						currentScope = typedValue
						index++
						continue
					}
					if typedValue, ok := value.(int); ok {
						return s.getIntLimit(index, pathSegments, typedValue)
					}
					if typedValue, ok := value.(float64); ok {
						return s.getFloat64Limit(index, pathSegments, typedValue)
					}
					if typedValue, ok := value.(float32); ok {
						return s.getFloat32Limit(index, pathSegments, typedValue)
					}
				}
			}
		}
	}
	return nil
}

func (s FUPScope) HasLimit(key string) bool {
	currentScope := s
	pathSegments := strings.Split(key, ".")
	index := 0
	for _, segment := range pathSegments {
		if value, ok := currentScope[segment]; ok {
			if typedValue, ok := value.(map[string]any); ok {
				currentScope = typedValue
				if index == len(pathSegments)-1 {
					return true
				}
				index++
				continue
			}
			return false
		}
		for scopeEntry, value := range currentScope {
			// regex-enabled scope keys must start with `r#`
			if strings.Index(scopeEntry, "r#") != 0 {
				continue
			}
			scopeEntryRegex, err := regexp.Compile(scopeEntry[2:])
			if nil == err {
				if scopeEntryRegex.MatchString(segment) {
					if typedValue, ok := value.(map[string]any); ok {
						currentScope = typedValue
						if index == len(pathSegments)-1 {
							return true
						}
						index++
						continue
					}
					return false
				}
			}
		}
	}
	return false
}

type OneOffToken struct {
	Value   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

type ApiClientInterface interface {
	GetClientId() string
	GetClientSecret() string
	GetApiKey() string
	GetCurrentApiKey() ApiClientKeyInterface
	SetCurrentApiKey(apiClientKey ApiClientKeyInterface)
	GetClientScope() *AccessScope
	GetFUPScope() *FUPScope
}
type ApiClientKeyInterface interface {
	GetKey() string
	GetClientScope() *AccessScope
	GetFUPScope() *FUPScope
	GetApiClient() ApiClientInterface
	GetExpirationDate() *time.Time
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
	GetFUPScope() *FUPScope
	GetID() string
}
type ApiUserTokenInterface interface {
	SetToken(token string)
	GetToken() string
	SetExpirationDate(expirationDate time.Time)
	GetExpirationDate() time.Time
	SetApiUser(apiUser ApiUserInterface)
	GetApiUser() ApiUserInterface
}
