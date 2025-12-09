package server

import (
	"fmt"
	"log"
	"sync"
	"syscall"
	"context"
	"net/http"
	"os/signal"
)


var wg sync.WaitGroup

func Server() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()
	LoadConfig(CONFIG_FILE)
	setupHandlers()
	go start()
	<-ctx.Done()	// Wait for Ctrl+C or other stop signal
	wg.Wait()		// Ensures registered processes are complete before exiting
}

func setupHandlers() {
	http.HandleFunc("/", defaultRequest)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))
	http.HandleFunc("GET /api/events", getEvents)
	http.HandleFunc("POST /api/event", getEvent)
	http.HandleFunc("POST /api/session/auth", authenticateSession)
	http.HandleFunc("POST /api/user/login", loginUser)
	http.HandleFunc("POST /api/user/logout", logoutUser)
	http.HandleFunc("POST /api/user/list", userReservations)
	http.HandleFunc("POST /api/user/reserve", makeReservation)
	http.HandleFunc("POST /api/user/register", registerUser)
	http.HandleFunc("POST /api/res/login", loginWithReservation)
	http.HandleFunc("POST /api/admin/create", createEvent)
}

func start() {
	log.Println("Server UP")
	if CONFIG.ENABLE_TLS {
		err := startTLS()
		log.Fatal(err)
	} else {
		err := startNonTLS()
		log.Fatal(err)
	}
}

func startTLS() error {
	err := http.ListenAndServeTLS(
		fmt.Sprintf(":%s", CONFIG.SERVER_PORT),
		CONFIG.CERT_FILE,
		CONFIG.PRIVATE_KEY_FILE,
		nil,
	)
	return err
}

func startNonTLS() error {
	err := http.ListenAndServe(
		fmt.Sprintf(":%s", CONFIG.SERVER_PORT),
		nil,
	)
	return err
}
