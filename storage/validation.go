package storage

import (
	"errors"
	example "expamle"
	"expamle/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type validationService struct {
	schemas map[string][]byte
	secret  string
	opts    example.StorageOption
	next    example.StorageService
}

var ErrInvalidAccessToken = errors.New("invalid photo access token")
var ErrBadUploadRequest = errors.New("upload request is illegal")

func NewValidationService(opts example.StorageOption, service example.StorageService, secret, schemaPath string) (example.StorageService, error) {
	schemas := make(map[string][]byte)
	err := validation.LoadSchemas(schemaPath, schemas)
	if err != nil {
		return nil, err
	}

	return &validationService{
		schemas: schemas,
		secret:  secret,
		opts:    opts,
		next:    service,
	}, nil
}

func (v *validationService) Upload(input *example.UploadInput, ctx *gin.Context) *example.BaseResult {
	tag := example.ObjectTag(input.Tag)
	input.Tag = string(tag)

	file, err := ctx.FormFile("file")
	if err != nil {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{ErrBadUploadRequest.Error()},
		}
	}

	ct, ok := file.Header["Content-Type"]
	if !ok {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{ErrBadUploadRequest.Error()},
		}
	}

	input.ContentType = ct[0]
	input.Size = file.Size

	validationError := validation.Validate(input, v.schemas)
	if validationError != nil {
		return validationError
	}

	return v.next.Upload(input, ctx)
}
