package server

import (
	"fmt"
	"log"
	"net/http"
)


func Server() {
	log.Println("Server UP")
	LoadConfig(CONFIG_FILE)
	http.HandleFunc("/", defaultRequest)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	http.Handle("/js/",  http.StripPrefix("/js/",  http.FileServer(http.Dir("./static/js"))))
	http.Handle("/img/", http.StripPrefix("/img/",  http.FileServer(http.Dir("./static/img"))))
	http.HandleFunc("GET /api/events", 		getEvents)
	http.HandleFunc("POST /api/event", 			getEvent)
	http.HandleFunc("POST /api/session/auth", 	authenticateSession)
	http.HandleFunc("POST /api/user/login", 	loginUser)
	http.HandleFunc("POST /api/user/list", 		userReservations)
	http.HandleFunc("POST /api/user/reserve", 	makeReservation)
	http.HandleFunc("POST /api/user/register",	registerUser)
	http.HandleFunc("POST /api/res/login", 		loginWithReservation)
	http.HandleFunc("POST /api/admin/create", 	createEvent)
	if CONFIG.ENABLE_TLS {
		log.Fatal(http.ListenAndServeTLS(
			fmt.Sprintf( ":%s", CONFIG.SERVER_PORT),
			CONFIG.CERT_FILE,
			CONFIG.PRIVATE_KEY_FILE, nil))
	} else {
		log.Fatal(http.ListenAndServe( fmt.Sprintf(":%s", CONFIG.SERVER_PORT), nil))
	}
}
