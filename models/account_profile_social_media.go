package models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type AccountProfileSocialMedia struct {
	ID          uint   `gorm:"primarykey"`
	AccountID   uint   `gorm:"not null"`
	AccountUUID string `gorm:"not null"`
	Provider    string `gorm:"not null"`
	UserName    string `gorm:"not null"`
}

func (accountProfileSocialMedia *AccountProfileSocialMedia) UpdateOrCreateByAccountUUIDAndProvider(m map[string]interface{}) *gorm.DB {
	result := db.First(&AccountProfileSocialMedia{},
		"account_uuid = ? AND provider = ?",
		accountProfileSocialMedia.AccountUUID,
		accountProfileSocialMedia.Provider,
	)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return db.Create(&accountProfileSocialMedia)
	}
	if result.Error != nil {
		return result
	}
	return db.Model(&accountProfileSocialMedia).Where(
		"account_uuid = ? AND provider = ?",
		accountProfileSocialMedia.AccountUUID,
		accountProfileSocialMedia.Provider,
	).Updates(m)
}

func (a AccountProfileSocialMedia) ReadAllByAccountUUID(accountUUID string, list *[]AccountProfileSocialMedia) *gorm.DB {
	return db.Where("account_uuid = ?", accountUUID).Find(&list)
}

func (AccountProfileSocialMedia) DeleteByAccountIDs(ids []interface{}) *gorm.DB {
	var fields []string

	for i := 0; i < len(ids); i++ {
		fields = append(fields, "account_id = ?")
	}

	return db.Unscoped().Where(strings.Join(fields, " OR "), ids...).Delete(&AccountProfileSocialMedia{})
}
