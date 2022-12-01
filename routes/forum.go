package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisteForumRoutes(publicR *gin.RouterGroup, privateR *gin.RouterGroup) {
	// Get all forum
	publicR.GET("/forums")
	// Create new forum
	privateR.POST("/forums")
	// Get forum info
	publicR.GET("/forums/:forumID")
	// Update forum info
	privateR.PATCH("/forums/:forumID")

	// Get all post from forum
	publicR.GET("/forums/:forumID/posts")
	// Create new post
	privateR.POST("/forums/:forumID/posts")
	// Update post info
	privateR.PATCH("/forums/:forumID/posts/:postID")
	// Get all cooments
	publicR.GET("/forums/:forumID/posts/:postID/comments")
	// Create comment
	privateR.POST("/forums/:forumID/posts/:postID/comments")
	// Update comment
	privateR.POST("/forums/:forumID/posts/:postID/comments/:commentID")
}
