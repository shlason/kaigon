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

	ErrCodeRequestContentTypeNotFormDataFormat    string = "err-400-rctnfdf"
	ErrMessageRequestContentTypeNotFormDataFormat string = "Content-Type is not form-data format"

	ErrCodeRequestQueryParamsNotValid    string = "err-400-rqpnv"
	ErrMessageRequestQueryParamsNotValid string = "request query params is not valid"

	ErrCodeRequestPayloadCaptchaFieldNotValid    string = "err-400-rpcfnv"
	ErrMessageRequestPayloadCaptchaFieldNotValid string = "captcha not valid"

	ErrCodeRequestPayloadFieldNotValid    string = "err-400-rpfnv"
	ErrMessageRequestPayloadFieldNotValid string = "request payload some field is not valid"

	ErrCodeRequestPermissionUnauthorized    string = "err-401-rpu"
	ErrMessageRequestPermissionUnauthorized string = "unauthorized request"

	ErrCodeRequestPermissionForbidden    string = "err-403-rpf"
	ErrMessageRequestPermissionForbidden string = "no permission to access"

	// 5XX
	ErrCodeServerGeneralFunctionGotError  string = "err-500-sgfge"
	ErrCodeServerDatabaseCreateGotError   string = "err-500-sdbcge"
	ErrCodeServerDatabaseQueryGotError    string = "err-500-sdbqge"
	ErrCodeServerDatabaseUpdateGotError   string = "err-500-sdbuge"
	ErrCodeServerDatabaseDeleteGotError   string = "err-500-sdbdge"
	ErrCodeServerRedisSetNXKeyGotError    string = "err-500-srsnxkge"
	ErrCodeServerRedisSetKeyGotError      string = "err-500-srskge"
	ErrCodeServerRedisGetKeyGotError      string = "err-500-srgkge"
	ErrCodeServerRedisDeleteKeyGotError   string = "err-500-srdkge"
	ErrCodeServerSendEmailGotError        string = "err-500-ssege"
	ErrCodeServerGenerateJWTTokenGotError string = "err-500-sgjwttge"
)
