package account

import example "expamle"

type service struct {
	repo example.AccountRepository
}

func NewService(repo example.AccountRepository) example.AccountService {
	return &service{
		repo: repo,
	}
}
