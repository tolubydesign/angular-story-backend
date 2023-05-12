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

// HTTPError struct to HTTPError object.
type HTTPError struct {
	ErrorMessage interface{} `json:"errorMessage"`
}

type HTTPResponseStatusCode struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}
