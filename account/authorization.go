package account

import example "expamle"

type authorizationService struct {
	next example.AccountService
}

func NewAuthorizationService(service example.AccountService) example.AccountService {
	return &authorizationService{
		next: service,
	}
}
