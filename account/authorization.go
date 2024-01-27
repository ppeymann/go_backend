package account

import (
	"github.com/gin-gonic/gin"
	example "github.com/ppeymann/go_backend"
)

type authorizationService struct {
	next example.AccountService
}

func NewAuthorizationService(service example.AccountService) example.AccountService {
	return &authorizationService{
		next: service,
	}
}

func (a *authorizationService) SignUp(input *example.SignUpInput, ctx *gin.Context) *example.BaseResult {
	return a.next.SignUp(input, ctx)
}
