package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	types "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
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

func generateSecretHash(email string, config *CognitoConfig) string {
	hash := hmac.New(sha256.New, []byte(config.AppClientSecret))
	hash.Write([]byte(email + config.AppClientID))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
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
		fmt.Println("Error getting user by token:", err)
		return nil, err
	}

	return result, nil
}

func (service *CognitoService) listUsers(email string, paginationToken string) (*cognito.ListUsersOutput, error) {
	input := &cognito.ListUsersInput{
		UserPoolId:      aws.String(service.Config.UserPoolID),
		AttributesToGet: []string{"email"},
		Filter:          aws.String(fmt.Sprintf("email=\"%s\"", email)),
		Limit:           aws.Int32(1),
		PaginationToken: aws.String(paginationToken),
	}

	result, err := service.client.ListUsers(context.TODO(), input)

	if err != nil {
		fmt.Println("Error getting user by email:", err)
		return nil, err
	}

	if result.Users[0].Username != nil {
		return result, nil
	}

	if result.PaginationToken != nil {
		return service.listUsers(email, *result.PaginationToken)
	}

	return nil, fmt.Errorf("user with email %s not found", email)
}

func (service *CognitoService) GetUserByEmail(email string) (types.UserType, error) {
	result, err := service.listUsers(email, "")

	if err != nil {
		log.Printf("failed to list users: %v", err)
		return types.UserType{}, err
	}

	user := result.Users[0]

	if user.Username == nil {
		return types.UserType{}, fmt.Errorf("user with email %s not found", email)
	}

	return user, nil
}

type LinkProviderInput struct {
	UserPoolId     string
	NativeUserId   string
	ProviderName   string
	ProviderUserId string
}

func (service *CognitoService) LinkProvider(input LinkProviderInput) error {
	params := &cognito.AdminLinkProviderForUserInput{
		UserPoolId: aws.String(input.UserPoolId),
		DestinationUser: &types.ProviderUserIdentifierType{
			ProviderName:           aws.String("Cognito"),
			ProviderAttributeValue: aws.String(input.NativeUserId),
			ProviderAttributeName:  aws.String("Cognito_Subject"),
		},
		SourceUser: &types.ProviderUserIdentifierType{
			ProviderName:           aws.String(input.ProviderName),
			ProviderAttributeValue: aws.String(input.ProviderUserId),
			ProviderAttributeName:  aws.String("Cognito_Subject"),
		},
	}

	_, err := service.client.AdminLinkProviderForUser(context.TODO(), params)

	if err != nil {
		log.Printf("failed to link provider: %v", err)
		return err
	}

	return nil
}

type CreateUserInput struct {
	UserPoolId   string
	Email        string
	FirstName    string
	LastName     string
	ProfileImage string
}

func (service *CognitoService) CreateUser(input CreateUserInput) (types.UserType, error) {
	params := &cognito.AdminCreateUserInput{
		UserPoolId:    aws.String(input.UserPoolId),
		Username:      aws.String(input.Email),
		MessageAction: types.MessageActionTypeSuppress,
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(input.Email),
			},
			{
				Name:  aws.String("given_name"),
				Value: aws.String(input.FirstName),
			},
			{
				Name:  aws.String("family_name"),
				Value: aws.String(input.LastName),
			},
			{
				Name:  aws.String("profile_image"),
				Value: aws.String(input.ProfileImage),
			},
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"),
			},
		},
	}

	result, err := service.client.AdminCreateUser(context.TODO(), params)

	if err != nil {
		log.Printf("failed to create user: %v", err)
		return types.UserType{}, err
	}

	return *result.User, nil
}
