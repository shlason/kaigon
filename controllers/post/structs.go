package post

import "github.com/shlason/kaigon/controllers"

var reactionTypes = map[string]string{
	"like":  "like",
	"sad":   "sad",
	"happy": "happy",
}

var reactionStatsDefault = map[string]int{
	reactionTypes["like"]:  0,
	reactionTypes["sad"]:   0,
	reactionTypes["happy"]: 0,
}

type postCreateRequestPayload struct {
	IsAnonymous  bool    `json:"isAnonymous"`
	ForumID      *string `json:"forumId"`
	Type         string
	Title        string
	Content      string
	Thumbnail    string
	MasterVision string   `json:"masterVision"`
	Topics       []string `json:"topics"`
}

// TODO: 標題、內文字數限制
func (p postCreateRequestPayload) check() (errResp controllers.JSONResponse, isNotValid bool) {
	if p.ForumID == nil || p.Type == "" || p.Title == "" || p.Content == "" {
		return controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
			Message: controllers.ErrMessageRequestPayloadFieldNotValid,
			Data:    nil,
		}, true
	}

	return controllers.JSONResponse{}, false
}
