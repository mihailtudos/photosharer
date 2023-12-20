package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dropboxAppId := os.Getenv("DROPBOX_APP_ID")
	dropboxAppSecret := os.Getenv("DROPBOX_APP_SECRET")
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     dropboxAppId,
		ClientSecret: dropboxAppSecret,
		Scopes:       []string{"files.content.read", "files.metadata.read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/oauth2/authorize",
			TokenURL: "https://api.dropboxapi.com/oauth2/token",
		},
	}

	// use PKCE to protect against CSRF attacks
	// https://www.ietf.org/archive/id/draft-ietf-oauth-security-topics-22.html#name-countermeasures-6
	verifier := oauth2.GenerateVerifier()

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
	fmt.Printf("Once you have a code paste it\n")

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	res, err := client.Post("https://api.dropboxapi.com/2/files/list_folder", "application/json", strings.NewReader(`
		{
			"path": ""
		}
	`))

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	_, _ = io.Copy(os.Stdout, res.Body)
}
