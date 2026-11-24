// Command cyrus-lyrics-api is the backend for the CyrusLyrics iOS app.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alkrauss48/cyrus-lyrics-api/oauth"
	"github.com/alkrauss48/cyrus-lyrics-api/public"
	"github.com/alkrauss48/cyrus-lyrics-api/sheets"
	"github.com/gin-gonic/gin"
)

const (
	addr            = ":8000"
	shutdownTimeout = 10 * time.Second
)

func initRouter() *gin.Engine {
	router := gin.New()
	router.Use(
		// Query strings can hold OAuth codes and legacy-auth tokens, so
		// don't write them to the request log.
		gin.LoggerWithConfig(gin.LoggerConfig{SkipQueryString: true}),
		gin.Recovery(),
	)

	// Public, unauthenticated routes
	router.GET("/", public.Root)
	router.GET("/sheets/default", public.Default)

	// OAuth-driven routes
	router.GET("/oauth/google", oauth.Login)
	router.GET("/oauth/google/callback", oauth.Callback)

	// Authenticated Google Sheets driven routes
	router.GET("/sheets", sheets.Index)
	router.POST("/sheets", sheets.Create)
	router.GET("/sheets/:id", sheets.Show)
	router.DELETE("/sheets/:id", sheets.Delete)

	return router
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:              addr,
		Handler:           initRouter(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", addr)
		serveErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	case <-ctx.Done():
		log.Print("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown: %v", err)
		}
	}
}
