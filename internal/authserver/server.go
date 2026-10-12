package authserver

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const sessionTTL = 10 * time.Minute

type SessionStatus string

const (
	StatusPending  SessionStatus = "pending"
	StatusComplete SessionStatus = "complete"
	StatusError    SessionStatus = "error"
)

type StartAuthResponse struct {
	SessionID string `json:"session_id"`
	AuthURL   string `json:"auth_url"`
}

type PollAuthResponse struct {
	Status SessionStatus `json:"status"`
	Token  *oauth2.Token `json:"token,omitempty"`
	Error  string        `json:"error,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type authSession struct {
	status    SessionStatus
	token     *oauth2.Token
	errMsg    string
	createdAt time.Time
}

type Server struct {
	auth     *spotifyauth.Authenticator
	mu       sync.Mutex
	sessions map[string]*authSession
}

func NewServer(redirectURI string) *Server {
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserModifyPlaybackState,
		),
	)

	s := &Server{
		auth:     auth,
		sessions: make(map[string]*authSession),
	}

	go s.cleanupLoop()
	return s
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/start", s.handleStartAuth)
	mux.HandleFunc("GET /callback", s.handleCallback)
	mux.HandleFunc("GET /auth/token", s.handlePollToken)
	mux.HandleFunc("POST /auth/refresh", s.handleRefreshToken)
	return mux
}

func (s *Server) handleStartAuth(w http.ResponseWriter, r *http.Request) {
	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, "failed to generate session", http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	s.sessions[sessionID] = &authSession{
		status:    StatusPending,
		createdAt: time.Now(),
	}
	s.mu.Unlock()

	resp := StartAuthResponse{
		SessionID: sessionID,
		AuthURL:   s.auth.AuthURL(sessionID),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode start auth response: %v", err)
	}
}

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state == "" {
		http.Error(w, "missing state parameter", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	session, exists := s.sessions[state]
	if !exists || time.Since(session.createdAt) > sessionTTL {
		delete(s.sessions, state)
		s.mu.Unlock()
		http.Error(w, "invalid or expired auth session", http.StatusBadRequest)
		return
	}
	s.mu.Unlock()

	tok, err := s.auth.Token(r.Context(), state, r)
	if err != nil {
		log.Printf("failed to exchange token for session %s: %v", state, err)
		s.mu.Lock()
		if sess, ok := s.sessions[state]; ok {
			sess.status = StatusError
			sess.errMsg = "failed to exchange authorization code for token"
		}
		s.mu.Unlock()
		http.Error(w, "Authentication failed. Please return to the app and try again.", http.StatusForbidden)
		return
	}

	s.mu.Lock()
	if sess, ok := s.sessions[state]; ok {
		sess.status = StatusComplete
		sess.token = tok
	}
	s.mu.Unlock()

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><p>Authentication complete! You can now close this window.</p></body></html>")
}

func (s *Server) handlePollToken(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "missing session_id parameter", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	session, exists := s.sessions[sessionID]
	if !exists || time.Since(session.createdAt) > sessionTTL {
		delete(s.sessions, sessionID)
		s.mu.Unlock()
		http.Error(w, "session not found or expired", http.StatusNotFound)
		return
	}

	resp := PollAuthResponse{
		Status: session.status,
		Token:  session.token,
		Error:  session.errMsg,
	}

	if session.status == StatusComplete || session.status == StatusError {
		delete(s.sessions, sessionID)
	}
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if resp.Status == StatusPending {
		w.WriteHeader(http.StatusAccepted)
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode poll response: %v", err)
	}
}

func (s *Server) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		http.Error(w, "invalid refresh request", http.StatusBadRequest)
		return
	}

	expiredToken := &oauth2.Token{
		RefreshToken: req.RefreshToken,
	}

	newToken, err := s.auth.RefreshToken(r.Context(), expiredToken)
	if err != nil {
		log.Printf("failed to refresh token: %v", err)
		http.Error(w, "failed to refresh token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newToken); err != nil {
		log.Printf("failed to encode refreshed token: %v", err)
	}
}

func (s *Server) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, sess := range s.sessions {
			if now.Sub(sess.createdAt) > sessionTTL {
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()
	}
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
