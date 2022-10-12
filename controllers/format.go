package controllers

import "github.com/gin-gonic/gin"

func BindJSON(c *gin.Context, r interface{}) (JSONResponse, error) {
	err := c.BindJSON(&r)

	return JSONResponse{
		Code:    ErrCodeRequestContentTypeNotJSONFormat,
		Message: ErrMessageContentTypeNotJSONFromat,
		Data:    nil,
	}, err

}
