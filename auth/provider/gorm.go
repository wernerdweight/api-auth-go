package provider

import (
	"errors"
	"github.com/google/uuid"
	"github.com/wernerdweight/api-auth-go/v2/auth/config"
	"github.com/wernerdweight/api-auth-go/v2/auth/constants"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
	"github.com/wernerdweight/api-auth-go/v2/auth/encoder"
	"github.com/wernerdweight/api-auth-go/v2/auth/entity"
	generator "github.com/wernerdweight/token-generator-go"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

// GormApiClientProvider is an implementation of the ApiClientProviderInterface for GORM
type GormApiClientProvider struct {
	newApiClient    func() contract.ApiClientInterface
	newApiClientKey func() contract.ApiClientKeyInterface
	getConnection   func() *gorm.DB
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
		return nil, contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	return apiClient, nil
}

func (p GormApiClientProvider) provideByAdditionalKey(apiKey string) (contract.ApiClientInterface, *contract.AuthError) {
	conn := p.getConnection()
	apiClientKey := p.newApiClientKey()
	result := conn.Joins("ApiClient").First(&apiClientKey, entity.GormApiClientKey{
		Key: apiKey,
	})
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, contract.NewAuthError(contract.ClientNotFound, nil)
		}
		return nil, contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}

	if apiClientKey.GetExpirationDate() != nil && apiClientKey.GetExpirationDate().Before(time.Now()) {
		return nil, contract.NewAuthError(contract.ApiKeyExpired, map[string]time.Time{"expiredAt": *apiClientKey.GetExpirationDate()})
	}

	// ApiClient needs to be fetched separately to return user defined model (otherwise it would be GormApiClient)
	apiClient := p.newApiClient()
	result = conn.First(&apiClient, apiClientKey.GetApiClient().(*entity.GormApiClient).ID)
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, contract.NewAuthError(contract.ClientNotFound, nil)
		}
		return nil, contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	apiClient.SetCurrentApiKey(apiClientKey)
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
			if config.ProviderInstance.IsAdditionalApiKeysEnabled() {
				return p.provideByAdditionalKey(apiKey)
			}
			return nil, contract.NewAuthError(contract.ClientNotFound, nil)
		}
		return nil, contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	return apiClient, nil
}

func (p GormApiClientProvider) Save(client contract.ApiClientInterface) *contract.AuthError {
	conn := p.getConnection()
	result := conn.Save(client)
	if nil != result.Error {
		return contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	return nil
}

func NewGormApiClientProvider(newApiClient func() contract.ApiClientInterface, newApiClientKey func() contract.ApiClientKeyInterface, getConnection func() *gorm.DB) *GormApiClientProvider {
	return &GormApiClientProvider{
		newApiClient:    newApiClient,
		newApiClientKey: newApiClientKey,
		getConnection:   getConnection,
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
	result := conn.First(&apiUser, "lower(email) = lower(?)", login)
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, contract.NewAuthError(contract.UserNotFound, nil)
		}
		return nil, contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
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
		return nil, contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	if apiUserToken.GetExpirationDate().Before(time.Now()) {
		return nil, contract.NewAuthError(contract.UserTokenExpired, map[string]time.Time{"expiredAt": apiUserToken.GetExpirationDate()})
	}
	if !apiUserToken.GetApiUser().IsActive() {
		return nil, contract.NewAuthError(contract.UserNotActive, nil)
	}
	// ApiUser needs to be fetched separately to return user defined model (otherwise it would be GormApiUser)
	apiUser := p.newApiUser()
	result = conn.First(&apiUser, apiUserToken.GetApiUser().(*entity.GormApiUser).ID)
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, contract.NewAuthError(contract.UserNotFound, nil)
		}
		return nil, contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	return apiUser, nil
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
		return nil, contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
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
		return nil, contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
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

func (p GormApiUserProvider) InvalidateTokens(user contract.ApiUserInterface) *contract.AuthError {
	slog.Debug("invalidating tokens for user", slog.String("user", user.GetLogin()))
	conn := p.getConnection()
	id, err := uuid.Parse(user.GetID())
	if nil != err {
		return contract.NewInternalError(contract.DatabaseError, map[string]string{"details": err.Error()})
	}
	var tokens []entity.GormApiUserToken
	result := conn.Where(&entity.GormApiUserToken{ApiUserID: id}).Where("expiration_date >= ?", time.Now()).Find(&tokens)
	if nil != result.Error {
		return contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	if 0 == len(tokens) {
		slog.Debug("no tokens to invalidate for user", slog.String("user", user.GetLogin()))
		return nil
	}
	slog.Debug("tokens to invalidate for user", slog.Int("tokens", len(tokens)), slog.String("user", user.GetLogin()))
	for index, token := range tokens {
		tokens[index].SetExpirationDate(time.Now())
		slog.Debug("invalidating token", slog.String("token", token.Token))
		if config.ProviderInstance.IsCacheEnabled() {
			cacheErr := config.ProviderInstance.GetCacheDriver().InvalidateToken(token.Token)
			if nil != cacheErr {
				slog.Error("can't invalidate token in cache", slog.String("token", token.Token), slog.String("user", user.GetLogin()), slog.String("error", cacheErr.Err.Error()))
			}
		}
	}
	slog.Debug("saving invalidated tokens for user", slog.String("user", user.GetLogin()), slog.Int("tokens", len(tokens)))
	result = conn.Save(&tokens)
	if nil != result.Error {
		return contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
	}
	slog.Debug("tokens invalidated")
	return nil
}

func (p GormApiUserProvider) Save(user contract.ApiUserInterface) *contract.AuthError {
	conn := p.getConnection()
	result := conn.Save(user)
	if nil != result.Error {
		return contract.NewInternalError(contract.DatabaseError, map[string]string{"details": result.Error.Error()})
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
