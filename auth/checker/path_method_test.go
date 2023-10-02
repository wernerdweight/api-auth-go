package checker

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"net/http"
	"net/url"
	"testing"
)

func TestPathAndMethodAccessScopeChecker_Check(t *testing.T) {
	type args struct {
		scope *contract.AccessScope
		c     *gin.Context
	}
	tests := []struct {
		name string
		ch   PathAndMethodAccessScopeChecker
		args args
		want constants.ScopeAccessibility
	}{
		{
			name: "Nil scope",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: nil, c: nil},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: nil},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope and context",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope, context with empty request",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{Request: &http.Request{}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope, context with request with empty URL",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{Request: &http.Request{URL: nil}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope, context with request with empty URL path",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{Request: &http.Request{URL: &url.URL{}}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope, context with request with empty URL path and method",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{Request: &http.Request{URL: &url.URL{}, Method: ""}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope, context with request with an URL and empty method",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: ""}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Empty scope, context with request with an URL and method",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "GET"}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Context with request with an URL and method, scope with different path and method",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{"post:/other": true}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "GET"}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Context with request with an URL and method, scope with different path and correct method",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{"get:/other": true}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "GET"}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Context with request with an URL and method, scope with correct path and different method",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{"post:/path": true}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "GET"}}},
			want: constants.ScopeAccessibilityForbidden,
		},
		{
			name: "Context with request with an URL and method, scope with correct path and method",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{"get:/path": true}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "GET"}}},
			want: constants.ScopeAccessibilityAccessible,
		},
		{
			name: "Context with request with an URL and method, scope with correct path and method, but different case",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{"get:/path": true}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "GET"}}},
			want: constants.ScopeAccessibilityAccessible,
		},
		{
			name: "Context with request with an URL and method, scope with correct path and method, but different case",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{"get:/path": true}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "get"}}},
			want: constants.ScopeAccessibilityAccessible,
		},
		{
			name: "Context with request with an URL and method, scope with correct path and method, but different case",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{"get:/path": true}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "Get"}}},
			want: constants.ScopeAccessibilityAccessible,
		},
		{
			name: "Context with request with an URL and method, scope with correct path and method, on-behalf",
			ch:   PathAndMethodAccessScopeChecker{},
			args: args{scope: &contract.AccessScope{"get:/path": "on-behalf"}, c: &gin.Context{Request: &http.Request{URL: &url.URL{Path: "/path"}, Method: "GET"}}},
			want: constants.ScopeAccessibilityOnBehalf,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ch.Check(tt.args.scope, tt.args.c); got != tt.want {
				t.Errorf("PathAndMethodAccessScopeChecker.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
