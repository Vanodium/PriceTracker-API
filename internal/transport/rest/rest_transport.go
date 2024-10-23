package rest_transport

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"encoding/hex"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const OauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo"

type userRequest struct {
	UserHash  string
	TrackerId int64
	Url       string
	XPath     string
}

func OauthCfg() *oauth2.Config {
	err := godotenv.Load("../../internal/transport/oauthcfg.env")
	if err != nil {
		log.Fatal(err)
	}
	var OauthCfg = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8989/auth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	return OauthCfg
}

func HashString(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ApiResponceJson[T any](w http.ResponseWriter, data T, isError bool, message string) {
	type ApiResponse struct {
		Data    interface{} `json:"data"`
		Error   bool        `json:"error,omitempty"`
		Message string      `json:"message,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json")
	response := ApiResponse{
		Data: data,
	}
	if isError {
		response.Error = true
		response.Message = message
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DecodeRequest(r *http.Request) userRequest {
	var req userRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.Decode(&req)
	return req
}
