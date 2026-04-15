package contract

import (
	"cmp"
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
)

// regexScopePrefix marks a scope key as a regex pattern. The prefix is stripped
// before compilation.
const regexScopePrefix = "r#"

// cachedScopeRegex holds a compiled regex behind a sync.Once so that
// concurrent lookups of the same pattern compile it exactly once, even under
// contention. A nil re means the pattern failed to compile.
type cachedScopeRegex struct {
	once sync.Once
	re   *regexp.Regexp
}

// regexCache memoizes compiled regex patterns so that access-scope checking
// does not recompile the same regex on every request. Patterns that fail to
// compile are stored as a cachedScopeRegex with a nil re, so the failure is
// also cached.
var regexCache sync.Map // map[string]*cachedScopeRegex

// getCompiledScopeRegex returns the compiled regex for the given pattern,
// using a process-wide cache. A nil *regexp.Regexp is returned (and cached)
// for patterns that fail to compile. Concurrent callers requesting the same
// pattern will share a single compilation.
func getCompiledScopeRegex(pattern string) *regexp.Regexp {
	entryAny, _ := regexCache.LoadOrStore(pattern, &cachedScopeRegex{})
	entry := entryAny.(*cachedScopeRegex)
	entry.once.Do(func() {
		// A compile error leaves entry.re as nil, which is the documented
		// "invalid pattern" signal to callers.
		entry.re, _ = regexp.Compile(pattern)
	})
	return entry.re
}

// sortedRegexScopeKeys returns the regex-enabled keys (those starting with
// regexScopePrefix) from the given map, sorted by pattern length descending so
// that "more specific" regexes are tried first. Ties are broken
// lexicographically to give deterministic ordering.
func sortedRegexScopeKeys(scope map[string]any) []string {
	var keys []string
	for k := range scope {
		if strings.HasPrefix(k, regexScopePrefix) {
			keys = append(keys, k)
		}
	}
	slices.SortFunc(keys, func(a, b string) int {
		if d := len(b) - len(a); d != 0 {
			return d
		}
		return cmp.Compare(a, b)
	})
	return keys
}

// lookupScopeEntry finds the value associated with segment in scope. Exact
// matches always win; otherwise regex-enabled keys are tried in
// most-specific-first order (see sortedRegexScopeKeys).
func lookupScopeEntry(scope map[string]any, segment string) (any, bool) {
	if value, ok := scope[segment]; ok {
		return value, true
	}
	for _, scopeEntry := range sortedRegexScopeKeys(scope) {
		re := getCompiledScopeRegex(scopeEntry[len(regexScopePrefix):])
		if re == nil {
			continue
		}
		if re.MatchString(segment) {
			return scope[scopeEntry], true
		}
	}
	return nil, false
}

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
	if hierarchySeparator == "" {
		hierarchySeparator = "|"
	}
	pathSegments := strings.Split(key, hierarchySeparator)
	currentScope := s
	for index, segment := range pathSegments {
		value, ok := lookupScopeEntry(currentScope, segment)
		if !ok {
			return constants.ScopeAccessibilityForbidden
		}
		if nested, ok := value.(AccessScope); ok {
			currentScope = nested
			continue
		}
		if typedValue, ok := value.(string); ok {
			return s.getStringAccessibility(index, pathSegments, typedValue)
		}
		if typedValue, ok := value.(bool); ok {
			return s.getBoolAccessibility(index, pathSegments, typedValue)
		}
		return constants.ScopeAccessibilityForbidden
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
	pathSegments := strings.Split(key, ".")
	currentScope := map[string]any(s)
	for index, segment := range pathSegments {
		value, ok := lookupScopeEntry(currentScope, segment)
		if !ok {
			return nil
		}
		if nested, ok := value.(map[string]any); ok {
			currentScope = nested
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
		return nil
	}
	return nil
}

func (s FUPScope) HasLimit(key string) bool {
	pathSegments := strings.Split(key, ".")
	currentScope := map[string]any(s)
	for index, segment := range pathSegments {
		value, ok := lookupScopeEntry(currentScope, segment)
		if !ok {
			return false
		}
		if nested, ok := value.(map[string]any); ok {
			if index == len(pathSegments)-1 {
				return true
			}
			currentScope = nested
			continue
		}
		return false
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
	SetCurrentToken(apiToken ApiUserTokenInterface)
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
