package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/controllers/forum"
	"github.com/shlason/kaigon/controllers/post"
)

// TODO: forum 相關 private 操作要加上版主權限管理
func RegisteForumAndPostRoutes(publicR *gin.RouterGroup, privateR *gin.RouterGroup) {
	// Get all forum
	publicR.GET("/forums", forum.ReadFroums)
	// Create new forum
	// TODO: 記得改回 private
	publicR.POST("/forums", forum.CreateForum)
	// Get forum info
	publicR.GET("/forums/:forumID", forum.ReadForumByID)
	// Update forum info
	privateR.PATCH("/forums/:forumID", forum.PatchForum)

	// Get all post from forum
	publicR.GET("/forums/:forumID/posts")
	// Create new post
	privateR.POST("/forums/:forumID/posts", post.CreatePost)
	// Update post info
	privateR.PATCH("/forums/:forumID/posts/:postID")
	privateR.POST("/forums/:forumID/posts/:postID/favorite")
	privateR.POST("/forums/:forumID/posts/:postID/followed")
	privateR.PATCH("/forums/:forumID/posts/:postID/reactions")
	// Get all cooments
	publicR.GET("/forums/:forumID/posts/:postID/comments")
	// Create comment
	privateR.POST("/forums/:forumID/posts/:postID/comments")
	// Update comment
	privateR.POST("/forums/:forumID/posts/:postID/comments/:commentID")
	privateR.POST("/forums/:forumID/posts/:postID/comments/:commentID/like")
}
