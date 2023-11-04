package account

import (
	example "expamle"
	"github.com/go-kit/log"
)

type loggingService struct {
	logger log.Logger
	next   example.AccountService
}

func NewLoggingService(logger log.Logger, service example.AccountService) example.AccountService {
	return &loggingService{
		logger: logger,
		next:   service,
	}
}
