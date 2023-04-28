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

// HTTPError struct to HTTPError object.
type HTTPError struct {
	ErrorMessage interface{} `json:"errorMessage"`
}
