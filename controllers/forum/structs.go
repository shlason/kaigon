package forum

import (
	"time"

	"github.com/shlason/kaigon/controllers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const maximumForumNameLength int = 32

type forumInfo struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
	Icon string             `json:"icon"`
}

type readForumsResponse []forumInfo

type createForumRequestPayload struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Banner      string `json:"banner"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

func (p createForumRequestPayload) check() (errResp controllers.JSONResponse, isNotValid bool) {
	if len(p.Name) > maximumForumNameLength || p.Name == "" {
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

type readForumByIDResponse struct {
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

type patchForumRequestPayload struct {
	Name        *string `json:"name" bson:"name,omitempty"`
	Icon        *string `json:"icon" bson:"icon,omitempty"`
	Banner      *string `json:"banner" bson:"banner,omitempty"`
	Rule        *string `json:"rule" bson:"rule,omitempty"`
	Description *string `json:"description" bson:"description,omitempty"`
}

func (p patchForumRequestPayload) check() (errResp controllers.JSONResponse, isNotValid bool) {
	if p.Name != nil {
		if len(*p.Name) > maximumForumNameLength || *p.Name == "" {
			return controllers.JSONResponse{
				Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
				Message: controllers.ErrMessageRequestPayloadFieldNotValid,
				Data:    nil,
			}, true
		}
	}
	if p.Icon != nil {
		if *p.Icon == "" {
			return controllers.JSONResponse{
				Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
				Message: controllers.ErrMessageRequestPayloadFieldNotValid,
				Data:    nil,
			}, true
		}
	}
	if p.Banner != nil {
		if *p.Banner == "" {
			return controllers.JSONResponse{
				Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
				Message: controllers.ErrMessageRequestPayloadFieldNotValid,
				Data:    nil,
			}, true
		}
	}

	return controllers.JSONResponse{}, false
}
