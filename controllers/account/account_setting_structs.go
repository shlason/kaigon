package account

import "github.com/shlason/kaigon/controllers"

type getSettingResponsePayload struct {
	controllers.GormModelResponse
	Name   string `json:"name"`
	Locale string `json:"locale"`
}
