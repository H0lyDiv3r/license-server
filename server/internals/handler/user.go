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
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{Message: "successfull signedup"})
}

func (h *UserHandler) SigninHandler(w http.ResponseWriter, r *http.Request) {

	var req domain.SigninRequest
	json.NewDecoder(r.Body).Decode(&req)

	token, err := h.service.Signin(r.Context(), &req)
	if err != nil {
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(struct {
		Token string `json:"token"`
	}{Token: token})
}
