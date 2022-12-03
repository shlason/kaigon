package forum

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ReadAll(c *gin.Context) {
	forums, err := models.Forum{}.Find()
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusOK, controllers.JSONResponse{
			Code:    controllers.SuccessCode,
			Message: controllers.SuccessMessage,
			Data:    forumReadAllResponse{},
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

	var response forumReadAllResponse

	for _, forum := range forums {
		response = append(response, forumInfo{
			ID:   forum.ID,
			Name: forum.Name,
			Icon: forum.Icon,
		})
	}

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data:    response,
	})
}

func Create(c *gin.Context) {
	var requestPayload forumCreateRequestPayload

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

func ReadByID(c *gin.Context) {
	forumID := c.Param("forumID")
	convertedForumID, err := primitive.ObjectIDFromHex(forumID)

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
			Code:    errCodeRequestParamForumIDNoDocument,
			Message: errMessageRequestParamForumIDNoDocument,
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

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data: forumReadByIDResponse{
			ID:            forumModel.ID,
			CreatedAt:     forumModel.CreatedAt,
			UpdatedAt:     forumModel.UpdatedAt,
			DeletedAt:     forumModel.DeletedAt,
			Name:          forumModel.Name,
			Icon:          forumModel.Icon,
			Banner:        forumModel.Banner,
			Rule:          forumModel.Rule,
			Description:   forumModel.Description,
			PopularTopics: forumModel.PopularTopics,
		},
	})
}
