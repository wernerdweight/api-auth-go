package contract

type ApiClientProviderInterface[T ApiClientInterface] interface {
	ProvideByIdAndSecret(id string, secret string) (T, error)
}
type ApiUserProviderInterface[T ApiUserInterface] interface {
	ProvideByLoginAndPassword(login string, password string) (T, error)
	ProvideByToken(token string) (T, error)
}
type ApiUserTokenProviderInterface[T ApiUserTokenInterface] interface {
	// TODO: add methods
	// TODO: do we need this?
}
