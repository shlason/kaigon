package controllers

const (
	SuccessMessage string = "success"
	SuccessCode    string = "success-200-s"
)

// All common error codes
const (
	// 4XX
	ErrCodeRequestContentTypeNotJSONFormat string = "err-400-rctnjsonf"
	ErrMessageContentTypeNotJSONFromat     string = "Content-Type is not JSON format"

	ErrCodeRequestQueryParamsNotValid    string = "err-400-rqpnv"
	ErrMessageRequestQueryParamsNotValid string = "request query params is not valid"

	ErrCodeRequestPayloadCaptchaFieldNotValid    string = "err-400-rpcfnv"
	ErrMessageRequestPayloadCaptchaFieldNotValid string = "captcha not valid"

	ErrCodeRequestPermissionForbidden    string = "err-403-rpf"
	ErrMessageRequestPermissionForbidden string = "no permission to access"

	// 5XX
	ErrCodeServerGeneralFunctionGotError  string = "err-500-sgfge"
	ErrCodeServerDatabaseQueryGotError    string = "err-500-sdbqge"
	ErrCodeServerDatabaseUpdateGotError   string = "err-500-sdbuge"
	ErrCodeServerRedisSetNXKeyGotError    string = "err-500-srsnxkge"
	ErrCodeServerRedisSetKeyGotError      string = "err-500-srskge"
	ErrCodeServerRedisGetKeyGotError      string = "err-500-srgkge"
	ErrCodeServerSendEmailGotError        string = "err-500-ssege"
	ErrCodeServerGenerateJWTTokenGotError string = "err-500-sgjwttge"
)
