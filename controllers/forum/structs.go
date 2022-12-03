package forum

import (
	"time"

	"github.com/shlason/kaigon/controllers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type forumInfo struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
	Icon string             `json:"icon"`
}

type forumReadAllResponse []forumInfo

type forumCreateRequestPayload struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Banner      string `json:"banner"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

func (p forumCreateRequestPayload) check() (errResp controllers.JSONResponse, isNotValid bool) {
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

type forumReadByIDResponse struct {
	ID            primitive.ObjectID `json:"id"`
	CreatedAt     time.Time          `json:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt"`
	DeletedAt     *time.Time         `json:"deletedAt"`
	Name          string             `json:"name"`
	Icon          string             `json:"icon"`
	Banner        string             `json:"banner"`
	Rule          string             `json:"rule"`
	Description   string             `json:"description"`
	PopularTopics []string           `json:"popularTopics"`
}
