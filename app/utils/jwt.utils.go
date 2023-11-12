package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tolubydesign/angular-story-backend/app/config"
	"github.com/tolubydesign/angular-story-backend/app/models"
)

// asdf
func JWTVerification(token string) (interface{}, error) {
	err := ValidateString(token)
	if err != nil {
		return nil, err
	}

	// TODO: complete verification

	return nil, nil
}

/*
Construct JWT token based on user information

resource - https://www.sohamkamani.com/golang/jwt-authentication/ |
*/
func BuildUserJWTToken(u *models.DatabaseUserStruct) (string, error) {
	// TODO: log security event
	// Get secret key from config
	configuration, err := config.GetConfiguration()
	if err != nil {
		return "", err
	}

	secret := configuration.Configuration.JWTSecretKey

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(24 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := models.UserClaim{
		Id:           u.Id,
		Email:        u.Email,
		Name:         u.Name,
		Surname:      u.Surname,
		Username:     u.AccountLevel,
		AccountLevel: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "delta",
			Subject:   "login",
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT token string
	ts, err := token.SignedString(secret)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", err
	}

	return ts, err
}

func GetJWTSecretKey() ([]byte, error) {
	configuration, err := config.GetConfiguration()
	if err != nil {
		return nil, err
	}

	var secret = []byte(configuration.Configuration.JWTSecretKey)
	return secret, nil
}
