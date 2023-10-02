package checker

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"net/http"
	"net/url"
	"testing"
)

func TestPathAccessScopeChecker_Check(t *testing.T) {
	type args struct {
		scope *contract.AccessScope
		c     *gin.Context
	}
	tests := []struct {
		name string
		ch   PathAccessScopeChecker
		args args
		want constants.ScopeAccessibility
	}{
		{
			name: "Nil scope",
			ch:   PathAccessScopeChecker{},
			args: args{scope: nil, c: nil},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope",
			ch:   PathAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: nil},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope and context",
			ch:   PathAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope, context with empty request",
			ch:   PathAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{Request: &http.Request{}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope, context with request with empty URL",
			ch:   PathAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{Request: &http.Request{URL: nil}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope, context with request with empty URL path",
			ch:   PathAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{Request: &http.Request{URL: &url.URL{}}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope, context with request with empty URL path",
			ch:   PathAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{Request: &http.Request{URL: &url.URL{}, Method: ""}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope, context with request with an URL",
			ch:   PathAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: ""}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope, context with request with an URL and method",
			ch:   PathAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "GET"}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Context with request with an URL and method, scope with different path and method",
			ch:   PathAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{"/other": true}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "GET"}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Context with request with an URL and method, scope with correct path, GET",
			ch:   PathAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{"/path": true}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "GET"}}},
			want: constants.ScopeAccessibilityAccessible,
		},
		{
			name: "Context with request with an URL and method, scope with correct path, POST",
			ch:   PathAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{"/path": true}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "POST"}}},
			want: constants.ScopeAccessibilityAccessible,
		},
		{
			name: "Context with request with an URL and method, scope with correct path, on-behalf",
			ch:   PathAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{"/path": "on-behalf"}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "GET"}}},
			want: constants.ScopeAccessibilityOnBehalf,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ch.Check(tt.args.scope, tt.args.c); got != tt.want {
				t.Errorf("PathAccessScopeChecker.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
