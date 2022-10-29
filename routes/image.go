package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/image"
)

func RegisteImageRoutes(privateR *gin.RouterGroup) {
	privateR.POST("/image", image.UploadToS3)
}
