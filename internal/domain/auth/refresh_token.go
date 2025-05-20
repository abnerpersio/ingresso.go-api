package auth

import (
	"encoding/json"
	"net/http"

	"ingresso.go/internal/infra/services/responses"
)

func (auth *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email        string `json:"email"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responses.SendError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	authResp, err := auth.Cognito.RefreshToken(body.Email, body.RefreshToken)

	if err != nil {
		responses.SendError(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(authResp.AuthenticationResult)
}
