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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scope.GetAccessibility(tt.args.path, ""); got != tt.want {
				t.Errorf("AccessScope.GetAccessibility() = %v, want %v", got, tt.want)
			}
		})
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
