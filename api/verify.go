package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
)

var (
	verifier = emailverifier.NewVerifier().EnableDomainSuggest()
)

type VerificationRequest struct {
	Email string `json:"email"`
}

type VerificationResponse struct {
	Email          string `json:"email"`
	IsValid        bool   `json:"is_valid"`
	Status         string `json:"status"`
	Disposable     bool   `json:"is_disposable"`
	RoleAccount     bool   `json:"is_role_account"`
	HasMxRecords   bool   `json:"has_mx_records"`
	Suggestion     string `json:"suggestion"`
	SyntaxValid    bool   `json:"syntax_valid"`
	GravatarFound  bool   `json:"gravatar_found"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// 1. Set CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-KEY")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Security: API Key check
	apiKey := r.Header.Get("X-API-KEY")
	masterKey := os.Getenv("MASTER_API_KEY")
	
	if masterKey != "" && apiKey != masterKey {
		http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
		return
	}

	// 3. Parse Request
	var req VerificationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 4. Verification Logic (NO SMTP HANDSHAKE as requested)
	email := strings.TrimSpace(req.Email)
	res, err := verifier.Verify(email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Verification failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Check Gravatar (Simulated/Optional check as per request)
	gravatarFound := false
	if res.Gravatar != nil {
		gravatarFound = true
	}

	// Determine overall validity
	// We consider it valid if it has MX records, syntax is correct, and it's not disposable
	isValid := res.Syntax.Valid && res.HasMXRecords && !res.Disposable

	response := VerificationResponse{
		Email:         email,
		IsValid:       isValid,
		Status:        getFriendlyStatus(isValid, res.Disposable, res.HasMXRecords),
		Disposable:    res.Disposable,
		RoleAccount:    res.RoleAccount,
		HasMxRecords:  res.HasMXRecords,
		Suggestion:    verifier.SuggestDomain(res.Syntax.Domain),
		SyntaxValid:   res.Syntax.Valid,
		GravatarFound: gravatarFound,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getFriendlyStatus(isValid bool, isDisposable bool, hasMx bool) string {
	if isDisposable {
		return "Disposable"
	}
	if !hasMx {
		return "No MX Records"
	}
	if isValid {
		return "Valid"
	}
	return "Invalid"
}
