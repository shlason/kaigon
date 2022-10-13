package account

const (
	errCodeRequestPayloadEmailFieldNotValid    = "err-400-rpefnv"
	errMessageRequestPayloadEmailFieldNotValid = "email field is not valid"

	errCodeRequestPayloadPasswordFieldNotValid    = "err-400-rppfnv"
	errMessageRequestPayloadPasswordFieldNotValid = "password field is not valid"

	errCodeRequestPayloadEmailFieldDatabaseRecordNotFound    = "err-400-rpefdbrnf"
	errMessageRequestPayloadEmailFieldDatabaseRecordNotFound = "email record not found"

	errCodeRequestPayloadPasswordFieldCompareMismatch    = "err-400-rppfcm"
	errMessageRequestPayloadPasswordFieldCompareMismatch = "password mismatch"

	errCodeRequestPayloadCaptchaFieldCompareMismatch    = "err-rpcfcm"
	errMessageRequestPayloadCaptchaFieldCompareMismatch = "capcha uuid or capcha code incorrect"
)
