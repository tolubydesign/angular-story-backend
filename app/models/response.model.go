package models

type JSONResponse struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Story `json:"data"`
	Message string  `json:"message"`
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
