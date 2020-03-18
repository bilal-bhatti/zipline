package main

import (
	"context"
	"flag"

	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/bilal-bhatti/zipline/example/web"
)

func main() {
	// global log settings
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	listenAddr := flag.String("listen", ":5678", "spcifiy port to listen on")
	flag.Parse()

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	// Listen for OS signals
	signal.Notify(quit, os.Interrupt)

	// Configure a new logger
	logger := log.New(os.Stdout, "SRVR: ", log.LstdFlags)

	env := &connectors.Env{}

	// Init http server
	server := &http.Server{
		Addr:         *listenAddr,
		Handler:      web.NewRouter(env),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start shutdown go routine
	go gracefullShutdown(server, logger, quit, done)

	// Ready to serve
	logger.Println("Server is ready to handle requests at", *listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", *listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}

func gracefullShutdown(server *http.Server, logger *log.Logger, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	logger.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(done)
}
