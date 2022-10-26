package account

import (
	"net/url"

	"github.com/shlason/kaigon/controllers"
)

type getProfileResponsePayload struct {
	Avatar    string `json:"avatar"`
	Banner    string `json:"banner"`
	Signature string `json:"signature"`
}

type patchProfileRequestPayload struct {
	Avatar    *string `json:"avatar"`
	Banner    *string `json:"banner"`
	Signature *string `json:"signature"`
}

func (p *patchProfileRequestPayload) check() (errResponse controllers.JSONResponse, isNotValid bool) {
	if p.Avatar != nil && *p.Avatar != "" {
		_, err := url.ParseRequestURI(*p.Avatar)
		if err != nil {
			return controllers.JSONResponse{
				Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
				Message: controllers.ErrMessageRequestPayloadFieldNotValid,
				Data:    nil,
			}, true
		}
	}

	if p.Banner != nil && *p.Banner != "" {
		_, err := url.ParseRequestURI(*p.Banner)
		if err != nil {
			return controllers.JSONResponse{
				Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
				Message: controllers.ErrMessageRequestPayloadFieldNotValid,
				Data:    nil,
			}, true
		}
	}

	if p.Signature != nil && *p.Signature != "" {
		_, err := url.ParseRequestURI(*p.Signature)
		if err != nil {
			return controllers.JSONResponse{
				Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
				Message: controllers.ErrMessageRequestPayloadFieldNotValid,
				Data:    nil,
			}, true
		}
	}

	return controllers.JSONResponse{}, false
}
