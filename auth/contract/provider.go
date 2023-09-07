package contract

type ApiClientProviderInterface[T ApiClientInterface] interface {
	ProvideByIdAndSecret(id string, secret string) (ApiClientInterface, *AuthError)
	ProvideByApiKey(apiKey string) (ApiClientInterface, *AuthError)
}
type ApiUserProviderInterface[T ApiUserInterface] interface {
	ProvideByLoginAndPassword(login string, password string) (*T, *AuthError)
	ProvideByToken(token string) (*T, *AuthError)
	ProvideByApiKey(apiKey string) (*T, *AuthError)
}
type ApiUserTokenProviderInterface[T ApiUserTokenInterface] interface {
	// TODO: add methods
	// TODO: do we need this?
}
