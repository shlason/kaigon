package controllers

import (
	"time"

	"gorm.io/gorm"
)

type JSONResponse struct {
	Code    string      `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type GormModelResponse struct {
	ID        uint           `json:"id,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}
