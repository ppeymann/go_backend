package example

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type (
	// AccountService represents method signatures for api account endpoint.
	// so any object that stratifying this interface can be used as account service for api endpoint.
	AccountService interface {
		SignUp(input *SignUpInput, ctx *gin.Context) *BaseResult
	}

	// AccountRepository represents method signatures for account domain repository.
	// so any object that stratifying this interface can be used as account domain repository

	AccountRepository interface {
		Create(input *SignUpInput) (*AccountEntity, error)
		Update(input *AccountEntity) error

		BaseRepository
	}

	// SignUpInput is DTO for parsing register request params.
	//
	// swagger:model SignUpInput
	SignUpInput struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	// AccountEntity Contains account info and is entity of user account that stored on database.
	//
	// swagger:model AccountEntity
	AccountEntity struct {
		gorm.Model `swaggerignore:"true"`

		Mobile    string               `json:"mobile" gorm:"mobile"`
		Password  string               `json:"password" gorm:"password"`
		Role      pq.StringArray       `json:"role" gorm:"type:varchar(64)[]"`
		Email     string               `json:"email" gorm:"email"`
		UserName  string               `json:"user_name" gorm:"user_name"`
		Suspended bool                 `json:"suspended" gorm:"suspended"`
		Tokens    []RefreshTokenEntity `json:"-" gorm:"foreignKey:AccountID;references:ID"`
	}

	// RefreshTokenEntity is entity to store accounts active session
	RefreshTokenEntity struct {
		gorm.Model
		AccountID uint
		TokenId   string `json:"token_id" gorm:"column:token_id;index"`
		UserAgent string `json:"user_agent" gorm:"column:user_agent"`
		IssuedAt  int64  `json:"issued_at" bson:"issued_at" gorm:"column:issued_at"`
		ExpiredAt int64  `json:"expired_at" bson:"expired_at" gorm:"column:expired_at"`
	}

	// TokenBundleOutput Contains Token, Refresh Token and Token Expire time for Login/Verify response DTO.
	//
	// swagger:model TokenBundleOutput
	TokenBundleOutput struct {
		// Token is JWT/PASETO token staring for storing in client side as access token
		Token string `json:"token"`

		// Refresh token string used for refreshing authentication and give fresh token
		Refresh string `json:"refresh"`

		// CentrifugeToken is JWT token for connecting to realtime messaging server
		CentrifugeToken string `json:"centrifuge_token"`

		// Expire time of Token and CentrifugeToken
		Expire time.Time `json:"expire"`

		// Roles contains user access available roles list
		Roles []string `json:"roles"`
	}
)
