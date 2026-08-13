package auth

import (
	"testing"

	"github.com/wernerdweight/api-auth-go/v3/auth/config"
	"github.com/wernerdweight/api-auth-go/v3/auth/contract"
	"github.com/wernerdweight/api-auth-go/v3/auth/fup"
)

// TestMiddleware_DefaultsAnonymousFUPChecker covers the default of the anonymous FUP checker,
// which can't live among the config provider's defaults (see Middleware)
func TestMiddleware_DefaultsAnonymousFUPChecker(t *testing.T) {
	scope := contract.FUPScope{}
	tests := []struct {
		name    string
		checker contract.FUPCheckerInterface
		want    contract.FUPCheckerInterface
	}{
		{name: "not configured", want: fup.IPFUPChecker{}},
		{name: "configured", checker: fup.CookieFUPChecker{}, want: fup.CookieFUPChecker{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Middleware(nil, contract.Config{
				Client: contract.ClientConfig{
					AnonymousFUPScope:   &scope,
					AnonymousFUPChecker: tt.checker,
				},
			})
			if got := config.ProviderInstance.GetAnonymousFUPChecker(); got != tt.want {
				t.Errorf("GetAnonymousFUPChecker() = %T, want %T", got, tt.want)
			}
		})
	}
}
