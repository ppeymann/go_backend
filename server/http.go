package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"time"
)

// secure sets http security options for gin framework
func (s *Server) cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTION", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authenticate", "Authorization"},
		ExposeHeaders:    []string{"Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// secure sets http security options for gin framework
func (s *Server) secure() gin.HandlerFunc {
	return secure.New(secure.Config{
		AllowedHosts:          s.config.Listener.AllowedHosts,
		SSLRedirect:           true,
		SSLHost:               s.config.Listener.SSLHost,
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	})
}
