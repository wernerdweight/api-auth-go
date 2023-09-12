package contract

type ApiClientProviderInterface[T ApiClientInterface] interface {
	ProvideByIdAndSecret(id string, secret string) (ApiClientInterface, *AuthError)
	ProvideByApiKey(apiKey string) (ApiClientInterface, *AuthError)
	Save(client ApiClientInterface) *AuthError
}
type ApiUserProviderInterface[T ApiUserInterface] interface {
	ProvideByLoginAndPassword(login string, password string) (ApiUserInterface, *AuthError)
	ProvideByToken(token string) (ApiUserInterface, *AuthError)
	Save(user ApiUserInterface) *AuthError
}
