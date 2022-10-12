package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shlason/kaigon/controllers"
	"github.com/shlason/kaigon/models"
	"golang.org/x/crypto/bcrypt"
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
		return
	}
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(requestPayload.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}
	accountModel := &models.Account{
		Uuid:     uuid.NewString(),
		Email:    requestPayload.Email,
		Password: string(hashPwd),
	}
	result := accountModel.Create()
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
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
