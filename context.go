package example

import (
	"errors"
	"expamle/authorization"
	"github.com/gin-gonic/gin"
)

var ErrUserPrincipalsNotFount = errors.New("UserPrincipals not found in context")

// AuthMiddleware defines function pattern for using as authorization control middleware
type AuthMiddleware func(role []string) gin.HandlerFunc

func Principals(ctx *gin.Context) (*authorization.Claims, error) {
	claims, ok := ctx.Get(ContextUserKey)
	if !ok {
		return nil, ErrUserPrincipalsNotFount
	}

	pr, ok := claims.(*authorization.Claims)
	if !ok {
		return nil, ErrUserPrincipalsNotFount
	}

	return pr, nil
}
