package account

import (
	example "expamle"
	"expamle/validation"
)

//type validationService struct {
//	schemas map[string][]byte
//	next    mml_be.AccountService
//}
//
//func NewValidationService(schemaPath string, service mml_be.AccountService) (mml_be.AccountService, error) {
//	schemas := make(map[string][]byte)
//	err := validation.LoadSchemas(schemaPath, schemas)
//	if err != nil {
//		return nil, err
//	}
//
//	return &validationService{
//		schemas: schemas,
//		next:    service,
//	}, nil
//}
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
