package handler

import (
	"encoding/json"
	"fmt"
	"license-server/internals/domain"
	"license-server/internals/service"
	"net/http"
)

type LicenseHandler struct {
	service service.LicenseService
}

func NewLicenseHandler(service service.LicenseService) LicenseHandler {
	return LicenseHandler{service: service}
}

func (h *LicenseHandler) GenerateLicense(w http.ResponseWriter, r *http.Request) {

	license, err := h.service.GenerateKey(r.Context())
	if err != nil {
		http.Error(w, fmt.Errorf("error generating key: %s ", err.Error()).Error(), http.StatusInternalServerError)
	}

	resp, err := json.Marshal(license)

	if err != nil {
		http.Error(w, fmt.Errorf("error generating key: %s ", err.Error()).Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(resp))
}

func (h *LicenseHandler) ActivateLicense(w http.ResponseWriter, r *http.Request) {

	var req domain.ActivateLicenseRequest
	json.NewDecoder(r.Body).Decode(&req)

	token, err := h.service.ActivateLicense(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to generate a token: %s", err.Error()).Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(token))

}
