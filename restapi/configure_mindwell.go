// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	accountImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/account"
	admImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/adm"
	badgesImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/badges"
	chatsImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/chats"
	commentsImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/comments"
	complainsImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/complains"
	designImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/design"
	entriesImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/entries"
	favoritesImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/favorites"
	imagesImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/images"
	notificationsImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/notifications"
	oauth2Impl "github.com/sevings/mindwell-server/internal/app/mindwell-server/oauth2"
	relationsImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/relations"
	tagsImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/tags"
	usersImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/users"
	votesImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/votes"
	watchingsImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/watchings"
	wishesImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/wishes"

	"github.com/sevings/mindwell-server/restapi/operations"
	"github.com/sevings/mindwell-server/utils"
)

//go:generate swagger generate server --target .. --name  --spec ../web/swagger.yaml --principal models.UserID

func configureFlags(api *operations.MindwellAPI) {
	_ = api
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MindwellAPI) http.Handler {
	rand.Seed(time.Now().UTC().UnixNano())

	logger, err := zap.NewProduction(zap.WithCaller(false))
	if err != nil {
		log.Println(err)
	}

	systemLogger := logger.With(zap.String("type", "system"))
	_, err = zap.RedirectStdLogAt(systemLogger, zap.ErrorLevel)
	if err != nil {
		systemLogger.Error(err.Error())
	}

	srv := utils.NewMindwellServer(api, "configs/server")
	srv.Eac = utils.NewEmailChecker(srv)

	pm := &utils.Postman{
		BaseUrl:   srv.ConfigString("server.base_url"),
		Support:   srv.ConfigString("server.support"),
		Moderator: srv.ConfigString("server.moderator"),
		Logger:    srv.LogEmail(),
	}

	smtpHost := srv.ConfigString("email.host")
	smtpPort := srv.ConfigInt("email.port")
	smtpUsername := srv.ConfigString("email.username")
	smtpPassword := srv.ConfigString("email.password")
	smtpHelo := srv.ConfigString("email.helo")
	if smtpHost != "" && smtpPort > 0 {
		err = pm.Start(smtpHost, smtpPort, smtpUsername, smtpPassword, smtpHelo)
		if err != nil {
			systemLogger.Error(err.Error())
		}
	}

	srv.Ntf.Mail = pm

	accountImpl.ConfigureAPI(srv)
	admImpl.ConfigureAPI(srv)
	usersImpl.ConfigureAPI(srv)
	entriesImpl.ConfigureAPI(srv)
	votesImpl.ConfigureAPI(srv)
	favoritesImpl.ConfigureAPI(srv)
	watchingsImpl.ConfigureAPI(srv)
	commentsImpl.ConfigureAPI(srv)
	designImpl.ConfigureAPI(srv)
	relationsImpl.ConfigureAPI(srv)
	badgesImpl.ConfigureAPI(srv)
	notificationsImpl.ConfigureAPI(srv)
	complainsImpl.ConfigureAPI(srv)
	chatsImpl.ConfigureAPI(srv)
	tagsImpl.ConfigureAPI(srv)
	oauth2Impl.ConfigureAPI(srv)
	imagesImpl.ConfigureAPI(srv)
	wishesImpl.ConfigureAPI(srv)

	// configure the api here
	api.ServeError = errors.ServeError
	api.Logger = srv.LogSystem().Sugar().Infof
	api.UrlformConsumer = runtime.DiscardConsumer
	api.MultipartformConsumer = runtime.DiscardConsumer
	api.JSONProducer = runtime.JSONProducer()
	api.ServerShutdown = func() {
		srv.Ntf.Stop()
	}

	api.AddMiddlewareFor("GET", "/me", utils.CreateUserLog(srv.DB, srv.LogRequest()))

	regLmt := tollbooth.NewLimiter(1/3600.0, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	regLmt.SetIPLookups([]string{"X-Forwarded-For"})
	regLmt.SetMessage(`{"message":"Слишком много регистраций. Попробуйте позже."}`)
	regLmt.SetMessageContentType("application/json")
	api.AddMiddlewareFor("POST", "/account/register", func(handler http.Handler) http.Handler {
		return tollbooth.LimitHandler(regLmt, handler)
	})

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	_ = tlsConfig
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
	_ = s
	_ = scheme
	_ = addr
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	lmt := tollbooth.NewLimiter(3, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"X-Forwarded-For"})
	lmt.SetMessage(`{"message":"You have reached maximum request limit."}`)
	lmt.SetMessageContentType("application/json")

	return tollbooth.LimitHandler(lmt, handler)
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	logger, err := utils.LogHandler("api", handler)
	if err != nil {
		log.Println(err)
	}

	return logger
}
