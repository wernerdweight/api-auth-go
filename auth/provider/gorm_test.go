package provider

import (
	"bytes"
	"errors"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/wernerdweight/api-auth-go/v3/auth/config"
	"github.com/wernerdweight/api-auth-go/v3/auth/contract"
	"github.com/wernerdweight/api-auth-go/v3/auth/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type failingCacheDriver struct {
	contract.CacheDriverInterface
}

func (d failingCacheDriver) InvalidateToken(token string) *contract.AuthError {
	return &contract.AuthError{Err: errors.New("cache unavailable")}
}

func TestGormApiUserProvider_InvalidateTokens_DoesNotLogTokensOrLogin(t *testing.T) {
	var logs bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(previous) })

	cacheTTL := time.Hour
	config.ProviderInstance.Init(contract.Config{Cache: &contract.CacheConfig{Driver: failingCacheDriver{}, TTL: &cacheTTL}})

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	conn, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	userID := uuid.New()
	login := "reset.user@example.com"
	tokenValue := "secret-token-value-1234567890"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "api_user_token"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "token", "expiration_date", "api_user_id"}).
			AddRow(uuid.New(), tokenValue, time.Now().Add(time.Hour), userID))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "api_user_token"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(uuid.New(), time.Now()))
	mock.ExpectCommit()

	user := &entity.GormApiUser{ID: userID, Login: login}
	provider := NewGormApiUserProvider(
		func() contract.ApiUserInterface { return &entity.GormApiUser{} },
		func() contract.ApiUserTokenInterface { return &entity.GormApiUserToken{} },
		func() *gorm.DB { return conn },
	)

	authErr := provider.InvalidateTokens(user)

	assert.Nil(t, authErr)
	assert.NoError(t, mock.ExpectationsWereMet())
	output := logs.String()
	assert.Contains(t, output, "can't invalidate token in cache")
	assert.Contains(t, output, `"user_id":"`+userID.String()+`"`)
	assert.NotContains(t, output, tokenValue)
	assert.NotContains(t, output, login)
}
