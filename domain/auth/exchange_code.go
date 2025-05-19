package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"ingresso.go/infra/services/responses"
)

func (auth *AuthHandler) ExchangeCode(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("code", queryParams.Get("code"))
	params.Set("redirect_uri", queryParams.Get("redirect_uri"))

	fmt.Println("params", params)

	url := auth.Cognito.Config.AppPoolDomain + "/oauth2/token"
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(params.Encode()))
	req.SetBasicAuth(auth.Cognito.Config.AppClientID, auth.Cognito.Config.AppClientSecret)

	if err != nil {
		responses.SendError(w, "Failed to exchange code", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		responses.SendError(w, "Failed to exchange code", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var responseBody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		responses.SendError(w, "Failed to exchange code", http.StatusInternalServerError)
		return
	}

	responses.SendSuccess(w, responses.ResponseData{
		Data: map[string]any{
			"accessToken":  responseBody.AccessToken,
			"refreshToken": responseBody.RefreshToken,
		},
	}, http.StatusOK)
}
