package server

import (
	example "expamle"
	"github.com/gin-gonic/gin"
	"net/http"
)

type accountHandler struct {
	service example.AccountService
	config  *example.Configuration
}

// InitAccountHandlers injects mml_be.AccountService to Server account field.
func (s *Server) InitAccountHandlers(svc example.AccountService, config *example.Configuration) {
	s.account = svc
	handler := accountHandler{
		service: svc,
		config:  config,
	}

	group := s.Router.Group("api/v1/account")
	{
		group.POST("/signup", handler.SignUp)
	}
}

func (h *accountHandler) SignUp(ctx *gin.Context) {
	input := &example.SignUpInput{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.JSON(http.StatusBadRequest, example.BaseResult{
			Errors: []string{example.ProvidedRequiredJsonBody},
		})
	}

	result := h.service.SignUp(input, ctx)
	ctx.JSON(result.Status, result)
}
