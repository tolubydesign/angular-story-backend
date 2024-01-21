package helpers

import "github.com/tolubydesign/angular-story-backend/app/models"

var ResponseMessages = models.ResponseMessage{
	NilClient:        "Dynamo Database inaccessible.",
	InvalidUUID:      "Invalid ID provided.",
	InvalidCreatorID: "Invalid creator value provided.",
}
