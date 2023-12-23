package storage

import (
	example "expamle"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	authorizationService struct {
		opts example.StorageOption
		repo example.StorageRepository
		next example.StorageService
	}
)

func NewAuthorizationService(opts example.StorageOption, repo example.StorageRepository, service example.StorageService) example.StorageService {
	return &authorizationService{
		opts: opts,
		repo: repo,
		next: service,
	}
}

func (a *authorizationService) Upload(input *example.UploadInput, ctx *gin.Context) *example.BaseResult {
	cliams, err := example.Principals(ctx)
	if err != nil {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{ErrInvalidAccessToken.Error()},
		}
	}

	input.Claims = cliams

	return a.next.Upload(input, ctx)
}
