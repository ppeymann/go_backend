package postgres

import (
	example "expamle"
	"gorm.io/gorm"
)

type accountRepository struct {
	pg       *gorm.DB
	database string
	table    string
}

// NewAccountRepository creates new accountRepository instance and fill its properties with specific argsf
func NewAccountRepository(pg *gorm.DB, database string) example.AccountRepository {
	return &accountRepository{
		pg:       pg,
		database: database,
		table:    "account",
	}
}

func (r *accountRepository) Migrate() error {
	return r.pg.AutoMigrate(&example.AccountEntity{})
}

func (r *accountRepository) Name() string {
	return r.table
}

func (r *accountRepository) Model() *gorm.DB {
	return r.pg.Model(&example.AccountEntity{})
}
