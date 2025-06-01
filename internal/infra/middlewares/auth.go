package middlewares

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"ingresso.go/internal/infra/interfaces"
	"ingresso.go/internal/infra/services"
)

type Auth struct {
	Cognito *services.CognitoService
}

var UserContextKey = "user"

func getProvider(attributes map[string]any) string {
	input, ok := attributes["identities"].(string)

	if !ok {
		return ""
	}

	var identities []any
	err := json.Unmarshal([]byte(input), &identities)

	if err != nil {
		return ""
	}

	if identity, ok := identities[0].(map[string]any); ok {
		if providerName, ok := identity["providerName"].(string); ok {
			return providerName
		}
	}

	return ""
}

func (auth *Auth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		splited := strings.Split(c.GetHeader("Authorization"), "Bearer ")

		if len(splited) < 2 {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token := strings.TrimSpace(splited[1])

		if token == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		result, err := auth.Cognito.GetUserByToken(token)

		if err != nil || result == nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var attributes = make(map[string]any)
		for _, attr := range result.UserAttributes {
			if attr.Name != nil && attr.Value != nil {
				attributes[*attr.Name] = *attr.Value
			}
		}

		c.Set(UserContextKey, interfaces.User{
			Id:       attributes["sub"].(string),
			Name:     attributes["given_name"].(string),
			Email:    attributes["email"].(string),
			Provider: getProvider(attributes),
		})

		c.Next()
	}
}
