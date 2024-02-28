package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

var mySigningKey = []byte("AllYourBase")
var expireDuration time.Duration = 30 * time.Minute
var ErrNotRefreshTime = errors.New("not refresh time")

type Token interface {
	Generate(username string) (tokenStr string, expiretime time.Time, err error)
	Validate(tokenStr string) (claims *MyCustomClaims, err error)
	Refresh(tokenStr string) (newTokenStr string, expiretime time.Time, err error)
}

type token struct {
}

func NewToken() Token {
	return &token{}
}

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (t *token) Generate(username string) (tokenStr string, expiretime time.Time, err error) {
	// 过期时间
	expiretime = time.Now().Add(expireDuration)

	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		username,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(expiretime),
		},
	}

	fmt.Printf("username: %v\n", claims.Username)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString(mySigningKey)
	if err != nil {
		return
	}
	return
}

func (t *token) Validate(tokenStr string) (claims *MyCustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return
	}

	if !token.Valid {
		err = jwt.ErrTokenSignatureInvalid
		return
	}

	var ok bool
	if claims, ok = token.Claims.(*MyCustomClaims); !ok {
		err = errors.Errorf("unknown claims type, cannot proceed")
		return
	}
	return
}

func (t *token) Refresh(tokenStr string) (newTokenStr string, expiretime time.Time, err error) {
	// 验证旧token
	claims, err := t.Validate(tokenStr)
	if err != nil {
		return
	}

	// 我们确保在足够的时间之前不会发行新令牌。
	// 在这种情况下，仅当旧令牌在30秒到期时才发行新令牌。
	// 否则，返回错误的请求状态。
	if claims.ExpiresAt.Sub(time.Now()) > time.Second*30 {
		err = ErrNotRefreshTime
		return
	}

	// 生成新token
	return t.Generate(claims.Username)
}
