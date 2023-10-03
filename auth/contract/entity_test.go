package contract

import (
	"github.com/wernerdweight/api-auth-go/auth/constants"
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
			args:  args{path: "test.nested1"},
			want:  constants.ScopeAccessibilityAccessible,
		},
		{
			name:  "Nested scope, false",
			scope: AccessScope{"test": AccessScope{"nested1": true, "nested2": false, "nested3": "on-behalf"}},
			args:  args{path: "test.nested2"},
			want:  constants.ScopeAccessibilityForbidden,
		},
		{
			name:  "Nested scope, on-behalf",
			scope: AccessScope{"test": AccessScope{"nested1": true, "nested2": false, "nested3": "on-behalf"}},
			args:  args{path: "test.nested3"},
			want:  constants.ScopeAccessibilityOnBehalf,
		},
		{
			name:  "Nested scope, not in scope",
			scope: AccessScope{"test": AccessScope{"nested1": true, "nested2": false, "nested3": "on-behalf"}},
			args:  args{path: "test.nested4"},
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
			args:  args{path: "test.test.test.nope"},
			want:  constants.ScopeAccessibilityForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scope.GetAccessibility(tt.args.path); got != tt.want {
				t.Errorf("AccessScope.GetAccessibility() = %v, want %v", got, tt.want)
			}
		})
	}
}
