package auth

import "github.com/gin-gonic/gin"

type getCaptchaInfoResponsePayload struct {
	UUID string `json:"uuid"`
}

type getCaptchaImageRequestParamsPayload struct {
	CaptchaUUID string `json:"captchaUuid"`
}

func (p *getCaptchaImageRequestParamsPayload) BindParams(c *gin.Context) {
	p.CaptchaUUID = c.Param("captchaUUID")
}

type updateCaptchaImageRequestParamsPayload struct {
	CaptchaUUID string `json:"captchaUuid"`
}

func (p *updateCaptchaImageRequestParamsPayload) BindParams(c *gin.Context) {
	p.CaptchaUUID = c.Param("captchaUUID")
}
