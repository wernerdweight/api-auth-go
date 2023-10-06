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

Configuration
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
    // '*'            	# all handlers
    // '/v1/*'   		# all handlers starting with '/v1/'
    // '/v1/some/path'  # only '/v1/some/path' handler
    
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

### With API key authentication mode:

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

### Using GORM as data provider:

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

### With cache:

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

### With on-behalf access mode:

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

### With user registration:

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

### With FUP limits:

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

** Retrieving authenticated client/user in targeted handlers/routes **

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

**FUP limits**



**Errors**

The following errors can occur (you can check for specific code since different errors have different severity):

```go
TODO
```

License
-------
This package is under the MIT license. See the complete license in the root directory of the bundle.
