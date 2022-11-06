package account

const (
	errCodeRequestPayloadEmailFieldNotValid    = "err-400-rpefnv"
	errMessageRequestPayloadEmailFieldNotValid = "email field is not valid"

	errCodeRequestPayloadPasswordFieldNotValid    = "err-400-rppfnv"
	errMessageRequestPayloadPasswordFieldNotValid = "password field is not valid"

	errCodeRequesyPayloadOriginalPasswordFieldMismatch    = "err-400-rpopfm"
	errMessageRequesyPayloadOriginalPasswordFieldMismatch = "original password mismatch"

	errCodeRequestPayloadEmailFieldDatabaseRecordNotFound    = "err-400-rpefdbrnf"
	errMessageRequestPayloadEmailFieldDatabaseRecordNotFound = "email record not found"

	errCodeRequestPayloadPasswordFieldCompareMismatch    = "err-400-rppfcm"
	errMessageRequestPayloadPasswordFieldCompareMismatch = "password mismatch"

	errCodeRequestPayloadCaptchaFieldCompareMismatch    = "err-400-rpcfcm"
	errMessageRequestPayloadCaptchaFieldCompareMismatch = "capcha uuid or capcha code incorrect"

	errCodeRequestPayloadTokenCodeFieldsNotValid    = "err-400-eptcfnv"
	errMessageRequestPayloadTokenCodeFieldsNotValid = "auth token or code is not valid"

	errCodeRequestPayloadVerificationTypeFieldNotValid    = "err-400-rpvtnv"
	errMessageRequestPayloadVerificationTypeFieldNotValid = "verification type is not valid `[email]`"

	errCodeRequestParamAccountUUIDNotFound    = "err-400-rpaunf"
	errMessageRequestParamAccountUUIDNotFound = "account uuid record not found"

	errCodeRequestPayloadEmailFieldDatabaseRecordAlreadyExist    = "err-409-rpefdbrae"
	errMessageRequestPayloadEmailFieldDatabaseRecordAlreadyExist = "email record already exist"
)
