package handler

import (
	"encoding/json"
	"license-server/internals/domain"
	"license-server/internals/service"
	"net/http"
)

type UserHandler struct {
	service service.UserService
}

func NewHandler(service service.UserService) UserHandler {
	return UserHandler{service: service}
}

func (h *UserHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req domain.SignupRequest

	json.NewDecoder(r.Body).Decode(&req)
	err := h.service.Signup(r.Context(), &req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Write([]byte("user signed up successfully"))
}

func (h *UserHandler) SigninHandler(w http.ResponseWriter, r *http.Request) {

	var req domain.SigninRequest
	json.NewDecoder(r.Body).Decode(&req)

	token, err := h.service.Signin(r.Context(), &req)
	if err != nil {
		http.Error(w, "failed to login", http.StatusInternalServerError)
	}

	w.Write([]byte(token))
}
