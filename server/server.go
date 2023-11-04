package server

import (
	"context"
	example "expamle"
	"expamle/authorization"
	"expamle/env"
	"fmt"
	"github.com/gin-gonic/gin"
	kitLog "github.com/go-kit/kit/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server holds the dependencies for an HTTP server.
type Server struct {
	account example.AccountService

	Router        *gin.Engine
	config        *example.Configuration
	instrumenting serviceInstrumenting

	paseto authorization.TokenMaker
	logger kitLog.Logger
}

// EnvMode specified the running env 'release' represents production mode and ” represents development.
// it depended on gin GIN_MODE env for unifying and simplicity of setting.
var EnvMode = ""

func NewServer(logger kitLog.Logger, config *example.Configuration) *Server {
	svr := &Server{
		logger:        logger,
		config:        config,
		instrumenting: newServiceInstrumenting(),
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// determining environment
	EnvMode = os.Getenv("GIN_MODE")

	// setting swagger info if not in production mode
	if env.GetStringDefault("SWAGGER_ENABLE", "false") == "true" {
		docs.SwaggerInfo.Title = fmt.Sprintf("Go BackEnd [AuthMode: %s]", config.Listener.AuthMode)
		docs.SwaggerInfo.Description = "the Swagger Documentation For Go Backend Example Api Server"
		docs.SwaggerInfo.Version = "v1"
		docs.SwaggerInfo.Host = env.GetStringDefault("HOST_URL", "localhost:8080")
		docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	}

	// binding global metrics middleware
	router.Use(svr.metrics())

	if env.GetStringDefault("CORS_ENABLE", "false") == "true" {
		router.Use(svr.cors())
	}

	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Fatal(err)
	}

	svr.Router = router

	svr.Router.GET("/metric", svr.prometheus())

	svr.paseto, err = authorization.NewPasetoMaker(svr.config.JWT.Secret)
	if err != nil {
		log.Fatal(err)

	}

	return svr
}

// Listen start listening address for incoming request and handle gracefully shutdown
func (s *Server) Listen() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if env.GetStringDefault("SWAGGER_ENABLED", "false") == "true" {
		s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Listener.Host, s.config.Listener.Port),
		Handler: s.Router,
	}

	// start https server

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http listener stopped : ", err)
		}
	}()

	// Listen for the interrupt signal.=
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully mml_be server, press Ctrl+C again to force")

	// The context is used to inform the server it has 30 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shotdown: ", err)
	}
	log.Println("Go Backend service exiting")
}
