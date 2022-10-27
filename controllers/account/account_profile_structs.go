package account

import (
	"net/url"

	"github.com/shlason/kaigon/controllers"
)

type socialMediaResponsePayload struct {
	Provider string `json:"provider"`
	UserName string `json:"userName"`
}

type getProfileResponsePayload struct {
	Avatar       string                       `json:"avatar"`
	Banner       string                       `json:"banner"`
	Signature    string                       `json:"signature"`
	SocialMedias []socialMediaResponsePayload `json:"socialMedias"`
}

type patchProfileRequestPayload struct {
	Avatar       *string                      `json:"avatar"`
	Banner       *string                      `json:"banner"`
	Signature    *string                      `json:"signature"`
	SocialMedias []socialMediaResponsePayload `json:"socialMedias"`
}

func (p *patchProfileRequestPayload) check() (errResponse controllers.JSONResponse, isNotValid bool) {
	var acceptSocialMediaProviders = map[string]string{
		"facebook":  "facebook",
		"instagram": "instagram",
		"twitter":   "twitter",
	}

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

	if p.SocialMedias != nil {
		for _, socialMedia := range p.SocialMedias {
			if _, ok := acceptSocialMediaProviders[socialMedia.Provider]; !ok {
				return controllers.JSONResponse{
					Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
					Message: controllers.ErrMessageRequestPayloadFieldNotValid,
					Data:    nil,
				}, true
			}
		}
	}

	return controllers.JSONResponse{}, false
}
