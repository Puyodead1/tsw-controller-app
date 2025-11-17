package tswconnector

type TSWAPI_ErrorResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type TSWAPI_GetResponse struct {
	Result  string         `json:"Result" validate:"oneof=Success,Error"`
	Values  map[string]any `json:"Values"`
	Message string         `json:"Message"`
}
