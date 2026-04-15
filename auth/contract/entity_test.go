package contract

import (
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"regexp"
	"testing"
)

func TestAccessScope_GetAccessibility(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name  string
		scope AccessScope
		args  args
		want  constants.ScopeAccessibility
	}{
		{
			name:  "Empty scope",
			scope: AccessScope{},
			args:  args{path: "/"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Empty scope, empty path",
			scope: AccessScope{},
			args:  args{path: ""},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Empty scope, path is not in scope",
			scope: AccessScope{},
			args:  args{path: "/"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Path is in scope, root",
			scope: AccessScope{"/": true},
			args:  args{path: "/"},
			want:  constants.ScopeAccessibilityAccessible,
		},
		{
			name:  "Path is not in scope",
			scope: AccessScope{"/": true},
			args:  args{path: "/test"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Path is in scope, path not root",
			scope: AccessScope{"/test": true},
			args:  args{path: "/test"},
			want:  constants.ScopeAccessibilityAccessible,
		},
		{
			name:  "Path is not in scope, but superseded",
			scope: AccessScope{"/test": true},
			args:  args{path: "/"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Path is in scope, multiple keys 1",
			scope: AccessScope{"/test": true, "/test2": true},
			args:  args{path: "/test"},
			want:  constants.ScopeAccessibilityAccessible,
		},
		{
			name:  "Path is in scope, multiple keys 2",
			scope: AccessScope{"/test": true, "/test2": true},
			args:  args{path: "/test2"},
			want:  constants.ScopeAccessibilityAccessible,
		},
		{
			name:  "Path is not in scope, multiple keys",
			scope: AccessScope{"/test": true, "/test2": true},
			args:  args{path: "/test3"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Path is in scope, but false, multiple keys",
			scope: AccessScope{"/test": true, "/test2": false},
			args:  args{path: "/test2"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Path is in scope, on-behalf, multiple keys",
			scope: AccessScope{"/test": true, "/test2": "on-behalf"},
			args:  args{path: "/test2"},
			want:  constants.ScopeAccessibilityOnBehalf,
		},
		{
			name:  "Nested scope, true",
			scope: AccessScope{"test": AccessScope{"nested1": true, "nested2": false, "nested3": "on-behalf"}},
			args:  args{path: "test|nested1"},
			want:  constants.ScopeAccessibilityAccessible,
		},
		{
			name:  "Nested scope, false",
			scope: AccessScope{"test": AccessScope{"nested1": true, "nested2": false, "nested3": "on-behalf"}},
			args:  args{path: "test|nested2"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Nested scope, on-behalf",
			scope: AccessScope{"test": AccessScope{"nested1": true, "nested2": false, "nested3": "on-behalf"}},
			args:  args{path: "test|nested3"},
			want:  constants.ScopeAccessibilityOnBehalf,
		},
		{
			name:  "Nested scope, not in scope",
			scope: AccessScope{"test": AccessScope{"nested1": true, "nested2": false, "nested3": "on-behalf"}},
			args:  args{path: "test|nested4"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Nested scope, parent",
			scope: AccessScope{"test": AccessScope{"nested1": true, "nested2": false, "nested3": "on-behalf"}},
			args:  args{path: "test"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Nested scope, no path",
			scope: AccessScope{"test": AccessScope{"nested1": true, "nested2": false, "nested3": "on-behalf"}},
			args:  args{path: ""},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Nested scope, too deep",
			scope: AccessScope{"test": AccessScope{"nested1": true, "nested2": false, "nested3": "on-behalf"}},
			args:  args{path: "test|test|test|nope"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Path is in scope, on-behalf, multiple keys, with wrong (not-enabled) regex",
			scope: AccessScope{"/test": true, "/test/[^/]+$": "on-behalf"},
			args:  args{path: "/test/abC-1De23f"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Path is in scope, on-behalf, multiple keys, with regex",
			scope: AccessScope{"/test": true, "r#^/test/[^/]+$": "on-behalf"},
			args:  args{path: "/test/abC-1De23f"},
			want:  constants.ScopeAccessibilityOnBehalf,
		},
		{
			name:  "Path is in scope, multiple keys, with regex",
			scope: AccessScope{"/test": true, "r#^/test/[^/]+/?$": true},
			args:  args{path: "/test/abC-1De23f/"},
			want:  constants.ScopeAccessibilityAccessible,
		},
		{
			name:  "Path not in scope, multiple keys, with regex",
			scope: AccessScope{"/test": true, "r#^/test/[^/]+/?$": "on-behalf"},
			args:  args{path: "/test/abC-1De23f/abcd"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Path in scope, multiple keys, with regex, with unsafe chars",
			scope: AccessScope{"/test": true, "r#^/test/.*?$": "on-behalf"},
			args:  args{path: "/test/abC-1De23f/Test_2023-01-02T12:13:14.567Z"},
			want:  constants.ScopeAccessibilityOnBehalf,
		},
		{
			name:  "Nested scope, on-behalf, multiple keys, with regex",
			scope: AccessScope{"/test": true, "nested1": AccessScope{"r#^/test/[^/]+$": AccessScope{"nested2": "on-behalf"}}},
			args:  args{path: "nested1|/test/abC-1De23f|nested2"},
			want:  constants.ScopeAccessibilityOnBehalf,
		},
		// Specificity-ordering cases: longer (more specific) regex keys must
		// be tried before shorter ones, so deny rules can shadow broader allow
		// rules regardless of map iteration order.
		{
			name: "Regex ordering: specific deny shadows broad allow (chat)",
			scope: AccessScope{
				"r#^/similarity/v[0-9]+/.*$":                true,
				"r#^/similarity/v[0-9]+/(chat|sessions).*$": false,
			},
			args: args{path: "/similarity/v1/chat"},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Regex ordering: specific deny shadows broad allow (sessions nested)",
			scope: AccessScope{
				"r#^/similarity/v[0-9]+/.*$":                true,
				"r#^/similarity/v[0-9]+/(chat|sessions).*$": false,
			},
			args: args{path: "/similarity/v2/sessions/123/messages"},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Regex ordering: broad allow still matches when specific deny does not",
			scope: AccessScope{
				"r#^/similarity/v[0-9]+/.*$":                true,
				"r#^/similarity/v[0-9]+/(chat|sessions).*$": false,
			},
			args: args{path: "/similarity/v1/anything"},
			want: constants.ScopeAccessibilityAccessible,
		},
		{
			name: "Regex ordering: unrelated path still forbidden",
			scope: AccessScope{
				"r#^/similarity/v[0-9]+/.*$":                true,
				"r#^/similarity/v[0-9]+/(chat|sessions).*$": false,
			},
			args: args{path: "/other/v1/anything"},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Regex ordering: full example with exact allow",
			scope: AccessScope{
				"/token/generate":                           true,
				"r#^/similarity/v[0-9]+/.*$":                true,
				"r#^/similarity/v[0-9]+/(chat|sessions).*$": false,
			},
			args: args{path: "/token/generate"},
			want: constants.ScopeAccessibilityAccessible,
		},
		{
			name: "Regex ordering: exact match beats regex even when regex is more specific",
			scope: AccessScope{
				"/similarity/v1/chat":                       true,
				"r#^/similarity/v[0-9]+/(chat|sessions).*$": false,
			},
			args: args{path: "/similarity/v1/chat"},
			want: constants.ScopeAccessibilityAccessible,
		},
		{
			name: "Regex ordering: on-behalf from more specific regex wins over broad allow",
			scope: AccessScope{
				"r#^/admin/.*$":       true,
				"r#^/admin/users/.*$": "on-behalf",
			},
			args: args{path: "/admin/users/42"},
			want: constants.ScopeAccessibilityOnBehalf,
		},
		{
			name: "Regex ordering: broad on-behalf used when specific deny does not match",
			scope: AccessScope{
				"r#^/admin/.*$":                 "on-behalf",
				"r#^/admin/internal/secret/.+$": false,
			},
			args: args{path: "/admin/users/42"},
			want: constants.ScopeAccessibilityOnBehalf,
		},
		{
			name: "Regex ordering: specific deny wins over broad on-behalf",
			scope: AccessScope{
				"r#^/admin/.*$":                 "on-behalf",
				"r#^/admin/internal/secret/.+$": false,
			},
			args: args{path: "/admin/internal/secret/42"},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Regex ordering: three regexes, middle one wins",
			scope: AccessScope{
				"r#^/a/.*$":                    true,
				"r#^/a/b/.*$":                  false,
				"r#^/a/b/c/very/specific/.*$":  "on-behalf",
			},
			args: args{path: "/a/b/other"},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Regex ordering: three regexes, most specific wins",
			scope: AccessScope{
				"r#^/a/.*$":                    true,
				"r#^/a/b/.*$":                  false,
				"r#^/a/b/c/very/specific/.*$":  "on-behalf",
			},
			args: args{path: "/a/b/c/very/specific/thing"},
			want: constants.ScopeAccessibilityOnBehalf,
		},
		{
			name: "Regex ordering: three regexes, least specific wins when others miss",
			scope: AccessScope{
				"r#^/a/.*$":                    true,
				"r#^/a/b/.*$":                  false,
				"r#^/a/b/c/very/specific/.*$":  "on-behalf",
			},
			args: args{path: "/a/something"},
			want: constants.ScopeAccessibilityAccessible,
		},
		{
			name: "Regex ordering: invalid regex is skipped, valid specific regex wins",
			scope: AccessScope{
				"r#^/broken/(unterminated": true,
				"r#^/ok/.*$":               true,
			},
			args: args{path: "/ok/something"},
			want: constants.ScopeAccessibilityAccessible,
		},
		{
			name: "Regex ordering: invalid regex is skipped, no fallthrough to broken entry",
			scope: AccessScope{
				"r#^/broken/(unterminated": true,
			},
			args: args{path: "/broken/unterminated"},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Regex ordering: deterministic on length tie",
			scope: AccessScope{
				"r#^/tied/a$": true,
				"r#^/tied/b$": false,
			},
			args: args{path: "/tied/a"},
			want: constants.ScopeAccessibilityAccessible,
		},
		{
			name: "Regex ordering: deterministic on length tie, second path",
			scope: AccessScope{
				"r#^/tied/a$": true,
				"r#^/tied/b$": false,
			},
			args: args{path: "/tied/b"},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Regex ordering: nested scope with specificity at leaf level",
			scope: AccessScope{
				"api": AccessScope{
					"r#^/users/.*$":         true,
					"r#^/users/admin/.*$":   false,
				},
			},
			args: args{path: "api|/users/admin/delete"},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Regex ordering: nested scope with specificity at leaf level, allow case",
			scope: AccessScope{
				"api": AccessScope{
					"r#^/users/.*$":         true,
					"r#^/users/admin/.*$":   false,
				},
			},
			args: args{path: "api|/users/42"},
			want: constants.ScopeAccessibilityAccessible,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scope.GetAccessibility(tt.args.path, ""); got != tt.want {
				t.Errorf("AccessScope.GetAccessibility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortedRegexScopeKeys(t *testing.T) {
	tests := []struct {
		name  string
		scope map[string]any
		want  []string
	}{
		{
			name:  "Empty map",
			scope: map[string]any{},
			want:  nil,
		},
		{
			name:  "No regex keys",
			scope: map[string]any{"/a": true, "/b": true},
			want:  nil,
		},
		{
			name: "Mixed keys — only regex keys are returned",
			scope: map[string]any{
				"/literal":     true,
				"r#^/a/.*$":    true,
				"another-lit":  false,
				"r#^/longer/.$": true,
			},
			want: []string{"r#^/longer/.$", "r#^/a/.*$"},
		},
		{
			name: "Length-descending ordering",
			scope: map[string]any{
				"r#^/short$":                          true,
				"r#^/much/much/longer/pattern/here$":  true,
				"r#^/medium/length$":                  true,
			},
			want: []string{
				"r#^/much/much/longer/pattern/here$",
				"r#^/medium/length$",
				"r#^/short$",
			},
		},
		{
			name: "Length ties are broken lexicographically (deterministic)",
			scope: map[string]any{
				"r#^/a$": true,
				"r#^/c$": true,
				"r#^/b$": true,
			},
			want: []string{"r#^/a$", "r#^/b$", "r#^/c$"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortedRegexScopeKeys(tt.scope)
			if len(got) != len(tt.want) {
				t.Fatalf("sortedRegexScopeKeys() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("sortedRegexScopeKeys()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestGetCompiledScopeRegex_CachesAndHandlesInvalid(t *testing.T) {
	// Use unique patterns so other tests don't pollute this one.
	validPattern := "^/cache-test/valid/[0-9]+$"
	invalidPattern := "^/cache-test/invalid/(unterminated"

	// Clean up after ourselves so repeat runs are hermetic.
	defer regexCache.Delete(validPattern)
	defer regexCache.Delete(invalidPattern)

	// First call compiles and caches.
	re1 := getCompiledScopeRegex(validPattern)
	if re1 == nil {
		t.Fatalf("getCompiledScopeRegex(%q) = nil, want compiled regex", validPattern)
	}
	if !re1.MatchString("/cache-test/valid/42") {
		t.Errorf("compiled regex does not match expected input")
	}

	// Second call returns the SAME pointer (cache hit).
	re2 := getCompiledScopeRegex(validPattern)
	if re1 != re2 {
		t.Errorf("getCompiledScopeRegex(%q) returned different pointers on subsequent calls; cache is not working", validPattern)
	}

	// Invalid pattern returns nil and caches the failure.
	if got := getCompiledScopeRegex(invalidPattern); got != nil {
		t.Errorf("getCompiledScopeRegex(%q) = %v, want nil", invalidPattern, got)
	}
	if got := getCompiledScopeRegex(invalidPattern); got != nil {
		t.Errorf("getCompiledScopeRegex(%q) second call = %v, want nil (cached)", invalidPattern, got)
	}
	// The invalid entry should be present in the cache as a nil value.
	if cached, ok := regexCache.Load(invalidPattern); !ok {
		t.Errorf("invalid pattern was not cached")
	} else if cached != nil {
		// We explicitly store a typed nil, which survives as an interface
		// holding (*regexp.Regexp)(nil); it is not == nil in a generic sense,
		// but the value read out of it must itself be a nil pointer.
		if rp, ok := cached.(*regexp.Regexp); !ok || rp != nil {
			t.Errorf("invalid pattern cache entry = %#v, want typed nil *regexp.Regexp", cached)
		}
	}
}

func TestFUPScope_GetLimit_HasLimit(t *testing.T) {
	type args struct {
		path string
	}
	testValue := 123
	var tests = []struct {
		name  string
		scope FUPScope
		args  args
		want  *int
		has   bool
	}{
		{
			name:  "Empty scope",
			scope: FUPScope{},
			args:  args{path: "/"},
			want:  nil,
			has:   false,
		},
		{
			name:  "Empty scope, empty path",
			scope: FUPScope{},
			args:  args{path: ""},
			want:  nil,
			has:   false,
		},
		{
			name:  "Empty scope, path is not in scope",
			scope: FUPScope{},
			args:  args{path: "/"},
			want:  nil,
			has:   false,
		},
		{
			name:  "Path is in scope, period is not, root",
			scope: FUPScope{"/": map[string]any{"hourly": 123}},
			args:  args{path: "/.minutely"},
			want:  nil,
			has:   true,
		},
		{
			name:  "Path is in scope, period is in scope, root",
			scope: FUPScope{"/": map[string]any{"hourly": 123}},
			args:  args{path: "/.hourly"},
			want:  &testValue,
			has:   true,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root",
			scope: FUPScope{"/test": map[string]any{"hourly": 123}},
			args:  args{path: "/test.hourly"},
			want:  &testValue,
			has:   true,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root, multiple keys",
			scope: FUPScope{"/test": map[string]any{"hourly": 123, "minutely": 321}},
			args:  args{path: "/test.hourly"},
			want:  &testValue,
			has:   true,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root, multiple keys, wrong key",
			scope: FUPScope{"/test": map[string]any{"hourly": 123, "minutely": 321}},
			args:  args{path: "/test.daily"},
			want:  nil,
			has:   true,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root, multiple keys, with regex",
			scope: FUPScope{"r#^/test(ing)?x$": map[string]any{"hourly": 123, "minutely": 321}},
			args:  args{path: "/testx.hourly"},
			want:  &testValue,
			has:   true,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root, multiple keys, with regex, different key",
			scope: FUPScope{"r#^/test(ing)?x$": map[string]any{"hourly": 123, "minutely": 321}},
			args:  args{path: "/testingx.hourly"},
			want:  &testValue,
			has:   true,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root, multiple keys, with regex, wrong key",
			scope: FUPScope{"r#^/test(ing)?x$": map[string]any{"hourly": 123, "minutely": 321}},
			args:  args{path: "/testing.hourly"},
			want:  nil,
			has:   false,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root, multiple keys, with regex, wrong suffix key",
			scope: FUPScope{"r#^/test(ing)?x$": map[string]any{"hourly": 123, "minutely": 321}},
			args:  args{path: "/testx1.hourly"},
			want:  nil,
			has:   false,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root, multiple keys, nested",
			scope: FUPScope{"/test": map[string]any{"priority": map[string]any{"hourly": 123, "minutely": 321}}},
			args:  args{path: "/test.priority.hourly"},
			want:  &testValue,
			has:   true,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root, multiple keys, nested, wrong key",
			scope: FUPScope{"/test": map[string]any{"priority": map[string]any{"hourly": 123, "daily": 321}}},
			args:  args{path: "/test.priority.minutely"},
			want:  nil,
			has:   true,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root, multiple keys, nested, with regex",
			scope: FUPScope{"/test": map[string]any{"r#^(priority|prioritized)": map[string]any{"hourly": 123, "minutely": 321}}},
			args:  args{path: "/test.priority.hourly"},
			want:  &testValue,
			has:   true,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root, multiple keys, nested, with regex, different key",
			scope: FUPScope{"/test": map[string]any{"r#^(priority|prioritized)": map[string]any{"hourly": 123, "minutely": 321}}},
			args:  args{path: "/test.prioritized.hourly"},
			want:  &testValue,
			has:   true,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root, multiple keys, nested, with regex, longer key",
			scope: FUPScope{"/test": map[string]any{"r#^(priority|prioritized)": map[string]any{"hourly": 123, "minutely": 321}}},
			args:  args{path: "/test.priority-any.hourly"},
			want:  &testValue,
			has:   true,
		},
		{
			name:  "Path is in scope, period is in scope, path is not root, multiple keys, nested, with regex, wrong key",
			scope: FUPScope{"/test": map[string]any{"r#^(priority|prioritized)": map[string]any{"hourly": 123, "minutely": 321}}},
			args:  args{path: "/test.prio.hourly"},
			want:  nil,
			has:   false,
		},
		// FUPScope specificity ordering: more specific regex wins over broader one.
		{
			name: "FUP regex ordering: specific regex wins over broad",
			scope: FUPScope{
				"r#^/api/v[0-9]+$":          map[string]any{"hourly": 321},
				"r#^/api/v[0-9]+/premium$":  map[string]any{"hourly": 123},
			},
			args: args{path: "/api/v1/premium.hourly"},
			want: &testValue,
			has:  true,
		},
		{
			name: "FUP regex ordering: broad regex used when specific does not match",
			scope: FUPScope{
				"r#^/api/v[0-9]+$":         map[string]any{"hourly": 123},
				"r#^/api/v[0-9]+/premium$": map[string]any{"hourly": 321},
			},
			args: args{path: "/api/v1.hourly"},
			want: &testValue,
			has:  true,
		},
		{
			name: "FUP regex ordering: exact match wins over more specific regex",
			scope: FUPScope{
				"/api/v1/premium":          map[string]any{"hourly": 123},
				"r#^/api/v[0-9]+/premium$": map[string]any{"hourly": 999},
			},
			args: args{path: "/api/v1/premium.hourly"},
			want: &testValue,
			has:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scope.GetLimit(tt.args.path); got != tt.want && (got == nil || *got != *tt.want) {
				t.Errorf("FUPScope.GetLimit() = %v, want %v", got, tt.want)
			}
			pathWithoutPeriod := regexp.MustCompile(`^(.*)\.(minute|hour|dai|week|month)ly$`).ReplaceAllString(tt.args.path, "$1")
			if got := tt.scope.HasLimit(pathWithoutPeriod); got != tt.has {
				t.Errorf("FUPScope.HasLimit() = %v, want %v", got, tt.has)
			}
		})
	}
}
