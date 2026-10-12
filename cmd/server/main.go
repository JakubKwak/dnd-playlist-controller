package main

import (
	"context"
	"dnd-playlist-controller/internal/authserver"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	if os.Getenv("SPOTIFY_ID") == "" || os.Getenv("SPOTIFY_SECRET") == "" {
		log.Fatal("(╯°□°）╯︵ ┻━┻ SPOTIFY_ID and SPOTIFY_SECRET environment variables must be set")
	}

	redirectURI := os.Getenv("REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/callback"
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		if port := os.Getenv("PORT"); port != "" {
			listenAddr = ":" + port
		} else {
			listenAddr = ":8080"
		}
	}

	authSrv := authserver.NewServer(redirectURI)
	httpSrv := &http.Server{
		Addr:              listenAddr,
		Handler:           authSrv.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("Auth server listening on %s (redirect URI: %s)", listenAddr, redirectURI)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("(╯°□°）╯︵ ┻━┻ Failed to start server: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down auth server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatalf("(╯°□°）╯︵ ┻━┻ Server shutdown failed: %v", err)
	}
}
