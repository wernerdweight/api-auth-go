package contract

type ApiClientProviderInterface[T ApiClientInterface] interface {
	ProvideByIdAndSecret(id string, secret string) (ApiClientInterface, *AuthError)
	ProvideByApiKey(apiKey string) (ApiClientInterface, *AuthError)
	Save(client ApiClientInterface) *AuthError
}
type ApiUserProviderInterface[T ApiUserInterface] interface {
	ProvideByLoginAndPassword(login string, password string) (ApiUserInterface, *AuthError)
	ProvideByLogin(login string) (ApiUserInterface, *AuthError)
	ProvideByToken(token string) (ApiUserInterface, *AuthError)
	ProvideByConfirmationToken(token string) (ApiUserInterface, *AuthError)
	ProvideByResetToken(token string) (ApiUserInterface, *AuthError)
	ProvideNew(login string, encryptedPassword string) ApiUserInterface
	Save(user ApiUserInterface) *AuthError
}
