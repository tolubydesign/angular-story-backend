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

type APIEndpoints struct {
	Post   PostEndpoint   `json:"post"`
	Put    PutEndpoint    `json:"put"`
	Get    GetEndpoint    `json:"get"`
	Delete DeleteEndpoint `json:"delete"`
}

type GetEndpoint struct {
	Story       string `json:"story"`
	AllStories  string `json:"allStories"`
	HealthCheck string `json:"healthCheck"`
	Login       string `json:"login"`
	Users       string `json:"users"`
	Tables      string `json:"tables"`
}

type PostEndpoint struct {
	Story            string `json:"story"`
	PopulateDatabase string `json:"populateDatabase"` // Not for production
	SignUp           string `json:"signUp"`
}

type PutEndpoint struct {
	Story string `json:"story"`
}

type DeleteEndpoint struct {
	Story string `json:"story"`
}
