package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type CognitoConfig struct {
	UserPoolID      string
	AppClientID     string
	AppClientSecret string
	AppPoolDomain   string
}

type CognitoService struct {
	client *cognito.Client
	Config *CognitoConfig
}

func NewCognitoService(cognitoConfig CognitoConfig) *CognitoService {
	config, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	client := cognito.NewFromConfig(config)

	return &CognitoService{
		client: client,
		Config: &cognitoConfig,
	}
}

func generateSecretHash(username string, config *CognitoConfig) string {
	h := hmac.New(sha256.New, []byte(config.AppClientSecret))

	_, err := h.Write([]byte(username + config.AppClientID))

	if err != nil {
		return ""
	}

	secretHash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return secretHash
}

func (service *CognitoService) AuthenticateUser(email, password string) (*cognito.InitiateAuthOutput, error) {
	input := &cognito.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		ClientId: aws.String(service.Config.AppClientID),
		AuthParameters: map[string]string{
			"USERNAME":    email,
			"PASSWORD":    password,
			"SECRET_HASH": generateSecretHash(email, service.Config),
		},
	}

	return service.client.InitiateAuth(context.TODO(), input)
}

func (service *CognitoService) RefreshToken(email, refreshToken string) (*cognito.InitiateAuthOutput, error) {
	input := &cognito.InitiateAuthInput{
		AuthFlow: "REFRESH_TOKEN_AUTH",
		ClientId: aws.String(service.Config.AppClientID),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
			"SECRET_HASH":   generateSecretHash(email, service.Config),
		},
	}

	return service.client.InitiateAuth(context.TODO(), input)
}

func (service *CognitoService) GetUserByToken(accessToken string) (*cognito.GetUserOutput, error) {
	input := &cognito.GetUserInput{
		AccessToken: aws.String(accessToken),
	}

	result, err := service.client.GetUser(context.TODO(), input)

	if err != nil {
		return nil, err
	}

	return result, nil
}
