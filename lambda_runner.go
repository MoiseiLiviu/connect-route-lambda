//go:build lambda

package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

type Response struct {
	Body       string `json:"body"`
	StatusCode int    `json:"statusCode"`
}

func Handler(event events.APIGatewayWebsocketProxyRequest) (Response, error) {
	jwksUrl := os.Getenv("JWKS_URL")
	if jwksUrl == "" {
		log.Error().Msg("JWKS_URL is not set")
		return Response{
			Body:       "JWKS_URL is not set",
			StatusCode: 500,
		}, nil
	}

	authorizer, err := NewAuthorizer(jwksUrl)
	if err != nil {
		log.Err(err).Msg("Failed to create authorizer")
		return Response{
			Body:       "Failed to create authorizer",
			StatusCode: 500,
		}, nil
	}

	authHeader := event.Headers["Authorization"]
	if authHeader == "" {
		log.Error().Msg("Authorization header not found")
		return Response{
			Body:       "Authorization header not found",
			StatusCode: 401,
		}, nil
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		log.Error().Msg("Token not found")
		return Response{
			Body:       "Token not found",
			StatusCode: 401,
		}, nil
	}

	err = authorizer.Execute(token)
	if err != nil {
		log.Err(err).Msg("invalid token")
		return Response{
			Body:       err.Error(),
			StatusCode: 401,
		}, nil
	}

	return Response{
		Body:       "Authorized",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
