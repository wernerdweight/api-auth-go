API Auth for Gin (Go Framework)
====================================

A Gin middleware providing authentication, authorisation, access control, FUP and other functionality for APIs.

[![Build Status](https://www.travis-ci.com/wernerdweight/api-auth-go.svg?branch=master)](https://www.travis-ci.com/wernerdweight/api-auth-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/wernerdweight/api-auth-go)](https://goreportcard.com/report/github.com/wernerdweight/api-auth-go)
[![GoDoc](https://godoc.org/github.com/wernerdweight/api-auth-go?status.svg)](https://godoc.org/github.com/wernerdweight/api-auth-go)
[![go.dev](https://img.shields.io/badge/go.dev-pkg-007d9c.svg?style=flat)](https://pkg.go.dev/github.com/wernerdweight/api-auth-go)


Installation
------------

### 1. Installation

```bash
go get github.com/wernerdweight/api-auth-go
```

Configuration and Usage
------------

Full configuration structure (only client configuration is mandatory, default values are mentioned in comments):

```go
{
    // Client: api client configuration (mandatory)
    Client: {
        // Provider: your provider that implements ApiClientProviderInterface
        Provider ApiClientProviderInterface[ApiClientInterface]
        // UseScopeAccessModel: if set to true, client scope will be checked before granting access (see `scope access` below) - default false
        UseScopeAccessModel *bool
        // AccessScopeChecker: the checker used to check scope access that implements AccessScopeCheckerInterface - defaults to PathAccessScopeChecker
        AccessScopeChecker AccessScopeCheckerInterface
        // FUPChecker: the checker used to check FUP limits that implements FUPCheckerInterface (optional; if you omit FUP checker, FUP limits will not be checked)
        // NOTE: if you want to use FUP limits, you must also enable Cache (see below)
        FUPChecker FUPCheckerInterface
    }
    
    // User: api user configuration (optional; if you omit user configuration, you will not be able to use `on-behalf` access mode (see below))
    User: *{
        // Provider: your provider that implements ApiUserProviderInterface
        Provider ApiUserProviderInterface[ApiUserInterface]
        // TokenFactory: generates your token type that implements ApiUserTokenInterface
        TokenFactory func() ApiUserTokenInterface
        // ApiTokenExpirationInterval: token expiration in seconds - defaults to 2,592,000 (30 days)
        ApiTokenExpirationInterval *time.Duration
        // UseScopeAccessModel: if set to true, user scope will be checked before granting access (see `scope access` below) - default false
        UseScopeAccessModel *bool
        // AccessScopeChecker: the checker used to check scope access that implements AccessScopeCheckerInterface - defaults to PathAccessScopeChecker
        AccessScopeChecker AccessScopeCheckerInterface
        // WithRegistration: if set to true, user registration will be enabled - default false
        WithRegistration *bool
        // ConfirmationTokenExpirationInterval: confirmation token expiration in seconds - defaults to 43200 (12 hours)
        ConfirmationTokenExpirationInterval *time.Duration
        // FUPChecker: the checker used to check FUP limits that implements FUPCheckerInterface (optional; if you omit FUP checker, FUP limits will not be checked)
        // NOTE: if you want to use FUP limits, you must also enable Cache (see below)
        FUPChecker FUPCheckerInterface
    }
    
    // Mode: modes of authentication (client id + secret and user token vs. api key)
    Mode *{
        // ApiKey: api key authentication mode (optional; default false)
        ApiKey *bool
        // ClientIdAndSecret: client id and secret authentication mode (optional; default true)
        ClientIdAndSecret *bool
    }

    // TargetHandlers: list of handlers to target (optional; if you omit target handlers, all handlers will be targeted)
    TargetHandlers *[]string
    // '.*'            	# all handlers
    // '/v1/.*'   		# all handlers starting with '/v1/'
    // '/v1/some/path'  # only '/v1/some/path' handler

	// ExcludeHandlers: list of handlers to exclude (optional; if you omit exclude handlers, no handlers will be excluded)
	ExcludeHandlers *[]string
    
    // ExcludeOptionsRequests: if true, requests using the OPTIONS method will be ignored (authentication will be skipped) - default false
    ExcludeOptionsRequests *bool
    
    // Cache: cache configuration (optional)
    Cache *{
        // Driver: your cache driver that implements CacheDriverInterface
        Driver CacheDriverInterface
        // Prefix: prefix to use for cache keys - defaults to `api-auth-go:`
        Prefix *string
        // TTL: cache TTL in seconds - defaults to 3600 (1 hour)
        TTL *time.Duration
    }
}
```

### Minimal configuration:

Please note that the minimal examples are using in-memory data providers, which are not suitable for production use. You should use your own implementation of data providers (see below), or you can use included GORM data providers (also below).

```go
package main

import "github.com/wernerdweight/api-auth-go/auth/contract"

contract.Config{
    Client: contract.ClientConfig{
        Provider: provider.NewMemoryApiClientProvider([]entity.MemoryApiClient{
            {Id: "id", Secret: "secret", ApiKey: "api-key"},
            {Id: "another-id", Secret: "another-secret", ApiKey: "another-api-key"},
            ...
        }),
    },
}
```

### ApiClient

You need to create a model that implements ApiClientInterface. This package provides an abstract implementation of this for in-memory data provider (which should not be used in production), and an implementation for GORM (see below).

```go
type AccessScope map[string]any
type FUPScope map[string]any

type ApiClientInterface interface {
    GetClientId() string
    GetClientSecret() string
    GetApiKey() string
    GetClientScope() *AccessScope
    GetFUPScope() *FUPScope
}

// if you want to use GORM as data provider, you can extend this type
type ApiClient struct {
    entity.GormApiClient
    // TODO: your other fields here
}

```

### ApiUser

**OPTIONAL:** If you want to restrict certain actions within your API to certain users (see 'on behalf' access mode below), create a model that implements ApiUserInterface and another one that implements ApiUserTokenInterface.

```go
package main

import "github.com/wernerdweight/api-auth-go/auth/contract"

contract.Config{
    Client: contract.ClientConfig{
        Provider: provider.NewMemoryApiClientProvider([]entity.MemoryApiClient{
            {Id: "id", Secret: "secret", ApiKey: "api-key"},
            {Id: "another-id", Secret: "another-secret", ApiKey: "another-api-key"},
            ...
        }),
    },
    User: &contract.UserConfig{
        // only use one of the following, or implement your own provider
        Provider: provider.NewMemoryApiUserProvider([]entity.MemoryApiUser{
            // note: you only need to provide the AccessScope if you want to enable the UseScopeAccessModel below
            {Id: "user-id", Login: "user@domain.tld", Password: "not-so-secret", CurrentToken: &entity.MemoryApiUserToken{Token: "secret-user-token", ExpirationDate: time.Time{}}, AccessScope: &contract.AccessScope{"/v1/some/path": true, "/v1/some/false": false, "/v1/some/true": true}},
        }),
        Provider: provider.NewGormApiUserProvider(newApiUser, newApiUserToken, getDBConnection),
        TokenFactory: func() contract.ApiUserTokenInterface {
            return &YourApiUserTokenImplementation{}
        },
        // you can optionally enable user scope access model analogically to client scope access model
        UseScopeAccessModel: &useUserScopeAccessModel,
    },
}
```

```go
type AccessScope map[string]any
type FUPScope map[string]any

type ApiUserInterface interface {
    AddApiToken(apiToken ApiUserTokenInterface)
    GetCurrentToken() ApiUserTokenInterface
    GetUserScope() *AccessScope
    GetLastLoginAt() *time.Time
    SetLastLoginAt(lastLoginAt *time.Time)
    GetPassword() string
    SetPassword(password string)
    GetLogin() string
    SetLogin(login string)
    SetConfirmationToken(confirmationToken *string)
    GetConfirmationRequestedAt() *time.Time
    SetConfirmationRequestedAt(confirmationRequestedAt *time.Time)
    IsActive() bool
    SetActive(active bool)
    GetResetRequestedAt() *time.Time
    SetResetRequestedAt(resetRequestedAt *time.Time)
    GetResetToken() *string
    SetResetToken(resetToken *string)
    GetFUPScope() *FUPScope
}
type ApiUserTokenInterface interface {
    SetToken(token string)
    GetToken() string
    SetExpirationDate(expirationDate time.Time)
    GetExpirationDate() time.Time
    SetApiUser(apiUser ApiUserInterface)
    GetApiUser() ApiUserInterface
}

// if you want to use GORM as data provider, you can extend these types
type ApiUser struct {
    entity.GormApiUser
    // TODO: your other fields here
}
type ApiUserToken struct {
    entity.GormApiUserToken
    // TODO: your other fields here
}
```

### Authentication:

To authenticate any request, you need to provide the client id and secret in the `X-Client-Id` and `X-Client-Secret` headers of the request.

```http request
POST /some/path HTTP/1.1
X-Client-Id: some-client-id
X-Client-Secret: some-client-secret
Host: your-api-host.com
```

### Only targeting certain handlers:

By default, all handlers/routes will be targeted automatically (no configuration required).

You can target only certain handlers by providing a list of regular expressions for `TargetHandlers` configuration value. These are matched agains the URL path of the request.

```go
package main

import "github.com/wernerdweight/api-auth-go/auth/contract"

contract.Config{
    Client: contract.ClientConfig{
        Provider: provider.NewMemoryApiClientProvider(...),
    },
    TargetHandlers: &[]string{
        "^/v1/.*",
    },
}
```

### API key authentication mode:

By default, client id and secret authentication mode is used. You can enable API key authentication mode by setting `Mode.ApiKey` to `true`.
If set, the middleware will look for the API key in the `Authorization` header of the request.
Both modes can be used at the same time (in that case, client id and secret will be checked first, and if not found, API key will be checked).

```go
package main

import "github.com/wernerdweight/api-auth-go/auth/contract"

useApiKeyMode := true

contract.Config{
    Client: contract.ClientConfig{
        Provider: provider.NewMemoryApiClientProvider(...),
    },
    Mode: &contract.ModeConfig{
        ApiKey: &useApiKeyMode,
    },
}
```

You can then authenticate your requests like this:

```http request
POST /some/path HTTP/1.1
Authorization: your-api-key
Host: your-api-host.com
```

### Using GORM as data provider:

The implementation of GORM data provider is included in this package. You can use it by providing your own implementation of `ApiClient`, `ApiUser` and `ApiUserToken` types (see above), and then providing a function that returns a GORM connection (see below).

```go
package main

import (
    "github.com/wernerdweight/api-auth-go/auth/contract"
    "gorm.io/gorm"
)

getDBConnection := func() *gorm.DB { /* TODO: your implementation */ }
newApiClient := func() contract.ApiClientInterface { return &YoutApiClienImplementation{} }

contract.Config{
    Client: contract.ClientConfig{
        Provider: provider.NewGormApiClientProvider(newApiClient, getDBConnection),
    },
}
```

### With scoped access model:

By default, all clients and users have access to all paths. You can enable scoped access model by setting `UseScopeAccessModel` to `true` in `Client` and/or `User` configuration (see below).

If enabled, the authenticator will also (apart from api credentials) check the defined client/user scope using configured checker (if no checker is explicitly configured, the default `PathChecker` is used). This way, different ApiClients/Users can have different privileges.

```go
package main

import "github.com/wernerdweight/api-auth-go/auth/contract"

useClientScopeAccessModel := true

contract.Config{
    Client: contract.ClientConfig{
        Provider: provider.NewMemoryApiClientProvider([]entity.MemoryApiClient{
            {Id: "id", Secret: "secret", ApiKey: "api-key", AccessScope: &contract.AccessScope{"/v1/some/path": "on-behalf", "/v1/some/false": false, "/v1/some/true": true}},
            {Id: "another-id", Secret: "another-secret", ApiKey: "another-api-key", AccessScope: &contract.AccessScope{"/v1/some/path": "on-behalf"}},
            ...
        }),
        UseScopeAccessModel: &useClientScopeAccessModel,
    },
}
```

The scope is generally a JSON column on ApiClient/ApiUser entities. You can store any information in that column and then use any checker you want to read and evaluate the stored information.

The default PathChecker expects a structure like this:

```json5
{
  "/some/path": true,
  // following line is a no-op, the route doesn't have to be specified if it should not be accessible
  "/some/other/path": false,
  // see `on-behalf` access mode below
  "/yet/another/path": 'on-behalf',
  // regexes are also supported with `r#` prefix
  "r#^/some/regex/[^/]+/?$": true,
}
```

This package also includes a `PathAndMethodChecker`, which also checks based on the HTTP method, and expects this structure:
**Please note that you're supposed to keep the keys lowercase for these built-in checkers.**

```json5
{
  "get:/some/path": true,
  // following line is a no-op, the route doesn't have to be specified if it should not be accessible
  "post:/some/other/path": false,
  // see `on-behalf` access mode below
  "delete:/yet/another/path": 'on-behalf',
  // regexes are also supported with `r#` prefix
  "r#get:^/some/regex/[^/]+/?$": true,
}
```

You can also implement custom checker by implementing `AccessScopeCheckerInterface`.

```go
type AccessScopeCheckerInterface interface {
    Check(scope *AccessScope, c *gin.Context) constants.ScopeAccessibility
}
```


### "on-behalf" access mode

If the ApiClient/ApiUser scope is configured to be checked (see above) and the `'on-behalf'` value is set in the scope, another authentication is required.

The request must then contain the `X-Api-User-Token` header with a valid token. To obtain the token, the user must login using Basic Auth - the request should look as follows:

```http request
POST /authenticate HTTP/1.1
X-Client-Id: some-client-id
X-Client-Secret: some-client-secret
Authorization: Basic encodedBasicAuthInformation==
Host: your-api-host.com
```

The response contains the token and scope (and optionally any other information returned from your user entity via json serialization):

```json
{
  "id": "62d5de93-eccc-45a4-b971-5fb11be0d139",
  "lastLoginAt": "2023-09-26T22:14:45.230747+02:00",
  "token": {
    "expirationDate": "2023-11-01T19:51:27.787135+01:00",
    "token": "aBc37De4FgH_-abC08d7eF",
  },
  "userScope": {
    "/v1/some/false": false,
    "/v1/some/path": true,
    "/v1/some/true": true
  }
}
```

You can then use the obtained token in requests that require the `'ob-behalf'` access mode like this:

```http request
GET /your/api/path HTTP/1.1
X-Client-Id: some-client-id
X-Client-Secret: some-client-secret
X-Api-User-Token: aBc37De4FgH_-abC08d7eF
Host: your-api-host.com
```

**FYI:** The `'on-behalf'` value only makes sense for client scope. If you set `'on-behalf'` as value inside the user scope, the value is interpreted in the same way as `true`.

### With cache:

You can enable caching through one of the built-in cache drivers (memory, Redis) providing your own implementation of `CacheDriverInterface` (see below).
This speeds up the authentication process by caching the results of the authentication process (client id and secret, user token, user scope, etc.).
By default, the cache expires after 1 hour (3600 seconds), but you can change this by setting `Cache.TTL` to a different value.
You can also change the cache prefix by setting `Cache.Prefix`, which is useful if you want to use the same Redis instance for multiple applications or environments.

```go
package main

import "github.com/wernerdweight/api-auth-go/auth/contract"

redisDsn, _ := getenv.GetEnv("REDIS_URL")
newApiClient := func() contract.ApiClientInterface { return &YourApiClientImplementation{} }
newApiUser := func() contract.ApiUserInterface { return &YourApiUserImplementation{} }

contract.Config{
    Client: contract.ClientConfig{
        Provider: provider.NewMemoryApiClientProvider([]entity.MemoryApiClient{
            {Id: "id", Secret: "secret", ApiKey: "api-key"},
            {Id: "another-id", Secret: "another-secret", ApiKey: "another-api-key"},
            ...
        }),
    },
    Cache: &contract.CacheConfig{
        // only use one of the following, or implement your own driver
        Driver: cache.NewMemoryCacheDriver(),
        Driver: cache.NewRedisCacheDriver(redisDsn, newApiClient, newApiUser),
    },
}
```

### With user registration:

By default, user registration is disabled. If you don't already have registration process in place, you can enable built-in registration by setting `WithRegistration` to `true` in `User` configuration (see below).

```go
package main

import "github.com/wernerdweight/api-auth-go/auth/contract"

withRegistration := true

contract.Config{
    Client: contract.ClientConfig{
        Provider: provider.NewMemoryApiClientProvider([]entity.MemoryApiClient{
            {Id: "id", Secret: "secret", ApiKey: "api-key"},
            {Id: "another-id", Secret: "another-secret", ApiKey: "another-api-key"},
            ...
        }),
    },
    User: &contract.UserConfig{
        Provider: provider.NewMemoryApiUserProvider([]entity.MemoryApiUser{
            {Id: "user-id", Login: "user@domain.tld", Password: "not-so-secret", CurrentToken: &entity.MemoryApiUserToken{Token: "secret-user-token", ExpirationDate: time.Time{}}},
        }),
        TokenFactory: func() contract.ApiUserTokenInterface {
            return &YourApiUserTokenImplementation{}
        },
        WithRegistration:    &withRegistration,
    },
}
```

This will enable the following endpoints:

**Registration:** to register a new user, send a POST request to `/registration/request` with the following payload:

```http request
POST /registration/request HTTP/1.1
Content-Type: application/json
Host: your-api-host.com

{
	"email": "user@domain.tld",
	"password": "testPass123"
}
```

After successful registration, the user needs to confirm their e-mail address. You need to provide the functionality to deliver the confirmation link to the user (e.g. via email) and send the confirmation request.
You can subscribe to an event dispatched by this package (see below).
The request to confirm the e-mail address looks like this:

```http request
POST /registration/confirm/{confirmationToken} HTTP/1.1
Host: your-api-host.com
```

Enabling registration also enables the password reset functionality. To request a password reset, send a POST request to `/resetting/request` with the following payload:

```http request
POST /resetting/request HTTP/1.1
Content-Type: application/json
Host: your-api-host.com

{
	"email": "user@domain.tld"
}
```

After successful request, the user needs to confirm their e-mail address. You need to provide the functionality to deliver the confirmation link to the user (e.g. via email) and send the confirmation request with changed password.
You can subscribe to an event dispatched by this package (see below).
The request to reset the password looks like this:

```http request
POST /resetting/reset/{resetToken} HTTP/1.1
Content-Type: application/json
Host: your-api-host.com

{
	"password": "1234TestPass"
}
```

By default, the tokens above are only valid for 12 hours. You can change this by setting `ConfirmationTokenExpirationInterval` to a different value in `User` configuration (see above).
During this interval, you can not request another password reset for the same user.

### With FUP limits:

By default, FUP limits are disabled. If you want to enable FUP limits, you can configure one of the built-in FUP checkers (Path, PathAndMethod), or you can provide your own implementation of `FUPCheckerInterface` (see below). You then need to enable it in `Client` and/or `User` configuration (see below).
Please note that for this functionality to work, you also need to enable cache (see above).

```go
package main

import "github.com/wernerdweight/api-auth-go/auth/contract"

contract.Config{
    Client: contract.ClientConfig{
        Provider: provider.NewMemoryApiClientProvider([]entity.MemoryApiClient{
            {Id: "id", Secret: "secret", ApiKey: "api-key"},
            {Id: "another-id", Secret: "another-secret", ApiKey: "another-api-key"},
            ...
        }),
        // only use one of the following, or implement your own checker
        FUPChecker: fup.PathFUPChecker{},
        FUPChecker: fup.PathAndMethodFUPChecker{},
    },
	// note: user config is still optional, included here for completeness
    User: &contract.UserConfig{
        Provider: provider.NewMemoryApiUserProvider([]entity.MemoryApiUser{
            {Id: "user-id", Login: "user@domain.tld", Password: "not-so-secret", CurrentToken: &entity.MemoryApiUserToken{Token: "secret-user-token", ExpirationDate: time.Time{}}},
        }),
        TokenFactory: func() contract.ApiUserTokenInterface {
            return &YourApiUserTokenImplementation{}
        },
        FUPChecker: fup.PathFUPChecker{},
        FUPChecker: fup.PathAndMethodFUPChecker{},
    },
    Cache: &contract.CacheConfig{
        // only use one of the following, or implement your own driver
        Driver: cache.NewMemoryCacheDriver(),
        Driver: cache.NewRedisCacheDriver(redisDsn, newApiClient, newApiUser),
    },
}
```

The scope is generally another JSON column on ApiClient/ApiUser entities. You can store any information in that column and then use any checker you want to read and evaluate the stored information.

The PathChecker expects a structure like this:

```json5
{
  "/some/path": {
    "minutely": 123,
    "hourly": 456,
    "daily": 789,
    "weekly": 101112,
    "monthly": 131415,
  },
  // following line is a no-op, the route doesn't have to be specified if it should not be limited
  "/some/other/path": {},
  // you can only specify the limits you want to use
  "/yet/another/path": {
    "minutely": 123,
    "hourly": 456,
    "daily": 789,
  },
  // regexes are also supported with `r#` prefix
  "r#^/some/regex/[^/]+/?$": {
    "minutely": 123,
    "daily": 789,
  }
}
```

This package also includes a `PathAndMethodChecker`, which also checks based on the HTTP method, and expects this structure:
**Please note that you're supposed to keep the keys lowercase for these built-in checkers.**

```json5
{
  "get:/some/path": {
    "minutely": 123,
    "hourly": 456,
    "daily": 789,
    "weekly": 101112,
    "monthly": 131415,
  },
  // following line is a no-op, the route doesn't have to be specified if it should not be limited
  "post:/some/other/path": {},
  // you can only specify the limits you want to use
  "delete:/yet/another/path": {
    "minutely": 123,
    "hourly": 456,
    "daily": 789,
  },
  // regexes are also supported with `r#` prefix
  "r#get:^/some/regex/[^/]+/?$": {
    "minutely": 123,
    "daily": 789,
  }
}
```

You can also implement custom checker by implementing `FUPCheckerInterface`.

```go
type FUPCheckerInterface interface {
    Check(fup *FUPScope, c *gin.Context, key string) FUPScopeLimits
}
```

The intervals (minutely, hourly, daily, weekly, monthly) are checked calendarly (so the limits are reset at the beginning of the interval, not after the first request in the interval).

If any of the limits is reached, the middleware will return `429 Too Many Requests` response with the `Retry-After` header set to the time when the interval resets. The payload also contains the surpassed limit information.

If no limit is reached, each response to a request that has limits configured will contain the `X-Client-FUP-Limits`/`X-User-FUP-Limits` header (or both) with the limit values as JSON. E.g.:

```json
{"hourly":{"limit":200,"used":3},"minutely":{"limit":10,"used":1},"weekly":{"limit":100,"used":46}}
```

Usage
------------

Example of real-world usage follows. You can obviously use your own implementation of DB comms, models, etc.

```go
// main.go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go-example/pkg/db"
	"github.com/wernerdweight/api-auth-go-example/pkg/routes"
	"github.com/wernerdweight/api-auth-go/auth"
	"github.com/wernerdweight/api-auth-go/auth/cache"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"github.com/wernerdweight/api-auth-go/auth/entity"
	"github.com/wernerdweight/api-auth-go/auth/fup"
	"github.com/wernerdweight/api-auth-go/auth/provider"
	"github.com/wernerdweight/get-env-go/getenv"
	"gorm.io/gorm"
	"log"
	"os"
)

type ApiClient struct {
	entity.GormApiClient
	DeletedAt gorm.DeletedAt `json:"-"`
}

type ApiUser struct {
	entity.GormApiUser
	DeletedAt gorm.DeletedAt `json:"-"`
}

type ApiUserToken struct {
	entity.GormApiUserToken
	DeletedAt gorm.DeletedAt `json:"-"`
}

func init() {
	err := getenv.Init()
	if nil != err {
		if err.(*getenv.Error).Code != getenv.NoEnvFileError {
			log.Fatal(err.(*getenv.Error).Error())
		}
		log.Print("no .env file found. Expecting ENV variables to be already exported in the system")
	}

	db.Init()
}

func main() {
	r := gin.Default()

	useClientScopeAccessModel := true
	useUserScopeAccessModel := true
	withRegistration := true
	useApiKeyMode := true
	redisDsn, _ := getenv.GetEnv("REDIS_URL") // e.g. redis://localhost:6379/0
	getDBConnection := func() *gorm.DB { return db.GetConnection() }
	newApiClient := func() contract.ApiClientInterface { return &ApiClient{} }
	newApiUser := func() contract.ApiUserInterface { return &ApiUser{} }
	newApiUserToken := func() contract.ApiUserTokenInterface { return &ApiUserToken{} }

	r.Use(auth.Middleware(r, contract.Config{
		Client: contract.ClientConfig{
			Provider: provider.NewGormApiClientProvider(newApiClient, getDBConnection),
			UseScopeAccessModel: &useClientScopeAccessModel,
			FUPChecker:          fup.PathAndMethodFUPChecker{},
		},
		User: &contract.UserConfig{
			Provider: provider.NewGormApiUserProvider(newApiUser, newApiUserToken, getDBConnection),
			TokenFactory: func() contract.ApiUserTokenInterface {
				return &ApiUserToken{}
			},
			WithRegistration:    &withRegistration,
			UseScopeAccessModel: &useUserScopeAccessModel,
			FUPChecker: fup.PathAndMethodFUPChecker{},
		},
		TargetHandlers: &[]string{
			"^/v1/.*",
		},
		Cache: &contract.CacheConfig{
			Driver: cache.NewRedisCacheDriver(redisDsn, newApiClient, newApiUser),
		},
		Mode: &contract.ModesConfig{
			ApiKey: &useApiKeyMode,
		},
	}))

	routes.RegisterRoutes(r)    // register your routes
	
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
```

If you're using GORM, this is what the db package could look like:

```go
// pkg/db/main.go
package db

import (
	_ "github.com/lib/pq"
	"github.com/wernerdweight/api-auth-go-test/model"
	"github.com/wernerdweight/get-env-go/getenv"
	"github.com/wernerdweight/throw-catch-go/throw"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var connection *gorm.DB = nil

func Init() {
	dsn, _ := getenv.GetEnv("DATABASE_URL")
	conn, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	connection = conn
}

func GetConnection() *gorm.DB {
	if nil == connection {
		Init()
	}
	return connection
}

func Migrate() error {
	if nil == connection {
		Init()
	}
	return connection.AutoMigrate(
		&model.ApiClient{},
		&model.ApiUser{},
		&model.ApiUserToken{},
	)
}

func Close() error {
	if nil != connection {
		db, _ := connection.DB()
		return db.Close()
	}
	return nil
}

```

### Retrieving authenticated client/user in targeted handlers/routes

```go
// pkg/routes/main.go (or anywhere else)
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wernerdweight/api-auth-go/auth/constants"
	"log"
	"net/http"
)

func ExampleHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		client, _ := c.Get(constants.ApiClient)
		user, _ := c.Get(constants.ApiUser)
		log.Printf("client: %s, user: %s", client, user)

		// TODO: do the job and return response
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	}
}

func RegisterRoutes(r *gin.Engine) {    // you can do this in your main or wherever you want
	r.POST("/v1/some/path", ExampleHandler())
	...
}

```

### Events

This package dispatches events that you can subscribe to. You can use this to implement your own functionality (e.g. sending confirmation/reset emails, etc.).

The [events-go](https://github.com/wernerdweight/events-go) package is used to dispatch the events. Check the [documentation](https://github.com/wernerdweight/events-go) to see how to subscribe to them.

The following event are dispatched:

```go
// used to validate information provided during registration
// you can subscribe to this event to validate the information and return an error if the information is not valid
type ValidateLoginInformationEvent struct {
	Login    string
	Password string
}

// issued when a new ApiUser is created during registration (before the user is saved)
// you can subscribe to this event to do something with the ApiUser (e.g. set the scope or your custom fields)
type CreateNewApiUserEvent struct {
	ApiUser ApiUserInterface
}

// issued when a new ApiUser is created during registration (after the user is saved)
// you can subscribe to this event to do something with the ApiUser (e.g. send confirmation email)
// NOTE: this event is dispatched asynchronously (returning an error will not affect the registration process)
type RegistrationRequestCompletedEvent struct {
	ApiUser ApiUserInterface
}

// issued when a new ApiUser is confirmed after registration (before the user is saved)
// you can subscribe to this event to do something with the ApiUser (e.g. set the scope or your custom fields)
type ActivateApiUserEvent struct {
	ApiUser ApiUserInterface
}

// issued when a new ApiUser is confirmed after registration (after the user is saved)
// you can subscribe to this event to do something with the ApiUser (e.g. send a tutorial email or activate a plan)
// NOTE: this event is dispatched asynchronously (returning an error will not affect the confirmation process)
type RegistrationConfirmationCompletedEvent struct {
	ApiUser ApiUserInterface
}

// issued when an ApiUser requests a password reset (before the user is saved)
// you can subscribe to this event to do something with the ApiUser (e.g. check user's IP address and/or device)
type RequestResetApiUserPasswordEvent struct {
	ApiUser ApiUserInterface
}

// issued when an ApiUser requests a password reset (after the user is saved)
// you can subscribe to this event to do something with the ApiUser (e.g. send a password reset email)
// NOTE: this event is dispatched asynchronously (returning an error will not affect the reset process)
type ResettingRequestCompletedEvent struct {
    ApiUser ApiUserInterface
}

// issued when an ApiUser resets a password (before the user is saved)
// you can subscribe to this event to do something with the ApiUser (e.g. check user's IP address and/or device)
type ResetApiUserPasswordEvent struct {
    ApiUser ApiUserInterface
}

// issued when an ApiUser resets a password (after the user is saved)
// you can subscribe to this event to do something with the ApiUser
// NOTE: this event is dispatched asynchronously (returning an error will not affect the reset process
type ResettingCompletedEvent struct {
	ApiUser ApiUserInterface
}

```

### Errors

The following errors can occur (you can check for specific code since different errors have different severity):

```go
var AuthErrorCodes = map[AuthErrorCode]string{
    Unknown:                   "unknown error",
    Unauthorized:              "unauthorized",
    ClientNotFound:            "client not found",
    UserNotFound:              "user not found",
    NoCredentialsProvided:     "no credentials provided",
    UserTokenRequired:         "user token required but not provided",
    UserTokenNotFound:         "user token not found",
    UserTokenExpired:          "user token expired",
    ClientForbidden:           "client access forbidden",
    UserForbidden:             "user access forbidden",
    UnknownScopeAccessibility: "unknown scope accessibility",
    UserProviderNotConfigured: "user provider not configured",
    DatabaseError:             "database error",
    InvalidCredentials:        "invalid credentials",
    InvalidRequest:            "invalid request",
    UserAlreadyExists:         "user already exists",
    EncryptionError:           "encryption error",
    UserNotActive:             "user not active",
    ConfirmationTokenExpired:  "confirmation token expired",
    ResettingAlreadyRequested: "resetting already requested",
    ResetTokenExpired:         "reset token expired",
    CacheError:                "cache error",
    MarshallingError:          "marshalling error",
    FUPCacheDisabled:          "cache driver needs to be configured for the FUP checker to work",
    RequestLimitDepleted:      "request limit depleted",
}
```

The payload always has the same structure:

```json5
{
  "code": 8,  // unique error code
  "error": "client access forbidden", // error message for give code
  "payload": null // optional payload, this can be anything
}
```

License
-------
This package is under the MIT license. See the complete license in the root directory of the package.
