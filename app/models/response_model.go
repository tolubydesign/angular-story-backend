package models

type JSONResponse struct {
	Type    string  `json:"type"`
	Data    []Story `json:"data"`
	Message string  `json:"message"`
}
