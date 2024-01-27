package account

import (
	example "github.com/ppeymann/go_backend"

	"github.com/gin-gonic/gin"
	"github.com/ppeymann/go_backend/validation"
)

type validationService struct {
	schemas map[string][]byte
	next    example.AccountService
}

func NewValidationService(schemaPath string, service example.AccountService) (example.AccountService, error) {
	schemas := make(map[string][]byte)
	err := validation.LoadSchemas(schemaPath, schemas)
	if err != nil {
		return nil, err
	}

	return &validationService{
		schemas: schemas,
		next:    service,
	}, nil
}

func (v *validationService) SignUp(input *example.SignUpInput, ctx *gin.Context) *example.BaseResult {
	err := validation.Validate(input, v.schemas)
	if err != nil {
		return err
	}

	return v.next.SignUp(input, ctx)
}
