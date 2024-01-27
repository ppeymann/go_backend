package services

import (
	kitLog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/ppeymann/go_backend"
	"github.com/ppeymann/go_backend/postgres"
	"github.com/ppeymann/go_backend/storage"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"log"
)

func InitStorageService(db *gorm.DB, logger kitLog.Logger, config *example.Configuration) example.StorageService {
	storageRepo, err := postgres.NewStorageRepository(db, config.Storage, config.Database, config.JWT.Secret)
	if err != nil {
		log.Panicln(err)
	}

	err = storageRepo.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	storageService := storage.NewService(config.Storage, storageRepo)

	path := getSchemaPath("storage")
	storageService, err = storage.NewValidationService(config.Storage, storageService, config.JWT.Secret, path)
	if err != nil {
		log.Fatal(err)
	}

	// @Inject logging service to chain
	storageService = storage.NewLoggingService(kitLog.With(logger, "compnent", "storage"), storageService)

	// @Inject instrumenting service to chain
	storageService = storage.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "storage",
			Name:      "request_count",
			Help:      "num of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "storage",
			Name:      "request_latency_microseconds",
			Help:      "total duration of requests (ms).",
		}, fieldKeys),
		storageService,
	)

	// @Inject authorization service to chain
	return storage.NewAuthorizationService(config.Storage, storageRepo, storageService)

}
