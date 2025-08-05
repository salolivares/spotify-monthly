package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const (
	redirectURI = "http://127.0.0.1:8080/callback"
	tokenPath   = "./token.json"
	authState   = "spotify-monthly-cli"
)

var auth = spotifyauth.New(
	spotifyauth.WithRedirectURL(redirectURI),
	spotifyauth.WithScopes(
		spotifyauth.ScopeUserLibraryRead,
		spotifyauth.ScopePlaylistModifyPrivate,
	),
)

type persistedToken struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

// Try to load token from disk
func loadToken() (*oauth2.Token, error) {
	var pt persistedToken
	f, err := os.Open(tokenPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&pt); err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken:  pt.AccessToken,
		TokenType:    pt.TokenType,
		RefreshToken: pt.RefreshToken,
		Expiry:       pt.Expiry,
	}, nil
}

// Save token to disk
func saveToken(tok *oauth2.Token) error {
	pt := persistedToken{
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
	}
	f, err := os.Create(tokenPath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(&pt)
}

// Top-level: returns a usable Spotify client (handles refresh if needed)
func getClient(ctx context.Context) *spotify.Client {
	tok, err := loadToken()
	if err == nil {
		httpClient := auth.Client(ctx, tok)
		return spotify.New(httpClient)
	}

	// No saved token — do interactive login
	ch := make(chan *oauth2.Token)
	server := &http.Server{Addr: ":8080"}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		tok, err := auth.Token(r.Context(), authState, r)
		if err != nil {
			http.Error(w, "Auth failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Login complete. You can close this tab.")
		go server.Shutdown(context.Background())
		saveToken(tok)
		ch <- tok
	})

	url := auth.AuthURL(authState)
	fmt.Println("Open this URL in your browser to authorize:")
	fmt.Println(url)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	tok = <-ch
	httpClient := auth.Client(ctx, tok)
	return spotify.New(httpClient)
}
