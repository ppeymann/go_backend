package storage

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/metrics"
	example "github.com/ppeymann/go_backend"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           example.StorageService
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, service example.StorageService) example.StorageService {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		next:           service,
	}
}

func (i *instrumentingService) Upload(input *example.UploadInput, ctx *gin.Context) *example.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "Upload").Add(1)
		i.requestLatency.With("method", "Upload").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Upload(input, ctx)
}
