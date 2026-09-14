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

	"github.com/dirtydriver/devops-homework/app/internal/api"
)

func main() {

	environment := getEnv("ENVIRONMENT", "devops-homework")
	port := getEnv("PORT", "8080")

	server := api.NewServer()

	httpserver := &http.Server{
		Addr:              ":" + port,
		Handler:           server.Handler(environment),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf(
			"Starting server on port %s on %s environment",
			port,
			environment,
		)

		if err := httpserver.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}

	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	if err := httpserver.Shutdown(shutdownctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped gracefully")
}

func getEnv(env string, default_val string) string {
	value := os.Getenv(env)
	if value == "" {
		return default_val
	}

	return value
}
