package spotifyclient

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth  *spotifyauth.Authenticator
	ch    = make(chan *spotify.Client)
	state = "jk-dnd-playlist-controller" // love security n shit
)

func New() (*spotify.Client, error) {
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserModifyPlaybackState))

	// try to load an existing user from token.json if it exists
	existingClient, err := clientFromExistingToken()
	if err == nil && existingClient != nil {
		return existingClient, nil
	}
	fmt.Printf("Couldn't get an existing client :( have to make new one: %s \n", err)

	return clientFromNewToken()
}

func clientFromExistingToken() (*spotify.Client, error) {
	token, err := readToken()
	if err != nil {
		return nil, err
	}
	client := spotify.New(auth.Client(context.Background(), token))

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println("Spotify user ID:", user.ID)

	return client, nil
}

func clientFromNewToken() (*spotify.Client, error) {
	// To get the auth we need to allow spotify to redirect the user to our server, which is very cringe btw
	srv := &http.Server{Addr: ":8080", Handler: nil}
	http.HandleFunc("/callback", completeAuth)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Spotify login link:", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println("Spotify user ID:", user.ID)

	return client, nil
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// save the token so we can skip auth next time
	err = saveToken(tok)
	if err != nil {
		fmt.Printf("couldn't save the token: %s", err)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	ch <- client

	// output message to user
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><p>Authentication complete! You can now close this window.</p></body></html>")
}
