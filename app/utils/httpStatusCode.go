package utils

import (
	"github.com/tolubydesign/angular-story-backend/app/models"
)

// Get the relevant http response code and status.
// HTTP Status Response - https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
func GetHTTPResponseStatusCode(code int) models.Response {
	var status models.Response

	switch code {
	case 200:
		status = models.Response{
			Status:  "OK",
			Message: "The request has succeeded. An entity corresponding to the requested resource is sent in the response.",
		}

	case 201:
		status = models.Response{
			Status:  "CREATED",
			Message: "The request has been fulfilled and resulted in a new resource being created.",
		}

	case 204:
		status = models.Response{
			Status:  "NO CONTENT",
			Message: "The server successfully processed the request, but is not returning any content.",
		}

	case 400:
		// The server cannot or will not process the request due to something that is perceived to be a client error
		// (e.g., malformed request syntax, invalid request message framing, or deceptive request routing).
		// INVALID REQUEST
		status = models.Response{
			Status:  "BAD REQUEST",
			Message: "The request could not be understood by the server due to malformed syntax or request.",
		}

	case 403:
		// The client does not have access rights to the content; that is, it is unauthorized, so the server is
		// refusing to give the requested resource. Unlike 401 Unauthorized, the client's identity is known to the server.
		status = models.Response{
			Code:   code,
			Status: "FORBIDDEN",
		}

	case 404:
		// The server cannot find the requested resource. In the browser, this means the URL is not recognized.
		// In an API, this can also mean that the endpoint is valid but the resource itself does not exist.
		// Servers may also send this response instead of 403 Forbidden to hide the existence of a resource from an unauthorized client.
		// This response code is probably the most well known due to its frequent occurrence on the web.
		status = models.Response{
			Status:  "NOT FOUND",
			Message: "The server has not found anything matching the request.Boron",
		}

	case 500:
		// The server has encountered a situation it does not know how to handle.
		status = models.Response{
			Status:  "INTERNAL SERVER ERROR",
			Message: "The server encountered an unexpected condition which prevented it from fulfilling the request.",
		}

	default:
		status = models.Response{
			Status:  "INVALID REQUEST",
			Message: "The request could not be understood by the server due to malformed syntax or request.",
		}
	}

	status.Code = code
	return status
}

// Return response code base on status code provided
func HandleResponseMessage(code int) (res models.Response) {
	responses := map[int]models.Response{
		200: {
			Status:  "OK",
			Message: "The request has succeeded. An entity corresponding to the requested resource is sent in the response.",
		},
		201: {
			Status:  "CREATED",
			Message: "The request has been fulfilled and resulted in a new resource being created.",
		},
		204: {
			Status:  "NO CONTENT",
			Message: "The server successfully processed the request, but is not returning any content.",
		},
		400: {
			Status:  "INVALID REQUEST",
			Message: "The request could not be understood by the server due to malformed syntax or request.",
		},
		404: {
			Status:  "NOT FOUND",
			Message: "The server has not found anything matching the request.Boron",
		},
		500: {
			Status:  "INTERNAL SERVER ERROR",
			Message: "The server encountered an unexpected condition which prevented it from fulfilling the request.",
		},
	}

	response := responses[code]
	return response
}
