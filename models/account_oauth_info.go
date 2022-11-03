package models

import (
	"strings"

	"gorm.io/gorm"
)

type AccountOauthInfo struct {
	gorm.Model
	AccountID   uint   `gorm:"not null"`
	AccountUUID string `gorm:"not null"`
	Provider    string `gorm:"not null"`
	Email       string `gorm:"not null"`
}

func (accountOAuthInfo *AccountOauthInfo) Create() *gorm.DB {
	return db.Create(&accountOAuthInfo)
}

func (accountOAuthInfo *AccountOauthInfo) ReadByEmailAndProvider() *gorm.DB {
	return db.First(&accountOAuthInfo, "email = ? AND provider = ?", accountOAuthInfo.Email, accountOAuthInfo.Provider)
}

func (a AccountOauthInfo) ReadAllByAccountUUID(accountUUID string, list *[]AccountOauthInfo) *gorm.DB {
	return db.Where("account_uuid = ?", accountUUID).Find(&list)
}

func (AccountOauthInfo) DeleteByAccountIDs(ids []interface{}) *gorm.DB {
	var fields []string

	for i := 0; i < len(ids); i++ {
		fields = append(fields, "account_id = ?")
	}

	return db.Unscoped().Where(strings.Join(fields, " OR "), ids...).Delete(&AccountOauthInfo{})
}
