package routes

import (
	"encoding/json"

	"log"
	"net/http"

	"context"

	db_operations "github.com/Vanodium/pricetracker/internal/db"
	core_functions "github.com/Vanodium/pricetracker/internal/services"
	rest_transport "github.com/Vanodium/pricetracker/internal/transport/rest"

	"golang.org/x/oauth2"
)

func Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/oauth", oAuthHandler)
	mux.HandleFunc("/auth/callback", oAuthCallbackHandler)
	mux.HandleFunc("/auth/token", getTokenHandler)
	mux.HandleFunc("/trackers", trackersHandler)
	return mux
}

func oAuthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got oauth request")
	url := rest_transport.OAuthCfg().AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

var OAuthUserData struct {
	UserHash  string
	UserEmail string `json:"email"`
}

func oAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got oauth callback")

	code := r.URL.Query().Get("code")
	tok, err := rest_transport.OAuthCfg().Exchange(context.Background(), code)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Got oauth token")

	OAuthUserData.UserHash = rest_transport.HashString(tok.RefreshToken)
	resp, err := rest_transport.OAuthCfg().Client(context.Background(), tok).Get(rest_transport.OAuthGoogleUrlAPI)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&OAuthUserData)
	if err != nil {
		log.Fatal("Unable to decode user info")
	}
	err = db_operations.AddUser(OAuthUserData.UserHash, OAuthUserData.UserEmail)
	if err != nil {
		rest_transport.ApiResponceJson(w, "", true, "Can not add user")
		log.Fatal(err)
	}
	http.Redirect(w, r, "token", http.StatusTemporaryRedirect)
}

func getTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got TOKEN request")
	rest_transport.ApiResponceJson(w, OAuthUserData.UserHash, false, "")
}

func trackersHandler(w http.ResponseWriter, r *http.Request) {
	req := rest_transport.DecodeRequest(r)
	if !db_operations.UserExists(req.UserHash) {
		log.Fatal("Wrong user token!")
		rest_transport.ApiResponceJson(w, "", true, "Wrong user token!")
	}
	userId := db_operations.GetUserId(req.UserHash)

	switch r.Method {
	case http.MethodGet:
		log.Println("Received GET request")
		trackers := db_operations.GetUserTrackers(userId)
		rest_transport.ApiResponceJson(w, trackers, false, "")
	case http.MethodPost:
		log.Println("Received POST request")
		currentPrice := core_functions.ParseDigits(core_functions.ExtractPrice(req.TrackerUrl, req.CssSelector))
		if currentPrice == "" {
			log.Fatal("Corrupted link or path")
			rest_transport.ApiResponceJson(w, "", true, "Error parsing price")
		}

		trackerId, err := db_operations.AddTracker(userId, req.TrackerUrl, req.CssSelector)
		if err != nil {
			rest_transport.ApiResponceJson(w, "", true, "Error adding tracker")
			log.Fatal(err)
		}

		currentDate := core_functions.GetCurrentDate()
		err = db_operations.AddPrice(trackerId, currentPrice, currentDate)
		if err != nil {
			rest_transport.ApiResponceJson(w, "", true, "Error adding price")
			log.Fatal(err)
		}
		rest_transport.ApiResponceJson(w, "Tracker added", false, "")

	case http.MethodDelete:
		log.Println("Received DELETE request")
		err := db_operations.DeleteTracker(userId, req.TrackerId)
		if err != nil {
			rest_transport.ApiResponceJson(w, "", true, "Error deleting tracker")
			log.Fatal(err)
		}
		rest_transport.ApiResponceJson(w, "Tracker deleted", false, "")
	default:
		log.Println("Received wrong request")
		http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
	}
}
