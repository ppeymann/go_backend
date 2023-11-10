package account

import (
	example "expamle"
	"expamle/validation"
	"github.com/gin-gonic/gin"
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
