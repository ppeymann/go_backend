package account

import (
	example "expamle"
	"github.com/go-kit/kit/metrics"
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
