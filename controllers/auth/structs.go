package auth

import "github.com/gin-gonic/gin"

type getCaptchaInfoResponsePayload struct {
	UUID string
}

type getCaptchaImageRequestParamsPayload struct {
	CaptchaUUID string
}

func (p *getCaptchaImageRequestParamsPayload) BindParams(c *gin.Context) {
	p.CaptchaUUID = c.Param("captchaUUID")
}
