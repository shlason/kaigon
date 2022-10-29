package models

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/shlason/kaigon/models/constants"
	"golang.org/x/crypto/bcrypt"
)

// Redis
type Session struct {
	AccountID   uint
	AccountUUID string
	Email       string
	Token       string
}

func (s *Session) Create() error {
	s.Token = uuid.NewString()
	return rdb.SetNX(
		rctx,
		fmt.Sprintf("auth:session:refresh_token:%d/%s/%s", s.AccountID, s.AccountUUID, s.Email),
		s.Token,
		time.Duration(constants.RefreshTokenCookieInfo.MaxAge)*time.Second,
	).Err()
}

func (s *Session) Read() error {
	val, err := rdb.Get(rctx, fmt.Sprintf("auth:session:refresh_token:%d/%s/%s", s.AccountID, s.AccountUUID, s.Email)).Result()
	s.Token = val
	return err
}

func (s *Session) Delete() error {
	return rdb.Del(rctx, fmt.Sprintf("auth:session:refresh_token:%d/%s/%s", s.AccountID, s.AccountUUID, s.Email)).Err()
}

type JWTToken struct {
	AccountID   uint
	AccountUUID string
	Email       string
	jwt.StandardClaims
}

func (tk *JWTToken) Generate() (string, error) {
	ttl := 15 * time.Minute
	tk.ExpiresAt = time.Now().UTC().Add(ttl).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, err := token.SignedString([]byte("my_JWT_secret"))

	return tokenString, err
}

func ParseJWTToken(tk string) (*JWTToken, error) {
	jwtToken, err := jwt.ParseWithClaims(
		tk,
		&JWTToken{},
		func(token *jwt.Token) (i interface{}, e error) {
			return []byte("my_JWT_secret"), nil
		})
	if err != nil || jwtToken == nil {
		return nil, err
	}

	if claim, ok := jwtToken.Claims.(*JWTToken); ok && jwtToken.Valid {
		return claim, nil
	}

	return nil, err
}

// Redis
type AuthCaptcha struct {
	UUID string
	Code string
}

func (ac *AuthCaptcha) Create() error {
	ac.UUID = uuid.NewString()
	return rdb.SetNX(rctx, fmt.Sprintf("auth:captcha:%s", ac.UUID), ac.Code, 5*time.Minute).Err()
}

func (ac *AuthCaptcha) ReadByUUID() error {
	val, err := rdb.Get(rctx, fmt.Sprintf("auth:captcha:%s", ac.UUID)).Result()
	ac.Code = val
	return err
}

func (ac *AuthCaptcha) UpdateByUUID() error {
	return rdb.Set(rctx, fmt.Sprintf("auth:captcha:%s", ac.UUID), ac.Code, 5*time.Minute).Err()
}

func (ac *AuthCaptcha) Delete() error {
	return rdb.Del(rctx, fmt.Sprintf("auth:captcha:%s", ac.UUID)).Err()
}

// Redis
type AuthAccountEmailVerification struct {
	AccountUUID string
	Token       string
	Code        string
	Result      string
}

func (aaev *AuthAccountEmailVerification) Create() error {
	token, err := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), 14)
	if err != nil {
		return err
	}
	code, err := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), 14)
	if err != nil {
		return err
	}
	err = rdb.SetNX(rctx, fmt.Sprintf("auth:account:email:verification:%s", aaev.AccountUUID), fmt.Sprintf("%s/%s", string(token), string(code)), 15*time.Minute).Err()
	if err != nil {
		return err
	}
	aaev.Token = string(token)
	aaev.Code = string(code)
	return err
}

func (aaev *AuthAccountEmailVerification) Read() error {
	val, err := rdb.Get(rctx, fmt.Sprintf("auth:account:email:verification:%s", aaev.AccountUUID)).Result()
	aaev.Result = val
	return err
}

func (aaev *AuthAccountEmailVerification) Update() error {
	token, err := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), 14)
	if err != nil {
		return err
	}
	code, err := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), 14)
	if err != nil {
		return err
	}
	return rdb.Set(rctx, fmt.Sprintf("auth:account:email:verification:%s", aaev.AccountUUID), fmt.Sprintf("%s/%s", string(token), string(code)), 15*time.Minute).Err()
}

func (aaev *AuthAccountEmailVerification) Delete() error {
	return rdb.Del(rctx, fmt.Sprintf("auth:account:email:verification:%s", aaev.AccountUUID)).Err()
}

func (aaev *AuthAccountEmailVerification) IsMatch() bool {
	return aaev.Result == fmt.Sprintf("%s/%s", aaev.Token, aaev.Code)
}

// Redis
type AuthAccountResetPassword struct {
	AccountUUID string
	Token       string
	Code        string
	Result      string
}

func (aarp *AuthAccountResetPassword) Create() error {
	token, err := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), 14)
	if err != nil {
		return err
	}
	code, err := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), 14)
	if err != nil {
		return err
	}
	err = rdb.SetNX(rctx, fmt.Sprintf("auth:account:reset:password:%s", aarp.AccountUUID), fmt.Sprintf("%s/%s", string(token), string(code)), 10*time.Minute).Err()
	if err != nil {
		return err
	}
	aarp.Token = string(token)
	aarp.Code = string(code)
	return err
}

func (aarp *AuthAccountResetPassword) Read() error {
	val, err := rdb.Get(rctx, fmt.Sprintf("auth:account:reset:password:%s", aarp.AccountUUID)).Result()
	aarp.Result = val
	return err
}

func (aarp *AuthAccountResetPassword) Update() error {
	token, err := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), 14)
	if err != nil {
		return err
	}
	code, err := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), 14)
	if err != nil {
		return err
	}
	err = rdb.Set(rctx, fmt.Sprintf("auth:account:reset:password:%s", aarp.AccountUUID), fmt.Sprintf("%s/%s", string(token), string(code)), 10*time.Minute).Err()
	if err != nil {
		return err
	}
	aarp.Token = string(token)
	aarp.Code = string(code)
	return err
}

func (aarp *AuthAccountResetPassword) Delete() error {
	return rdb.Del(rctx, fmt.Sprintf("auth:account:reset:password:%s", aarp.AccountUUID)).Err()
}

func (aarp *AuthAccountResetPassword) IsMatch() bool {
	return aarp.Result == fmt.Sprintf("%s/%s", aarp.Token, aarp.Code)
}
