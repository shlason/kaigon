package controllers

const (
	SuccessMessage string = "success"
	SuccessCode    string = "success-200-s"
)

const (
	// All common error codes
	// 4XX
	ErrCodeRequestContentTypeNotJSONFormat    string = "err-400-rctnjsonf"
	ErrCodeRequestPayloadCaptchaFieldNotValid string = "err-400-rpcfnv"

	// 5XX
	ErrCodeServerGeneralFunctionGotError string = "err-500-sgfge"
	ErrCodeServerDatabaseQueryGotError   string = "err-500-sdbqge"
	ErrCodeServerRedisSetNXKeyGotError   string = "err-500-srsnxkge"
	ErrCodeServerRedisSetKeyGotError     string = "err-500-srskge"
	ErrCodeServerRedisGetKeyGotError     string = "err-500-srgkge"

	// All common error message
	ErrMessageContentTypeNotJSONFromat           string = "Content-Type is not JSON format"
	ErrMessageRequestPayloadCaptchaFieldNotValid string = "captcha not valid"
)
