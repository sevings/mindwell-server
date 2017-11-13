// Code generated by go-swagger; DO NOT EDIT.

package restapi

import (
	"crypto/tls"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	accountImpl "github.com/sevings/yummy-server/src/account"
	entriesImpl "github.com/sevings/yummy-server/src/entries"
	favoritesImpl "github.com/sevings/yummy-server/src/favorites"
	usersImpl "github.com/sevings/yummy-server/src/users"
	votesImpl "github.com/sevings/yummy-server/src/votes"
	watchingsImpl "github.com/sevings/yummy-server/src/watchings"
	commentsImpl "github.com/sevings/yummy-server/src/comments"

	"github.com/didip/tollbooth"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/sevings/yummy-server/gen/restapi/operations"
	"github.com/sevings/yummy-server/gen/restapi/operations/account"
	"github.com/sevings/yummy-server/gen/restapi/operations/comments"
	"github.com/sevings/yummy-server/gen/restapi/operations/entries"
	"github.com/sevings/yummy-server/gen/restapi/operations/favorites"
	"github.com/sevings/yummy-server/gen/restapi/operations/me"
	"github.com/sevings/yummy-server/gen/restapi/operations/relations"
	"github.com/sevings/yummy-server/gen/restapi/operations/users"
	"github.com/sevings/yummy-server/gen/restapi/operations/watchings"

	"github.com/sevings/yummy-server/gen/models"

	goconf "github.com/zpatrick/go-config"

	"database/sql"

	_ "github.com/lib/pq"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target ../gen --name  --spec ../swagger-ui/swagger.yaml

func configureFlags(api *operations.YummyAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func loadConfig() *goconf.Config {
	toml := goconf.NewTOMLFile("config.toml")
	loader := goconf.NewOnceLoader(toml)
	config := goconf.NewConfig([]goconf.Provider{loader})
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}
	return config
}

func openDatabase(config *goconf.Config) *sql.DB {
	driver, err := config.StringOr("database.driver", "postgres")
	if err != nil {
		log.Print(err)
	}

	user, err := config.String("database.user")
	if err != nil {
		log.Print(err)
	}

	pass, err := config.String("database.password")
	if err != nil {
		log.Print(err)
	}

	name, err := config.String("database.name")
	if err != nil {
		log.Print(err)
	}

	db, err := sql.Open(driver, "user="+user+" password="+pass+" dbname="+name)
	if err != nil {
		log.Fatal(err)
	}

	schema, err := config.String("database.schema")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("SET search_path = " + schema + ", public")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func configureAPI(api *operations.YummyAPI) http.Handler {
	rand.Seed(time.Now().UTC().UnixNano())

	config := loadConfig()
	db := openDatabase(config)

	accountImpl.ConfigureAPI(db, api)
	usersImpl.ConfigureAPI(db, api)
	entriesImpl.ConfigureAPI(db, api)
	votesImpl.ConfigureAPI(db, api)
	favoritesImpl.ConfigureAPI(db, api)
	watchingsImpl.ConfigureAPI(db, api)
	commentsImpl.ConfigureAPI(db, api)

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.UrlformConsumer = runtime.DiscardConsumer
	api.MultipartformConsumer = runtime.DiscardConsumer
	api.JSONProducer = runtime.JSONProducer()
	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleUi := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Index(r.URL.Path, "/help/api/") == 0 {
			http.StripPrefix("/help/api/", http.FileServer(http.Dir("swagger-ui"))).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})

	lmt := tollbooth.NewLimiter(10, time.Second, nil)
	lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	lmt.SetMessage("")
	lmt.SetMessageContentType("application/json")
	lmt.SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
		err := models.Error{Message: "Too many requests"}
		data, _ := err.MarshalBinary()
		w.Truncate()
		w.Write(data)
	})
	return tollbooth.LimitFuncHandler(lmt, handleUi)
}
