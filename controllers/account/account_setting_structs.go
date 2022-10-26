package account

import "github.com/shlason/kaigon/controllers"

type getSettingResponsePayload struct {
	controllers.GormModelResponse
	Name   string `json:"name"`
	Locale string `json:"locale"`
}

type patchSettingRequestPayload struct {
	Name   *string `json:"name"`
	Locale *string `json:"locale"`
}

func (p *patchSettingRequestPayload) check() (errResponse controllers.JSONResponse, isNotValid bool) {
	const nameMaximumLength = 12
	var acceptLocales = map[string]string{
		"TW": "Tw",
	}

	if p.Name != nil {
		if *p.Name == "" || len(*p.Name) > nameMaximumLength {
			return controllers.JSONResponse{
				Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
				Message: controllers.ErrMessageRequestPayloadFieldNotValid,
				Data:    nil,
			}, true
		}
	}

	if p.Locale != nil {
		if _, ok := acceptLocales[*p.Locale]; !ok {
			return controllers.JSONResponse{
				Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
				Message: controllers.ErrMessageRequestPayloadFieldNotValid,
				Data:    nil,
			}, true
		}
	}

	return controllers.JSONResponse{}, false
}

type patchSettingResponsePayload struct {
	controllers.GormModelResponse
	Name   string `json:"name"`
	Locale string `json:"locale"`
}
