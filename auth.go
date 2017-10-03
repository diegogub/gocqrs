package gocqrs

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/diegogub/lib"
	"time"
)

func AuthToken(t, secret string) (*SessionClaims, error) {
	token, err := jwt.ParseWithClaims(t, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if token == nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SessionClaims); ok && token.Valid {

		return claims, err
	} else {
		return nil, err
	}
}

func BuildToken(u User) string {
	claims := SessionClaims{
		u.Username,
		u.Role,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(runningApp.sduration).Unix(),
			Issuer:    runningApp.Name,
			Id:        lib.NewShortId(""),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte(runningApp.Secret))

	return tokenString
}

/*
import (
	"github.com/diegogub/lib"
	"time"
)

type Sessioner interface {
	Save(s *Session) error
	Valid(id string) (*Session, error)
}

type Session struct {
	ID         string                 `json:"id"`
	Username   string                 `json:"un"`
	Role       string                 `json:"role"`
	ValidUntil time.Time              `json:"ttl"`
	Data       map[string]interface{} `json:"session"`
}

func NewSession(u *User, validity string) *Session {
	var s Session
	s.ID = lib.NewLongId("S")
	s.Username = u.Username
	s.Role = u.Role

	d, err := time.ParseDuration(validity)
	if err != nil {
		d = time.Duration(time.Minute * 10)
	}

	s.ValidUntil = time.Now().UTC().Add(d)

	return &s
}

type ReadAuther interface {
	AuthRead(e *Entity, username, role string) error
}

*/
