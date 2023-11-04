package example

import (
	"errors"
	"gorm.io/gorm"
)

type (
	// BaseRepository is abstract interface that all repository must implement its method
	BaseRepository interface {
		// Migrate runs AutoMigrate for expected repository model
		Migrate() error

		// Name repository associated table name
		Name() string

		// Model returns *gorm.DB instance for repository
		Model() *gorm.DB
	}

	// BaseResult a basic Golang struct which includes the following field: Success, Errors, Messages, ResultCount, Result
	// It is the unified response model for entire service api calls
	//
	// swagger:model BaseResult
	BaseResult struct {
		Status int `json:"-"`

		// Errors provides list of error that occurred in processing request
		Errors []string `json:"errors" mapstructure:"errors"`

		// ResultCount specified number of records that returned in result_count field expected result been array.
		ResultCount int64 `json:"result_count,omitempty" mapstructure:"result_count"`

		// Result single/array of any type (object/number/string/boolean) that returns as response
		Result interface{} `json:"result" mapstructure:"result"`
	}
)

var (
	ErrUnimplementedRequest = errors.New("request is not implemented")
	ErrUnhandled            = errors.New("an unhandled error occurred during processing the request")
	ErrNotFound             = errors.New("not found")
	ErrInternalServer       = errors.New("internal server error")
	ErrEntityAlreadyExist   = errors.New("entity with specified properties already exist")
)

var AllRoles = []string{UserRole}

const (
	ContextUserKey           string = "CONTEXT_USER"
	UserSessionKey           string = "USER_SESSION"
	AuthorizationFailed      string = "authorization failed"
	ProvidedRequiredJsonBody string = "please provided required JSOn body"
)

const (
	UserRole string = "USER"
)
