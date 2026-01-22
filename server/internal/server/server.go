package server

import (
	"fmt"
	"log"
	"sync"
	"syscall"
	"context"
	"net/http"
	"os/signal"

	"github.com/JValtteri/qure/server/internal/state"
	"github.com/JValtteri/qure/server/internal/crypt"
	c "github.com/JValtteri/qure/server/internal/config"
)


var wg sync.WaitGroup

func Server() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()
	c.LoadConfig(c.CONFIG_FILE)
	setupFirstAdminUser("admin", crypt.CreateHumanReadableKey)
	setupHandlers()
	state.InitWaitGroup(&wg)	// Adds state.presistance_api to WaitGroup
	go start()
	<-ctx.Done()				// Wait for Ctrl+C or other stop signal
	state.Save(c.CONFIG.DB_FILE_NAME)
	wg.Wait()					// Ensures registered processes are complete before exiting
}

func setupHandlers() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(fmt.Sprintf("%s/css", c.CONFIG.SOURCE_DIR)))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(fmt.Sprintf("%s/js", c.CONFIG.SOURCE_DIR)))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(fmt.Sprintf("%s/img", c.CONFIG.SOURCE_DIR)))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(fmt.Sprintf("%s/assets", c.CONFIG.SOURCE_DIR)))))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(fmt.Sprintf("%s", c.CONFIG.SOURCE_DIR)))))
	http.HandleFunc("POST /api/events", getEvents)
	http.HandleFunc("POST /api/event", getEvent)
	http.HandleFunc("POST /api/session/auth", authenticateSession)
	http.HandleFunc("POST /api/user/login", loginUser)
	http.HandleFunc("POST /api/user/logout", logoutUser)
	http.HandleFunc("POST /api/user/list", userReservations)
	http.HandleFunc("POST /api/user/reserve", makeReservation)
	http.HandleFunc("POST /api/user/register", registerUser)
	http.HandleFunc("POST /api/user/change", changePassword)
	http.HandleFunc("POST /api/user/delete", deleteUser)
	http.HandleFunc("POST /api/admin/create", createEvent)
	http.HandleFunc("PUT /api/admin/edit", editEvent)
	http.HandleFunc("POST /api/admin/remove", deleteEvent)
}

func start() {
	log.Println("Server UP")
	if c.CONFIG.ENABLE_TLS {
		err := startTLS()
		log.Fatal(err)
	} else {
		err := startNonTLS()
		log.Fatal(err)
	}
}

func startTLS() error {
	err := http.ListenAndServeTLS(
		fmt.Sprintf(":%s", c.CONFIG.SERVER_PORT),
		c.CONFIG.CERT_FILE,
		c.CONFIG.PRIVATE_KEY_FILE,
		nil,
	)
	return err
}

func startNonTLS() error {
	err := http.ListenAndServe(
		fmt.Sprintf(":%s", c.CONFIG.SERVER_PORT),
		nil,
	)
	return err
}
