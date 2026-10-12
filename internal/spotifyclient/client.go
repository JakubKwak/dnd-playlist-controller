package spotifyclient

import (
	"bytes"
	"context"
	"dnd-playlist-controller/internal/authserver"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

const (
	defaultAuthServerURL = "http://localhost:8080"
	pollInterval         = 2 * time.Second
	pollTimeout          = 10 * time.Minute
	httpTimeout          = 15 * time.Second
)

func New() (*spotify.Client, error) {
	serverURL := strings.TrimRight(os.Getenv("AUTH_SERVER_URL"), "/")
	if serverURL == "" {
		serverURL = defaultAuthServerURL
	}

	// try to load an existing user from token.json if it exists
	existingClient, err := clientFromExistingToken(serverURL)
	if err == nil && existingClient != nil {
		return existingClient, nil
	}
	fmt.Printf("Couldn't get an existing client: %s \n", err)

	return clientFromNewToken(serverURL)
}

func clientFromExistingToken(serverURL string) (*spotify.Client, error) {
	token, err := readToken()
	if err != nil {
		return nil, err
	}

	client := spotify.New(newOAuthClient(context.Background(), serverURL, token))

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println("Spotify user ID:", user.ID)

	return client, nil
}

func clientFromNewToken(serverURL string) (*spotify.Client, error) {
	httpClient := &http.Client{Timeout: httpTimeout}

	resp, err := httpClient.Post(serverURL+"/auth/start", "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to contact auth server at %s: %w", serverURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth server returned status %d when starting auth", resp.StatusCode)
	}

	var startResp authserver.StartAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&startResp); err != nil {
		return nil, fmt.Errorf("failed to decode auth server response: %w", err)
	}

	fmt.Println("Spotify login link:", startResp.AuthURL)

	tok, err := pollForToken(httpClient, serverURL, startResp.SessionID)
	if err != nil {
		return nil, err
	}

	if err := saveToken(tok); err != nil {
		fmt.Printf("couldn't save the token: %s\n", err)
	}

	client := spotify.New(newOAuthClient(context.Background(), serverURL, tok))

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println("Spotify user ID:", user.ID)

	return client, nil
}

func pollForToken(httpClient *http.Client, serverURL, sessionID string) (*oauth2.Token, error) {
	deadline := time.Now().Add(pollTimeout)
	pollURL := fmt.Sprintf("%s/auth/token?session_id=%s", serverURL, url.QueryEscape(sessionID))

	for time.Now().Before(deadline) {
		time.Sleep(pollInterval)

		resp, err := httpClient.Get(pollURL)
		if err != nil {
			continue
		}

		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			return nil, fmt.Errorf("auth session expired or not found")
		}

		var pollResp authserver.PollAuthResponse
		err = json.NewDecoder(resp.Body).Decode(&pollResp)
		resp.Body.Close()
		if err != nil {
			continue
		}

		switch pollResp.Status {
		case authserver.StatusComplete:
			if pollResp.Token == nil {
				return nil, fmt.Errorf("auth completed without token")
			}
			return pollResp.Token, nil
		case authserver.StatusError:
			return nil, fmt.Errorf("authentication failed: %s", pollResp.Error)
		case authserver.StatusPending:
			continue
		}
	}

	return nil, fmt.Errorf("timed out waiting for Spotify authentication")
}

func newOAuthClient(ctx context.Context, serverURL string, token *oauth2.Token) *http.Client {
	ts := oauth2.ReuseTokenSource(token, &remoteTokenSource{
		serverURL:    serverURL,
		refreshToken: token.RefreshToken,
		httpClient:   &http.Client{Timeout: httpTimeout},
	})
	return oauth2.NewClient(ctx, ts)
}

type remoteTokenSource struct {
	mu           sync.Mutex
	serverURL    string
	refreshToken string
	httpClient   *http.Client
}

func (r *remoteTokenSource) Token() (*oauth2.Token, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.refreshToken == "" {
		return nil, fmt.Errorf("no refresh token available")
	}

	body, err := json.Marshal(authserver.RefreshRequest{
		RefreshToken: r.refreshToken,
	})
	if err != nil {
		return nil, err
	}

	resp, err := r.httpClient.Post(r.serverURL+"/auth/refresh", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token via auth server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth server returned status %d on token refresh", resp.StatusCode)
	}

	var newToken oauth2.Token
	if err := json.NewDecoder(resp.Body).Decode(&newToken); err != nil {
		return nil, fmt.Errorf("failed to decode refreshed token: %w", err)
	}

	if newToken.RefreshToken == "" {
		newToken.RefreshToken = r.refreshToken
	} else {
		r.refreshToken = newToken.RefreshToken
	}

	if err := saveToken(&newToken); err != nil {
		fmt.Printf("couldn't save refreshed token: %s\n", err)
	}

	return &newToken, nil
}
