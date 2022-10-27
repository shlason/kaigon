package image

import (
	"fmt"
	"mime/multipart"

	"github.com/shlason/kaigon/controllers"
)

type uploadToS3RequestPayload struct {
	File   *multipart.FileHeader `form:"file"`
	Folder string                `form:"folder"`
}

func (p *uploadToS3RequestPayload) check() (errRespnse controllers.JSONResponse, isNotValid bool) {
	var acceptFolders = map[string]string{
		"account": "account",
		"article": "article",
		"chat":    "chat",
	}
	var acceptImageTypes = map[string]string{
		"image/png":  "image/png",
		"image/jpg":  "image/jpg",
		"image/jpeg": "image/jpeg",
		"image/gif":  "image/gif",
	}
	// 5 MB Limit
	const imageMaximumSize int64 = 5242880

	if _, ok := acceptFolders[p.Folder]; !ok {
		return controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
			Message: controllers.ErrMessageRequestPayloadFieldNotValid,
			Data:    nil,
		}, true
	}

	if p.File == nil {
		return controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
			Message: controllers.ErrMessageRequestPayloadFieldNotValid,
			Data:    nil,
		}, true
	}
	fmt.Println(p.File.Header.Get("Content-Type"))
	if _, ok := acceptImageTypes[p.File.Header.Get("Content-Type")]; !ok || p.File.Size > imageMaximumSize {
		return controllers.JSONResponse{
			Code:    errCodeRequestPayloadImageFileNotValid,
			Message: errMessageRequestPayloadImageFileNotValid,
			Data:    nil,
		}, true
	}

	return controllers.JSONResponse{}, false
}

type uploadToS3ResponsePayload struct {
	S3URL string `json:"s3Url"`
	URL   string `json:"url"`
}
