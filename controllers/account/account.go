package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers"
)

func SignUp(c *gin.Context) {
	var requestPayload signUpRequestPayload
	errResp, err := controllers.BindJSON(c, &requestPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	errResp, isNotValid := requestPayload.check()
	if isNotValid {
		c.JSON(http.StatusBadRequest, errResp)
	}
}

func SignIn(c *gin.Context) {

}

func CreateVerifySession(c *gin.Context) {

}

func VerifyWithEmail(c *gin.Context) {

}

func CreateResetPasswordSession(c *gin.Context) {

}

func ResetPassword(c *gin.Context) {

}
