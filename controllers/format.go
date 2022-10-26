package controllers

import (
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/utils"
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
