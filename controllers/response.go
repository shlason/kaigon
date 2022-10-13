package controllers

const SuccessMessage string = "success"
const SuccessCode string = "success-200-s"

// All common error codes
// 4XX
const ErrCodeRequestContentTypeNotJSONFormat string = "err-400-rctnjsonf"
const ErrCodeRequestPayloadCaptchaFieldNotValid string = "err-400-rpcfnv"

// 5XX
const ErrCodeServerGeneralFunctionGotError string = "err-500-sgfge"
const ErrCodeServerDatabaseQueryGotError string = "err-500-sdbqge"

// All common error message
const ErrMessageContentTypeNotJSONFromat string = "Content-Type is not JSON format"
const ErrMessageRequestPayloadCaptchaFieldNotValid string = "captcha not valid"
