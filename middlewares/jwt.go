package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
)

type authHeader struct {
	Token string `header:"Authorization"`
}

const (
	errCodeRequestHeaderAuthorizationFieldNotValid    = "err-401-rhafnv"
	errMessageRequestHeaderAuthorizationFieldNotValid = "authorization failed or format error `Bearer {token}`"
)

func JWT(c *gin.Context) {
	h := authHeader{}
	err := c.ShouldBindHeader(&h)
	if err != nil {
		c.JSON(http.StatusUnauthorized, controllers.JSONResponse{
			Code:    errCodeRequestHeaderAuthorizationFieldNotValid,
			Message: errMessageRequestHeaderAuthorizationFieldNotValid,
			Data:    nil,
		})
		c.Abort()
		return
	}

	tokenHeader := strings.Split(h.Token, "Bearer ")

	if len(tokenHeader) < 2 {
		c.JSON(http.StatusUnauthorized, controllers.JSONResponse{
			Code:    errCodeRequestHeaderAuthorizationFieldNotValid,
			Message: errMessageRequestHeaderAuthorizationFieldNotValid,
			Data:    nil,
		})
		c.Abort()
		return
	}

	_, err = models.ParseJWTToken(tokenHeader[1])

	if err != nil {
		c.JSON(http.StatusUnauthorized, controllers.JSONResponse{
			Code:    errCodeRequestHeaderAuthorizationFieldNotValid,
			Message: errMessageRequestHeaderAuthorizationFieldNotValid,
			Data:    nil,
		})
		c.Abort()
		return
	}
}
