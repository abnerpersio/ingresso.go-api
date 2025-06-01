package domain_auth

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"ingresso.go/internal/infra/services"
)

func (auth *AuthHandler) ExchangeCode(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	params := url.Values{
		"grant_type":   []string{"authorization_code"},
		"code":         []string{queryParams.Get("code")},
		"redirect_uri": []string{queryParams.Get("redirect_uri")},
	}

	url := auth.Cognito.Config.AppPoolDomain + "/oauth2/token"
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(params.Encode()))
	req.SetBasicAuth(auth.Cognito.Config.AppClientID, auth.Cognito.Config.AppClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		services.SendError(c, "Failed to exchange code", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK || err != nil {
		services.SendError(c, "Failed to exchange code", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	var result map[string]string
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	if result["access_token"] == "" || result["refresh_token"] == "" {
		services.SendError(c, "Failed to exchange code", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": result["access_token"], "refreshToken": result["refresh_token"]})
}
