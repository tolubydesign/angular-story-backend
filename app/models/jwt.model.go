package models

import "github.com/golang-jwt/jwt/v5"

type UserClaim struct {
	Id           string `json:"id" validate:"uuid"`
	Email        string `json:"email" validate:"required"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Username     string `json:"username" validate:"required"`
	AccountLevel string `json:"accountLevel" validate:"required"`
	jwt.RegisteredClaims
}
