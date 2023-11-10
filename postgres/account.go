package postgres

import (
	example "expamle"
	"expamle/env"
	"expamle/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"strings"
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

func (r *accountRepository) Create(input *example.SignUpInput) (*example.AccountEntity, error) {
	role := []string{
		input.Role,
	}

	if role[0] != example.UserRole {
		role = append(role, example.UserRole)
	}

	account := &example.AccountEntity{
		Model:     gorm.Model{},
		Mobile:    input.Mobile,
		Role:      role,
		Email:     input.Mobile,
		UserName:  input.Mobile,
		Suspended: false,
	}

	if env.IsProduction() {
		hash, err := utils.HashString(input.Password)
		if err != nil {
			return nil, example.ErrInternalServer
		}

		account.Password = hash
	} else {
		account.Password = input.Password
	}

	err := r.pg.Transaction(func(tx *gorm.DB) error {
		if res := r.Model().Create(account).Error; res.Error != nil {
			str := res.(*pgconn.PgError).Message
			if strings.Contains(str, "duplicate key value") {
				return res
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return account, nil

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
