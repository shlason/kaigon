package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
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
