package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ppeymann/go_backend/utils"
	"github.com/thoas/go-funk"
)

// Authenticate is authentication and Authenticate middleware for http request
func (s *Server) pasetoAuth(roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// catch Authenticate header from context
		authHeader := ctx.GetHeader("Authorization")
		if len(authHeader) == 0 {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("authorization header is not provided"))

			return
		}

		// Bearer token format validation
		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("invalid Authorization header format"))

			return
		}

		at := strings.ToLower(fields[0])
		if at != "Bearer" {
			_ = ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unsupported Authenticate format : %s", fields[0]))

			return
		}

		token := fields[1]
		claims, err := s.paseto.VerifyToken(token)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)

			return
		}

		if claims.ExpiredAt.Before(time.Now().UTC()) {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("authorization token is expired"))
		}

		if len(roles) > 0 {
			if len(funk.IntersectString(roles, claims.Roles)) == 0 {
				_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("permission denied"))

				return
			}
		}

		ctx.Set(utils.ContextUserKey, claims)
		ctx.Set(utils.ContextRoleKey, claims.Roles)

	}
}

// Authenticate is authentication and Authenticate middleware for http request
func (s *Server) Authenticate(roles []string) gin.HandlerFunc {
	return s.pasetoAuth(roles)
}
