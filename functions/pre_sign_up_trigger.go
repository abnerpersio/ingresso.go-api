package main

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"ingresso.go/internal/infra/config"
	"ingresso.go/internal/infra/services"
)

var externalTriggerSource = "PreSignUp_ExternalProvider"

func handler(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignupResponse, error) {
	service := services.NewCognitoService(services.CognitoConfig{
		UserPoolID:      config.GetEnv("COGNITO_USER_POOL_ID"),
		AppClientID:     config.GetEnv("COGNITO_APP_CLIENT_ID"),
		AppClientSecret: config.GetEnv("COGNITO_APP_CLIENT_SECRET"),
		AppPoolDomain:   config.GetEnv("COGNITO_APP_POOL_DOMAIN"),
	})

	response := events.CognitoEventUserPoolsPreSignupResponse{
		AutoConfirmUser: true,
		AutoVerifyEmail: true,
	}

	if event.TriggerSource != externalTriggerSource {
		return response, nil
	}

	user, err := service.GetUserByEmail(event.UserName)

	if err != nil {
		user, _ = service.CreateUser(services.CreateUserInput{
			UserPoolId:   event.UserPoolID,
			Email:        event.Request.UserAttributes["email"],
			FirstName:    event.Request.UserAttributes["given_name"],
			LastName:     event.Request.UserAttributes["family_name"],
			ProfileImage: "",
		})
	}

	var attributes = make(map[string]any)
	for _, attr := range user.Attributes {
		if attr.Name != nil && attr.Value != nil {
			attributes[*attr.Name] = *attr.Value
		}
	}

	nativeUserId, ok := attributes["sub"].(string)

	if !ok {
		return response, nil
	}

	providerName := strings.Split(event.TriggerSource, "_")[0]
	providerUserId := strings.Split(event.UserName, "_")[1]

	service.LinkProvider(services.LinkProviderInput{
		UserPoolId:     event.UserPoolID,
		NativeUserId:   nativeUserId,
		ProviderName:   providerName,
		ProviderUserId: providerUserId,
	})

	return response, nil
}

func main() {
	lambda.Start(handler)
}
