package models

import (
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
	result := db.Model(&accountProfileSocialMedia).Where(
		"account_uuid = ? AND provider = ?",
		accountProfileSocialMedia.AccountUUID,
		accountProfileSocialMedia.Provider,
	).Updates(m)
	if result.RowsAffected == 0 {
		return db.Create(&accountProfileSocialMedia)
	}
	return result
}

func (a AccountProfileSocialMedia) ReadAllByAccountUUID(accountUUID string, list *[]AccountProfileSocialMedia) *gorm.DB {
	return db.Where("account_uuid = ?", accountUUID).Find(&list)
}
