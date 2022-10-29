package auth

const (
	ErrCodeRequestQueryParamsAccountUUIDFieldNotValid    string = "err-400-rqpaufnv"
	ErrMessageRequestQueryParamsAccountUUIDFieldNotValid string = "query string 'accountUuid' in not valid"

	ErrCodeRequestQueryParamsOAuthURLTypeFieldNotValid    string = "err-400-rqpoutfnv"
	ErrMessageRequestQueryParamsOAuthURLTypeFieldNotValid string = "query string 'type' in not valid ('login' or 'bind')"

	ErrCodeRequestQueryParamEmailFieldNotValid    string = "err-400-rqpefnv"
	ErrMessageRequestQueryParamEmailFieldNotValid string = "query string 'email' in not valid"

	ErrCodeRequestHeaderCookieRefreshTokenFieldUnauthorized    string = "err-401-rhcrtfu"
	ErrMessageRequestHeaderCookieRefreshTokenFieldUnauthorized string = "refresh token unauthorized, expired or undefined"

	ErrCodeRequestRecordAlreadyExistWhenOAuthBinding    string = "err-409-rraewoab"
	ErrMessageRequestRecordAlreadyExistWhenOAuthBinding string = "OAuth info already exist"

	ErrCodeRequestOAuthAccessTokenGotError string = "err-500-roatge"
)
