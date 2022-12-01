package forum

import "github.com/shlason/kaigon/controllers"

type createForumRequestPayload struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Banner      string `json:"banner"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

func (p createForumRequestPayload) check() (errResp controllers.JSONResponse, isNotValid bool) {
	const maximumNameLength int = 12

	if len(p.Name) > maximumNameLength || p.Name == "" {
		return controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
			Message: controllers.ErrMessageRequestPayloadFieldNotValid,
			Data:    nil,
		}, true
	}

	if p.Icon == "" || p.Banner == "" {
		return controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
			Message: controllers.ErrMessageRequestPayloadFieldNotValid,
			Data:    nil,
		}, true
	}

	return controllers.JSONResponse{}, false
}
