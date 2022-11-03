package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/develop"
)

func RegisteDevelopUtilsRoutes(r *gin.RouterGroup) {
	r.GET("/develop/utils/account/delete", develop.DeleteAccount)
}
