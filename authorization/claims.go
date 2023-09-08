package authorization

import "time"

type (

	// Claims specify JWT payload claims.
	Claims struct {
		Subject   uint      `json:"sub"`
		ID        string    `json:"jti"`
		Roles     []string  `json:"role"`
		Issuer    string    `json:"iss"`
		Audience  string    `json:"aud"`
		IssuedAt  time.Time `json:"issued_at"`
		ExpiredAt time.Time `json:"exp"`
	}

	TokenMaker interface {
		VerifyToken(token string) (*Claims, error)
		CreateToken(claims *Claims) (string, error)
	}
)

func NewTokenClaims() *Claims {
	return &Claims{}
}

func (c *Claims) withSubject(val uint) *Claims {
	c.Subject = val
	return c
}

func (c *Claims) withID(val string) *Claims {
	c.ID = val
	return c
}

func (c *Claims) withRole(val string) *Claims {
	c.Roles = append(c.Roles, val)
	return c
}

func (c *Claims) withIssuer(val string) *Claims {
	c.Issuer = val
	return c
}

func (c *Claims) withAudience(val string) *Claims {
	c.Audience = val
	return c
}

func (c *Claims) withIssuedAt(val time.Time) *Claims {
	c.IssuedAt = val
	return c
}

func (c *Claims) withExpiredAt(val time.Time) *Claims {
	c.ExpiredAt = val
	return c
}
