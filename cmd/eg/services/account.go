package services

import (
	example "expamle"
	"expamle/account"
	"expamle/postgres"
	kitLog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"log"
)

func InitAccountService(db *gorm.DB, logger kitLog.Logger, config *example.Configuration) example.AccountService {
	// initialize and configure mongodb account repository
	accountRepo := postgres.NewAccountRepository(db, config.Database)

	err := accountRepo.Migrate()
	if err != nil {
		log.Panicln(err)
	}

	service := account.NewService(accountRepo)

	path := getSchemaPath("account")

	service, err = account.NewValidationService(path, service)
	if err != nil {
		log.Panicln(err)
	}

	// @Inject logging service to chain
	service = account.NewLoggingService(kitLog.With(logger, "component", "account"), service)

	// @Inject instrumenting service to chain
	service = account.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "account",
			Name:      "request_count",
			Help:      "num of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "account",
			Name:      "request_latency_microseconds",
			Help:      "total duration of requests (ms).",
		}, fieldKeys),
		service,
	)

	// @Inject authorization service to chain
	service = account.NewAuthorizationService(service)

	return service
}
