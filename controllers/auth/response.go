package auth

const (
	ErrCodeRequestQueryParamAccountUUIDFieldNotValid    string = "err-400-rqpaufnv"
	ErrMessageRequestQueryParamAccountUUIDFieldNotValid string = "query string 'accountUuid' in not valid"

	ErrCodeRequestQueryParamEmailFieldNotValid    string = "err-400-rqpefnv"
	ErrMessageRequestQueryParamEmailFieldNotValid string = "query string 'email' in not valid"

	ErrCodeRequestHeaderCookieRefreshTokenFieldUnauthorized    string = "err-401-rhcrtfu"
	ErrMessageRequestHeaderCookieRefreshTokenFieldUnauthorized string = "refresh token unauthorized, expired or undefined"
)
