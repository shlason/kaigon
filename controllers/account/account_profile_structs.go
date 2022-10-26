package account

type getProfileResponsePayload struct {
	Avatar    string `json:"avatar"`
	Banner    string `json:"banner"`
	Signature string `json:"signature"`
}
