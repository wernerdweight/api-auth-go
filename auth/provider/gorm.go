package provider

import (
	"errors"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"github.com/wernerdweight/api-auth-go/auth/entity"
	"gorm.io/gorm"
	"time"
)

// GormApiClientProvider is an implementation of the ApiClientProviderInterface for GORM
type GormApiClientProvider struct {
	newApiClient  func() contract.ApiClientInterface
	getConnection func() *gorm.DB
}

func (p GormApiClientProvider) ProvideByIdAndSecret(id string, secret string) (contract.ApiClientInterface, *contract.AuthError) {
	apiClient := p.newApiClient()
	conn := p.getConnection()
	result := conn.First(&apiClient, entity.GormApiClient{
		ClientId:     id,
		ClientSecret: secret,
	})
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, contract.NewAuthError(contract.ClientNotFound, nil)
		}
		return nil, contract.NewAuthError(contract.DatabaseError, map[string]error{"details": result.Error})
	}
	return apiClient, nil
}

func (p GormApiClientProvider) ProvideByApiKey(apiKey string) (contract.ApiClientInterface, *contract.AuthError) {
	apiClient := p.newApiClient()
	conn := p.getConnection()
	result := conn.First(&apiClient, entity.GormApiClient{
		ApiKey: apiKey,
	})
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, contract.NewAuthError(contract.ClientNotFound, nil)
		}
		return nil, contract.NewAuthError(contract.DatabaseError, map[string]error{"details": result.Error})
	}
	return apiClient, nil
}

func (p GormApiClientProvider) Save(client contract.ApiClientInterface) *contract.AuthError {
	conn := p.getConnection()
	result := conn.Save(client)
	if nil != result.Error {
		return contract.NewAuthError(contract.DatabaseError, map[string]error{"details": result.Error})
	}
	return nil
}

func NewGormApiClientProvider(newApiClient func() contract.ApiClientInterface, getConnection func() *gorm.DB) *GormApiClientProvider {
	return &GormApiClientProvider{
		newApiClient:  newApiClient,
		getConnection: getConnection,
	}
}

// GormApiUserProvider is an implementation of the ApiUserProviderInterface for GORM
type GormApiUserProvider struct {
	newApiUser      func() contract.ApiUserInterface
	newApiUserToken func() contract.ApiUserTokenInterface
	getConnection   func() *gorm.DB
}

func (p GormApiUserProvider) ProvideByLoginAndPassword(login string, password string) (contract.ApiUserInterface, *contract.AuthError) {
	apiUser := p.newApiUser()
	conn := p.getConnection()
	result := conn.First(&apiUser, entity.GormApiUser{
		Login:    login,
		Password: password,
	})
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, contract.NewAuthError(contract.UserNotFound, nil)
		}
		return nil, contract.NewAuthError(contract.DatabaseError, map[string]error{"details": result.Error})
	}
	return apiUser, nil
}

func (p GormApiUserProvider) ProvideByToken(token string) (contract.ApiUserInterface, *contract.AuthError) {
	apiUserToken := p.newApiUserToken()
	conn := p.getConnection()
	result := conn.First(&apiUserToken, entity.GormApiUserToken{
		Token: token,
	})
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, contract.NewAuthError(contract.UserNotFound, nil)
		}
		return nil, contract.NewAuthError(contract.DatabaseError, map[string]error{"details": result.Error})
	}
	if apiUserToken.GetExpirationDate().After(time.Now()) {
		return nil, contract.NewAuthError(contract.UserTokenRequired, map[string]time.Time{"expiredAt": apiUserToken.GetExpirationDate()})
	}
	return apiUserToken.GetApiUser(), nil
}

func (p GormApiUserProvider) Save(user contract.ApiUserInterface) *contract.AuthError {
	conn := p.getConnection()
	result := conn.Save(user)
	if nil != result.Error {
		return contract.NewAuthError(contract.DatabaseError, map[string]error{"details": result.Error})
	}
	return nil
}

func NewGormApiUserProvider(newApiUser func() contract.ApiUserInterface, newApiUserToken func() contract.ApiUserTokenInterface, getConnection func() *gorm.DB) *GormApiUserProvider {
	return &GormApiUserProvider{
		newApiUser:      nil,
		newApiUserToken: nil,
		getConnection:   nil,
	}
}
