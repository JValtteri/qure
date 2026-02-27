package server

import (
	"fmt"
	"log"
	"time"
	"sync"
	"syscall"
	"context"
	"net/http"
	"os/signal"

	"github.com/JValtteri/qure/server/internal/state"
	"github.com/JValtteri/qure/server/internal/crypt"
	c "github.com/JValtteri/qure/server/internal/config"
	R "github.com/JValtteri/qure/server/internal/server/ratelimiter"
)


var wg sync.WaitGroup

func Server() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()
	c.LoadConfig(c.CONFIG_FILE)
	setupFirstAdminUser("admin", crypt.CreateHumanReadableKey)
	mux := http.NewServeMux()
	setupHandlers(mux)
	state.InitWaitGroup(&wg)	// Adds state.presistance_api to WaitGroup
	var srv = newServer(mux, c.CONFIG.SERVER_PORT)
	go start(srv)
	<-ctx.Done()				// Wait for Ctrl+C or other stop signal
	state.Save(c.CONFIG.DB_FILE_NAME)
	wg.Wait()					// Ensures registered processes are complete before exiting
}

func setupHandlers(mux *http.ServeMux) {
	// Set IP based rate limit
	baceRule := R.NewIPLimiterRule(3, 60, 1)
	fastRule := R.NewIPLimiterRule(3, 70, 1)
	fileHandler(mux, "/css/",		"/css")
	fileHandler(mux, "/js/",		"/js")
	fileHandler(mux, "/img/",		"/img")
	fileHandler(mux, "/assets/",	"/assets")
	fileHandler(mux, "/",			"")
	handlerFunc(mux, "POST /api/events",				getEvents,				fastRule)
	handlerFunc(mux, "POST /api/event",					getEvent,				fastRule)
	handlerFunc(mux, "POST /api/session/auth",			authenticateSession,	baceRule)
	handlerFunc(mux, "POST /api/user/login",			loginUser,				baceRule)
	handlerFunc(mux, "POST /api/user/logout",			logoutUser,				baceRule)
	handlerFunc(mux, "POST /api/user/list",				userReservations,		baceRule)
	handlerFunc(mux, "POST /api/user/reserve",			makeReservation,		baceRule)
	handlerFunc(mux, "POST /api/user/amend",			editReservation,		baceRule)
	handlerFunc(mux, "POST /api/user/cancel",			cancelReservation,		baceRule)
	handlerFunc(mux, "POST /api/user/register",			registerUser,			baceRule)
	handlerFunc(mux, "POST /api/user/change",			changePassword,			baceRule)
	handlerFunc(mux, "POST /api/user/delete",			deleteUser,				baceRule)
	handlerFunc(mux, "POST /api/admin/create",			createEvent,			baceRule)
	handlerFunc(mux, "PUT /api/admin/edit",				editEvent,				baceRule)
	handlerFunc(mux, "POST /api/admin/remove",			deleteEvent,			baceRule)
	handlerFunc(mux, "POST /api/admin/reservations",	getEventReservations,	baceRule)
}

func handlerFunc(
	mux *http.ServeMux,
	pattern string,
	handler func(w http.ResponseWriter, request *http.Request),
	limitRule *R.IPLimiter,
) {
	mux.HandleFunc(
		pattern,
		R.RateLimiter(limitRule,
			handler,
		),
	)
}

func fileHandler(mux *http.ServeMux, route string, path string) {
	mux.Handle(
		route,
		http.StripPrefix(
			route, http.FileServer(
				http.Dir(
					fmt.Sprintf(
						"%s%s",
						c.CONFIG.SOURCE_DIR,
						path,
					),
				),
			),
		),
	)
}

func start(srv *http.Server) {
	log.Println("Server UP")
	if c.CONFIG.ENABLE_TLS {
		err := startTLS(srv)
		log.Fatal(err)
	} else {
		err := startNonTLS(srv)
		log.Fatal(err)
	}
}

func startTLS(srv *http.Server) error {
	err := srv.ListenAndServeTLS(
		c.CONFIG.CERT_FILE,
		c.CONFIG.PRIVATE_KEY_FILE,
	)
	return err
}

func startNonTLS(srv *http.Server) error {
	err := srv.ListenAndServe()
	return err
}

func newServer(mux *http.ServeMux, port string) *http.Server {
	srv := &http.Server{
		Addr:			fmt.Sprintf(":%v", port),
		Handler:		mux,
		ReadTimeout:	1 * time.Second,
		WriteTimeout:	3 * time.Second,
		IdleTimeout:	5 * time.Second,
		ErrorLog: 		nil,				// nil = uses standard log package
	}
	return srv
}
