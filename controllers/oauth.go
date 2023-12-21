package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"strings"
)

type OAuth struct {
	ProviderConfigs map[string]*oauth2.Config
}

func (oa OAuth) Connect(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)
	config, ok := oa.ProviderConfigs[provider]
	if !ok {
		http.Error(w, "Invalid oauth provider", http.StatusBadRequest)
		return
	}

	state := csrf.Token(r)
	setCookie(w, "oauth_state", state)

	// count be added as an extra layer of security
	//verifier := oauth2.GenerateVerifier()
	//oauth2.S256ChallengeOption(verifier)

	url := config.AuthCodeURL(
		state,
		//oauth2.SetAuthURLParam("token_access_type", "offline"),
		oauth2.SetAuthURLParam("redirect_uri", redirectUri(r, provider)))

	http.Redirect(w, r, url, http.StatusFound)
}

func (oa OAuth) Callback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)
	config, ok := oa.ProviderConfigs[provider]
	if !ok {
		http.Error(w, "Invalid oauth provider", http.StatusBadRequest)
		return
	}

	state := r.FormValue("state")
	cookieState, err := readCookie(r, "oauth_state")
	if err != nil || cookieState != state {
		if err != nil {
			log.Println(err)
		}
		http.Error(w, "Invalid request 1", http.StatusBadRequest)
		return
	}

	//deleteCookie(w, "oauth_state")

	code := r.FormValue("code")
	token, err := config.Exchange(r.Context(), code, oauth2.SetAuthURLParam("redirect_uri", redirectUri(r, provider)))
	if !ok {
		http.Error(w, "Something went wrong 1", http.StatusBadRequest)
		return
	}
	// Persist token
	// Redirect user to where they started the process (could store the url in context or state)

	client := config.Client(r.Context(), token)
	res, err := client.Post(
		"https://api.dropboxapi.com/2/files/list_folder",
		"application/json",
		strings.NewReader(`
		{
			"path": ""
		}
	`))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	defer res.Body.Close()

	rawBytes, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, rawBytes, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	pretty.WriteTo(w)
}

func redirectUri(r *http.Request, provider string) string {
	if r.Host == "localhost:8080" {
		return fmt.Sprintf("http://localhost:8080/oauth/%s/callback", provider)
	}

	return fmt.Sprintf("http://photosharer.renect.co.uk/oauth/%s/callback", provider)
}
