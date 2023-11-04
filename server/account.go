package server

import example "expamle"

type accountHandler struct {
	service example.AccountService
	config  *example.Configuration
}

// InitAccountHandlers injects mml_be.AccountService to Server account field.
func (s *Server) InitAccountHandlers(svc example.AccountService, config *example.Configuration) {
	s.account = svc
	//handler := accountHandler{
	//	service: svc,
	//	config:  config,
	//}

	//group := s.Router.Group("api/v1/account")
	//{
	//
	//}
}
