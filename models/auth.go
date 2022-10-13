package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthCaptcha struct {
	UUID string
	Code string
}

func (ac *AuthCaptcha) Create() error {
	ac.UUID = uuid.NewString()
	return rdb.SetNX(rctx, fmt.Sprintf("auth:captcha:%s", ac.UUID), "captcahCode", 5*time.Minute).Err()
}

func (ac *AuthCaptcha) ReadByUUID() error {
	val, err := rdb.Get(rctx, fmt.Sprintf("auth:captcha:%s", ac.UUID)).Result()
	ac.Code = val
	return err
}

func (ac *AuthCaptcha) UpdateByUUID() error {
	return rdb.Set(rctx, fmt.Sprintf("auth:captcha:%s", ac.UUID), "captcahCode", 5*time.Minute).Err()
}

func (ac *AuthCaptcha) Delete() error {
	return rdb.Del(rctx, fmt.Sprintf("auth:captcha:%s", ac.UUID)).Err()
}

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
	err = rdb.SetNX(rctx, fmt.Sprintf("auth:account:reset:password:%s", aarp.AccountUUID), fmt.Sprintf("%s/%s", string(token), string(code)), 5*time.Minute).Err()
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
	err = rdb.Set(rctx, fmt.Sprintf("auth:account:reset:password:%s", aarp.AccountUUID), fmt.Sprintf("%s/%s", string(token), string(code)), 5*time.Minute).Err()
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
