package auth

import (
	"encoding/json"
	"net/http"

	"ingresso.go/infra/services/responses"
)

func (auth *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responses.SendError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	authResp, err := auth.Cognito.AuthenticateUser(body.Username, body.Password)
	if err != nil {
		responses.SendError(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(authResp.AuthenticationResult)
}
