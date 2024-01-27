package storage

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	example "github.com/ppeymann/go_backend"
)

type loggingService struct {
	logger log.Logger
	next   example.StorageService
}

func NewLoggingService(logger log.Logger, service example.StorageService) example.StorageService {
	return &loggingService{
		logger: logger,
		next:   service,
	}
}

func (l *loggingService) Upload(input *example.UploadInput, ctx *gin.Context) (result *example.BaseResult) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "Upload",
			"user", input.Claims.Subject,
			"tag", input.Tag,
			"size", input.Size,
			"errors", strings.Join(result.Errors, ","),
			"client_ip", ctx.GetHeader("X-Real-Ip"),
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.Upload(input, ctx)
}
