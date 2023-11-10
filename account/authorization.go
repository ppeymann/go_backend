package account

import (
	example "expamle"
	"github.com/gin-gonic/gin"
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
