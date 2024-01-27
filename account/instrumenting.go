package account

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/metrics"
	example "github.com/ppeymann/go_backend"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           example.AccountService
}

func NewInstrumentingService(count metrics.Counter, latency metrics.Histogram, service example.AccountService) example.AccountService {
	return &instrumentingService{
		requestCount:   count,
		requestLatency: latency,
		next:           service,
	}
}

func (i *instrumentingService) SignUp(input *example.SignUpInput, ctx *gin.Context) *example.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "SignUp").Add(1)
		i.requestLatency.With("method", "SignUp").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.SignUp(input, ctx)
}
