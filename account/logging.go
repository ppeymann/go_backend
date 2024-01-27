package account

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	example "github.com/ppeymann/go_backend"
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

func (l *loggingService) SignUp(input *example.SignUpInput, ctx *gin.Context) (result *example.BaseResult) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "SignUp",
			"errors", strings.Join(result.Errors, ","),
			"mobile", input.Mobile,
			"client_ip", ctx.ClientIP(),
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.SignUp(input, ctx)
}
