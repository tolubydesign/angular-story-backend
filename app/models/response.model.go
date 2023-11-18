package models

type HTTPResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   bool        `json:"error"`
}

type ErrorResponse struct {
	ErrorMessage interface{} `json:"errorMessage"`
	Code         int         `json:"code"`
}

type HTTPError struct {
	ErrorMessage interface{} `json:"errorMessage"`
}

type Response struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseMessage struct {
	NilClient        string
	InvalidUUID      string
	InvalidCreatorID string
}
