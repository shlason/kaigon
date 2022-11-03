package models

import (
	"strings"

	"gorm.io/gorm"
)

type AccountSetting struct {
	gorm.Model
	AccountID   uint   `gorm:"unique; not null;"`
	AccountUUID string `gorm:"unique; not null;"`
	Name        string
	Locale      string
}

func (accountSetting *AccountSetting) Create() *gorm.DB {
	return db.Create(&accountSetting)
}

func (accountSetting *AccountSetting) ReadByAccountUUID() *gorm.DB {
	return db.First(&accountSetting, "account_uuid = ?", accountSetting.AccountUUID)
}

func (accountSetting *AccountSetting) UpdateByAccountUUID(m map[string]interface{}) *gorm.DB {
	return db.Model(&accountSetting).Where("account_uuid = ?", accountSetting.AccountUUID).Updates(m)
}

func (AccountSetting) DeleteByAccountIDs(ids []interface{}) *gorm.DB {
	var fields []string

	for i := 0; i < len(ids); i++ {
		fields = append(fields, "account_id = ?")
	}

	return db.Unscoped().Where(strings.Join(fields, " OR "), ids...).Delete(&AccountSetting{})
}
