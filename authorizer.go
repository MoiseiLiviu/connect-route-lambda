package main

import (
	"errors"
	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type Authorizer interface {
	Execute(token string) error
}

type authorizer struct {
	jwks *keyfunc.JWKS
}

func NewAuthorizer(jwksURL string) (Authorizer, error) {
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		log.Err(err).Msg("Failed to create JWKS from resource at the given URL")
		return nil, err
	}

	return &authorizer{jwks: jwks}, nil
}

func (a *authorizer) Execute(tokenString string) error {
	token, err := jwt.Parse(tokenString, a.jwks.Keyfunc)
	if err != nil {
		log.Err(err).Msg("Failed to parse token")
		return err
	}
	if !token.Valid {
		log.Err(err).Msg("Token is not valid")
		return errors.New("token is not valid")
	}

	return nil
}
