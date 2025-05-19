package handlers

import (
	"net/http"

	"ingresso.go/infra/services/responses"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	responses.SendSuccess(w, responses.ResponseData{
		Message: "ok",
	}, http.StatusOK)
}
