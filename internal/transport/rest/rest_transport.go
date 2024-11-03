package rest_transport

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"encoding/hex"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const OAuthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo"
const OAuthRedirectUrl = "http://localhost:8989/auth/callback"

type userRequest struct {
	UserHash  string
	TrackerId int64
	Url       string
	XPath     string
}

func OAuthCfg() *oauth2.Config {
	err := godotenv.Load("../../internal/transport/enviroment.env")
	if err != nil {
		log.Fatal(err)
	}
	var oAuthCfg = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		RedirectURL:  OAuthRedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	return oAuthCfg
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
	log.Println("Returned trackers", response.Data)
}

func DecodeRequest(r *http.Request) userRequest {
	var req userRequest
	parameters := r.URL.Query()
	userHash, trackerId, url, xPath := parameters.Get("UserHash"), parameters.Get("TrackerId"), parameters.Get("Url"), parameters.Get("XPath")
	if userHash == "" {
		log.Fatal("Empty token from get request")
	}
	req.UserHash = userHash
	if trackerId != "" {
		log.Println("Got tracker id in request")
		req.TrackerId, _ = strconv.ParseInt(trackerId, 10, 64)
	} else if url != "" {
		if xPath != "" {
			log.Println("Got link and xPath request")
			req.Url, req.XPath = url, xPath
		} else {
			log.Fatal("Emty link or xPath")
		}
	}
	return req
}
