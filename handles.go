package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// 临时存储密码，真正的秘密应该存储在数据库中
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check password
	exceptPass, ok := users[creds.Username]
	if !ok || exceptPass != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 创建token
	token := NewToken()
	tokenStr, expireTime, err := token.Generate(creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("got token:%s \n", tokenStr)

	// 最后，我们将客户端cookie token设置为刚刚生成的JWT
	// 我们还设置了与令牌本身相同的cookie到期时间
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expireTime,
	})
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	// 从cookie中获取token
	c, err := r.Cookie("token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// 其他类型
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := c.Value

	// 验证token
	token := NewToken()
	claims, err := token.Validate(tokenStr)
	if err != nil {
		fmt.Printf("validate error:%v\n", err)
		if errors.Is(err, jwt.ErrSignatureInvalid) ||
			errors.Is(err, jwt.ErrTokenSignatureInvalid) ||
			errors.Is(err, jwt.ErrTokenExpired) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// 获取token
	// 从cookie中获取token
	c, err := r.Cookie("token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// 其他类型
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := c.Value

	// 验证token
	token := NewToken()
	tokenStr, expireTime, err := token.Refresh(tokenStr)
	if err != nil {
		fmt.Printf("%v\n", err)
		if errors.Is(err, ErrNotRefreshTime) {
			w.Write([]byte("not yet refresh time"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("got token:%s \n", tokenStr)

	// 最后，我们将客户端cookie token设置为刚刚生成的JWT
	// 我们还设置了与令牌本身相同的cookie到期时间
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expireTime,
	})
}
