package provider

import (
	"errors"
	"github.com/wernerdweight/api-auth-go/auth/config"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"github.com/wernerdweight/api-auth-go/auth/encoder"
	"github.com/wernerdweight/api-auth-go/auth/entity"
	generator "github.com/wernerdweight/token-generator-go"
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
		return nil, contract.NewAuthError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
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
		return nil, contract.NewAuthError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	return apiClient, nil
}

func (p GormApiClientProvider) Save(client contract.ApiClientInterface) *contract.AuthError {
	conn := p.getConnection()
	result := conn.Save(client)
	if nil != result.Error {
		return contract.NewAuthError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
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
	apiUser, err := p.ProvideByLogin(login)
	if nil != err {
		return nil, err
	}
	err = encoder.ComparePassword(apiUser, password)
	if nil != err {
		return nil, err
	}
	if !apiUser.IsActive() {
		return nil, contract.NewAuthError(contract.UserNotActive, nil)
	}
	return apiUser, nil
}

func (p GormApiUserProvider) ProvideByLogin(login string) (contract.ApiUserInterface, *contract.AuthError) {
	apiUser := p.newApiUser()
	conn := p.getConnection()
	result := conn.First(&apiUser, entity.GormApiUser{
		Login: login,
	})
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, contract.NewAuthError(contract.UserNotFound, nil)
		}
		return nil, contract.NewAuthError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	return apiUser, nil
}

func (p GormApiUserProvider) ProvideByToken(token string) (contract.ApiUserInterface, *contract.AuthError) {
	apiUserToken := p.newApiUserToken()
	conn := p.getConnection()
	result := conn.Joins("ApiUser").First(&apiUserToken, entity.GormApiUserToken{
		Token: token,
	})
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, contract.NewAuthError(contract.UserTokenNotFound, nil)
		}
		return nil, contract.NewAuthError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	if apiUserToken.GetExpirationDate().Before(time.Now()) {
		return nil, contract.NewAuthError(contract.UserTokenExpired, map[string]time.Time{"expiredAt": apiUserToken.GetExpirationDate()})
	}
	return apiUserToken.GetApiUser(), nil
}

func (p GormApiUserProvider) ProvideByConfirmationToken(token string) (contract.ApiUserInterface, *contract.AuthError) {
	apiUser := p.newApiUser()
	conn := p.getConnection()
	result := conn.First(&apiUser, entity.GormApiUser{
		ConfirmationToken: &token,
	})
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, contract.NewAuthError(contract.UserNotFound, nil)
		}
		return nil, contract.NewAuthError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	// check token expiration
	expirationInterval := config.ProviderInstance.GetConfirmationTokenExpirationInterval()
	expiresAt := apiUser.GetConfirmationRequestedAt().Add(expirationInterval)
	if expiresAt.Before(time.Now()) {
		return nil, contract.NewAuthError(contract.ConfirmationTokenExpired, map[string]time.Time{"expiredAt": expiresAt})
	}
	return apiUser, nil
}

func (p GormApiUserProvider) ProvideByResetToken(token string) (contract.ApiUserInterface, *contract.AuthError) {
	apiUser := p.newApiUser()
	conn := p.getConnection()
	result := conn.First(&apiUser, entity.GormApiUser{
		ResetToken: &token,
	})
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, contract.NewAuthError(contract.UserNotFound, nil)
		}
		return nil, contract.NewAuthError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	// check token expiration
	expirationInterval := config.ProviderInstance.GetConfirmationTokenExpirationInterval()
	expiresAt := apiUser.GetResetRequestedAt().Add(expirationInterval)
	if expiresAt.Before(time.Now()) {
		return nil, contract.NewAuthError(contract.ResetTokenExpired, map[string]time.Time{"expiredAt": expiresAt})
	}
	return apiUser, nil
}

func (p GormApiUserProvider) ProvideNew(login string, encryptedPassword string) contract.ApiUserInterface {
	token := generator.NewTokenGenerator("").Generate(constants.DefaultTokenLength)
	now := time.Now()
	apiUser := p.newApiUser()
	apiUser.SetLogin(login)
	apiUser.SetPassword(encryptedPassword)
	apiUser.SetConfirmationToken(&token)
	apiUser.SetConfirmationRequestedAt(&now)
	return apiUser
}

func (p GormApiUserProvider) Save(user contract.ApiUserInterface) *contract.AuthError {
	conn := p.getConnection()
	result := conn.Save(user)
	if nil != result.Error {
		return contract.NewAuthError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	return nil
}

func NewGormApiUserProvider(newApiUser func() contract.ApiUserInterface, newApiUserToken func() contract.ApiUserTokenInterface, getConnection func() *gorm.DB) *GormApiUserProvider {
	return &GormApiUserProvider{
		newApiUser:      newApiUser,
		newApiUserToken: newApiUserToken,
		getConnection:   getConnection,
	}
}
