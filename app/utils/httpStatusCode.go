package utils

import (
	"github.com/tolubydesign/angular-story-backend/app/models"
)

// Get the relevant http response code and status.
// Resource: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#server_error_responses
func GetHTTPResponseStatusCode(code int) models.HTTPResponseStatusCode {
	var status models.HTTPResponseStatusCode

	switch code {
	case 403:
		// The client does not have access rights to the content; that is, it is unauthorized, so the server is
		// refusing to give the requested resource. Unlike 401 Unauthorized, the client's identity is known to the server.
		status = models.HTTPResponseStatusCode{
			Code:   code,
			Status: "Forbidden",
		}
	case 400:
		// The server cannot find the requested resource. In the browser, this means the URL is not recognized.
		// In an API, this can also mean that the endpoint is valid but the resource itself does not exist.
		// Servers may also send this response instead of 403 Forbidden to hide the existence of a resource from an unauthorized client.
		// This response code is probably the most well known due to its frequent occurrence on the web.
		status = models.HTTPResponseStatusCode{
			Code:   code,
			Status: "Bad Request",
		}
	case 404:
		// The server cannot find the requested resource. In the browser, this means the URL is not recognized.
		// In an API, this can also mean that the endpoint is valid but the resource itself does not exist.
		// Servers may also send this response instead of 403 Forbidden to hide the existence of a resource from an unauthorized client.
		// This response code is probably the most well known due to its frequent occurrence on the web.
		status = models.HTTPResponseStatusCode{
			Code:   code,
			Status: "Not Found",
		}
	case 200:
		status = models.HTTPResponseStatusCode{
			Code:   code,
			Status: "OK",
		}
	default:
		// The server has encountered a situation it does not know how to handle.
		status = models.HTTPResponseStatusCode{
			Code:   500,
			Status: "Internal Server Error",
		}
	}

	return status
}
