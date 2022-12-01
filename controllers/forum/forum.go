package forum

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func Create(c *gin.Context) {
	var requestPayload createForumRequestPayload

	errResp, err := controllers.BindJSON(c, &requestPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	errResp, isNotValid := requestPayload.check()
	if isNotValid {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	forumModel := models.Forum{
		Name:          requestPayload.Name,
		Icon:          requestPayload.Icon,
		Banner:        requestPayload.Banner,
		Rule:          requestPayload.Rule,
		Description:   requestPayload.Description,
		PopularTopics: []string{},
	}
	err = forumModel.FindOneByName()

	if err == nil {
		c.JSON(http.StatusConflict, controllers.JSONResponse{
			Code:    errCodeRequestPayloadForumNameFieldAlreadyExist,
			Message: errMessageRequestPayloadForumNameFieldAlreadyExist,
			Data:    nil,
		})
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseQueryGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	_, err = forumModel.InsertOne()

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseQueryGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data:    nil,
	})
}
