package restapi_images

import "C"
import (
	"crypto/tls"
	"github.com/davidbyttow/govips/v2/vips"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	imagesImpl "github.com/sevings/mindwell-server/internal/app/mindwell-images"
	"github.com/sevings/mindwell-server/restapi_images/operations"
	"github.com/sevings/mindwell-server/restapi_images/operations/images"
	"github.com/sevings/mindwell-server/restapi_images/operations/me"
	"github.com/sevings/mindwell-server/restapi_images/operations/themes"
	"github.com/sevings/mindwell-server/utils"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name mindwell-images --spec ../../mindwell-server/web/swagger.yaml --operation PutMeAvatar --operation PutMeCover --principal models.UserID --model Avatar --model Cover --model UserID --model Error

func configureFlags(api *operations.MindwellImagesAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MindwellImagesAPI) http.Handler {
	rand.Seed(time.Now().UTC().UnixNano())

	config := utils.LoadConfig("configs/images")
	mi := imagesImpl.NewMindwellImages(config)

	// configure the api here
	api.ServeError = errors.ServeError
	api.Logger = mi.LogSystem().Sugar().Infof
	api.UrlformConsumer = runtime.DiscardConsumer
	api.MultipartformConsumer = runtime.DiscardConsumer
	api.JSONProducer = runtime.JSONProducer()

	api.OAuth2AppAuth = utils.NewOAuth2App(mi.TokenHash(), mi.DB())
	api.OAuth2PasswordAuth = utils.NewOAuth2User(mi.TokenHash(), mi.DB(), utils.PasswordFlow)
	api.OAuth2CodeAuth = utils.NewOAuth2User(mi.TokenHash(), mi.DB(), utils.CodeFlow)

	api.MePutMeAvatarHandler = me.PutMeAvatarHandlerFunc(imagesImpl.NewAvatarUpdater(mi))
	api.MePutMeCoverHandler = me.PutMeCoverHandlerFunc(imagesImpl.NewCoverUpdater(mi))

	api.ThemesPutThemesNameAvatarHandler = themes.PutThemesNameAvatarHandlerFunc(imagesImpl.NewThemeAvatarUpdater(mi))
	api.ThemesPutThemesNameCoverHandler = themes.PutThemesNameCoverHandlerFunc(imagesImpl.NewThemeCoverUpdater(mi))

	api.ImagesPostImagesHandler = images.PostImagesHandlerFunc(imagesImpl.NewImageUploader(mi))
	api.ImagesGetImagesIDHandler = images.GetImagesIDHandlerFunc(imagesImpl.NewImageLoader(mi))
	api.ImagesDeleteImagesIDHandler = images.DeleteImagesIDHandlerFunc(imagesImpl.NewImageDeleter(mi))

	api.ServerShutdown = func() {
		mi.Shutdown()
		vips.Shutdown()
	}

	vips.LoggingSettings(func(messageDomain string, messageLevel vips.LogLevel, message string) {
		switch messageLevel {
		case vips.LogLevelError:
			mi.TypedLog(messageDomain).Error(message)
		case vips.LogLevelCritical:
			mi.TypedLog(messageDomain).Error(message)
		case vips.LogLevelWarning:
			mi.TypedLog(messageDomain).Warn(message)
		case vips.LogLevelMessage:
			mi.TypedLog(messageDomain).Info(message)
		case vips.LogLevelInfo:
			mi.TypedLog(messageDomain).Info(message)
		case vips.LogLevelDebug:
			mi.TypedLog(messageDomain).Debug(message)
		}
	}, vips.LogLevelWarning)

	vips.Startup(nil)

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
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	getLmt := tollbooth.NewLimiter(3, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	getLmt.SetIPLookups([]string{"X-Forwarded-For"})
	getLmt.SetMessage(`{"message":"Превышено максимальное число запросов."}`)
	getLmt.SetMessageContentType("application/json")
	getLmtHandler := tollbooth.LimitHandler(getLmt, handler)

	postLmt := limiter.New(&limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour}).SetMax(1 / 3600.0).SetBurst(20)
	postLmt.SetIPLookups([]string{"X-Forwarded-For"})
	postLmt.SetMessage(`{"message":"Превышено максимальное число загрузок."}`)
	postLmt.SetMessageContentType("application/json")
	postLmtHandler := tollbooth.LimitHandler(postLmt, handler)

	lmtFunc := func(resp http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost || req.Method == http.MethodPut {
			postLmtHandler.ServeHTTP(resp, req)
		} else {
			getLmtHandler.ServeHTTP(resp, req)
		}
	}

	return http.HandlerFunc(lmtFunc)
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
