//go:build lambda

package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
)

type Response struct {
	Body       string `json:"body"`
	StatusCode int    `json:"statusCode"`
}

func Handler(event events.APIGatewayWebsocketProxyRequest) (Response, error) {
	authorizer, ok := event.RequestContext.Authorizer.(map[string]interface{})
	if !ok {
		log.Info().Msg("Authorizer not found in request context")
		return Response{
			StatusCode: 401,
			Body:       "Authorizer not found in request context",
		}, nil
	}

	log.Info().Msgf("Authorizer: %v", authorizer)

	userID, ok := authorizer["UserID"].(string)
	if !ok {
		log.Info().Msg("UserID not found in authorizer")
		return Response{
			StatusCode: 401,
			Body:       "UserID not found in authorizer",
		}, nil
	}

	err := Init(event.RequestContext.ConnectionID, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save connection details")
		return Response{
			StatusCode: 500,
			Body:       "Failed to save connection details",
		}, nil
	}

	return Response{
		StatusCode: 200,
		Body:       "Successfully created connection!",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
