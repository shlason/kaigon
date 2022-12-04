package controllers

import (
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func BindJSON(c *gin.Context, r interface{}) (JSONResponse, error) {
	err := c.BindJSON(&r)

	return JSONResponse{
		Code:    ErrCodeRequestContentTypeNotJSONFormat,
		Message: ErrMessageContentTypeNotJSONFromat,
		Data:    nil,
	}, err
}

func GetFilteredNilRequestPayloadMap(any interface{}) map[string]interface{} {
	m := structs.Map(any)
	utils.FilterNilMap(&m)
	return m
}

func GetFilteredNilRequestPayloadBsonM(any interface{}) (bson.M, error) {
	resultByte, err := bson.Marshal(any)
	if err != nil {
		return nil, err
	}
	bsonM := bson.M{}
	err = bson.Unmarshal(resultByte, &bsonM)

	return bsonM, err
}
