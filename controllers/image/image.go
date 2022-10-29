package image

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shlason/kaigon/configs"
	"github.com/shlason/kaigon/controllers"
)

// TODO: Doc
func UploadToS3(c *gin.Context) {
	var requestPayload *uploadToS3RequestPayload

	uploader := newS3Uploader()

	if !strings.HasPrefix(c.Request.Header.Get("Content-Type"), "multipart/form-data") {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestContentTypeNotFormDataFormat,
			Message: controllers.ErrMessageRequestContentTypeNotFormDataFormat,
			Data:    nil,
		})
		return
	}

	err := c.ShouldBind(&requestPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, controllers.JSONResponse{
			Code:    controllers.ErrCodeRequestPayloadFieldNotValid,
			Message: controllers.ErrMessageRequestPayloadFieldNotValid,
			Data:    nil,
		})
		return
	}

	errResp, isNotValid := requestPayload.check()
	if isNotValid {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	file, err := requestPayload.File.Open()

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    controllers.ErrCodeServerGeneralFunctionGotError,
			Message: err,
			Data:    nil,
		})
		return
	}

	fileName := fmt.Sprintf("%s_%s", uuid.NewString(), requestPayload.File.Filename)

	upInput := &s3manager.UploadInput{
		Bucket:      aws.String(configs.AWS.S3.BucketName),
		Key:         aws.String(fmt.Sprintf("%s/%s", requestPayload.Folder, fileName)),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(requestPayload.File.Header.Get("Content-Type")),
		ACL:         aws.String("public-read"),
	}
	res, err := uploader.UploadWithContext(context.Background(), upInput)

	if err != nil {
		c.JSON(http.StatusInternalServerError, controllers.JSONResponse{
			Code:    errCodeRequestGotErrorWhenUploadImageToS3,
			Message: err,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, controllers.JSONResponse{
		Code:    controllers.SuccessCode,
		Message: controllers.SuccessMessage,
		Data: uploadToS3ResponsePayload{
			S3URL: res.Location,
			URL: fmt.Sprintf(
				"%s://%s/%s/%s",
				configs.Server.Protocol,
				configs.AWS.S3.BucketName,
				requestPayload.Folder,
				fileName,
			),
		},
	})
}

func newS3Uploader() *s3manager.Uploader {
	s3Config := &aws.Config{
		Region:      aws.String(configs.AWS.S3.Region),
		Credentials: credentials.NewStaticCredentials(configs.AWS.S3.AccessKeyID, configs.AWS.S3.AccessSecretKey, ""),
	}
	s3Session, err := session.NewSession(s3Config)

	if err != nil {
		log.Fatal(err)
	}

	uploader := s3manager.NewUploader(s3Session)

	return uploader
}
