package account

import (
	example "expamle"
	"expamle/authorization"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"net/http"
	"time"
)

type service struct {
	repo   example.AccountRepository
	config *example.Configuration
}

func NewService(repo example.AccountRepository, config *example.Configuration) example.AccountService {
	return &service{
		repo:   repo,
		config: config,
	}
}

func (s *service) SignUp(input *example.SignUpInput, ctx *gin.Context) *example.BaseResult {
	account, err := s.repo.Create(input)
	if err != nil {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	refresh := example.RefreshTokenEntity{
		TokenId:   ksuid.New().String(),
		UserAgent: ctx.Request.UserAgent(),
		IssuedAt:  time.Now().UTC().Unix(),
		ExpiredAt: time.Now().Add(time.Duration(s.config.JWT.RefreshExpire) * time.Minute).UTC().Unix(),
	}

	account.Tokens = append(account.Tokens, refresh)

	err = s.repo.Update(account)
	if err != nil {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	paseto, err := authorization.NewPasetoMaker(s.config.JWT.Secret)
	if err != nil {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{example.ErrInternalServer.Error()},
		}
	}

	tokenClaims := &authorization.Claims{
		Subject:   account.ID,
		Roles:     account.Role,
		IssuedAt:  time.Unix(refresh.IssuedAt, 0),
		ExpiredAt: time.Now().Add(time.Duration(s.config.JWT.TokenExpire) * time.Minute).UTC(),
	}

	refreshToken := &authorization.Claims{
		Subject:   account.ID,
		ID:        refresh.TokenId,
		Roles:     account.Role,
		IssuedAt:  time.Unix(refresh.IssuedAt, 0),
		ExpiredAt: time.Unix(refresh.ExpiredAt, 0),
	}

	tokeStr, err := paseto.CreateToken(tokenClaims)
	if err != nil {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{example.ErrInternalServer.Error()},
		}
	}

	refreshStr, err := paseto.CreateToken(refreshToken)
	if err != nil {
		return &example.BaseResult{
			Status: http.StatusOK,
			Errors: []string{
				example.ErrInternalServer.Error(),
			},
		}
	}

	return &example.BaseResult{
		Status: http.StatusOK,
		Result: example.TokenBundleOutput{
			Token:           tokeStr,
			Refresh:         refreshStr,
			CentrifugeToken: "",
			Expire:          tokenClaims.ExpiredAt,
			Roles:           account.Role,
		},
	}
}
