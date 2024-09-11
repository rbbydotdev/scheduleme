package server

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"scheduleme/config"
	"scheduleme/services"
	"scheduleme/sqlite"
	"syscall"
)

// Easy initializes default everything when no inversion of control is needed
func RunMain() {
	address := Run(nil, nil)
	log.Printf("Server is listening on %s", address)

	// Set up channel on which to send signal notifications.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Block until a signal is received.
	sig := <-sigChan
	log.Printf("Received signal: %s. Shutting down...", sig)
}

// Runs http routed services on config.Port returns the address string
func Run(cfg *config.ConfigStruct, db *sqlite.Db) string {
	var err error
	if cfg == nil {
		cfg = config.InitConfig()
	}
	if db == nil {
		db, err = sqlite.NewOpenDB(cfg.Dsn)
		if err != nil {
			log.Fatal(err)
		}
	}
	topSvc := services.TopServices(cfg, db)
	topRte := services.TopRoutes(topSvc)
	return ServerFromHandler(topRte, cfg.Port)
}

// startServer starts a new HTTP server on an available port and returns the listening address
func ServerFromHandler(handler http.Handler, port string) string {
	server := &http.Server{Handler: handler}
	listener, err := net.Listen("tcp", ":"+port) // Listen on any available port
	if err != nil {
		log.Fatalf("Failed to listen on a port: %v", err)
	}

	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()
	return listener.Addr().String() // Returns the address including the port the server is listening on
}
