package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePost(c *gin.Context) {
	var requestPayload postCreateRequestPayload

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

	convertedForumID, err := primitive.ObjectIDFromHex(*requestPayload.ForumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	forumModel := models.Forum{
		MongoDBModel: models.MongoDBModel{
			ID: convertedForumID,
		},
	}
	err = forumModel.FindOneByID()
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    errCodeRequestPayloadForumIDNotExist,
			Message: errMessageRequestPayloadForumIDNotExist,
			Data:    nil,
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseQueryGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	authPayload := c.MustGet("authPayload").(*models.JWTToken)

	// TODO: topics 相關處理
	postModel := models.Post{
		AccountID:     authPayload.AccountID,
		IsAnonymous:   requestPayload.IsAnonymous,
		ForumID:       convertedForumID,
		Type:          requestPayload.Type,
		Title:         requestPayload.Title,
		Content:       requestPayload.Content,
		Thumbnail:     requestPayload.Thumbnail,
		MasterVision:  requestPayload.MasterVision,
		CommentCount:  0,
		Topics:        requestPayload.Topics,
		ReactionStats: reactionStatsDefault,
	}
	_, err = postModel.InsertOne()
	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerDatabaseCreateGotError,
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
