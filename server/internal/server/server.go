package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	ware "github.com/JValtteri/qure/server/internal/middleware"
	"github.com/JValtteri/qure/server/internal/utils"
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
	if CONFIG.ENABLE_TLS {
		log.Fatal(http.ListenAndServeTLS(
			fmt.Sprintf( ":%s", CONFIG.SERVER_PORT),
			CONFIG.CERT_FILE,
			CONFIG.PRIVATE_KEY_FILE, nil))
	} else {
		log.Fatal(http.ListenAndServe( fmt.Sprintf(":%s", CONFIG.SERVER_PORT), nil))
	}
}

func defaultRequest(w http.ResponseWriter, request *http.Request) {
	http.ServeFile(w, request, "./static/index.html")
}

func getEvents(w http.ResponseWriter, request *http.Request) {
	isAdmin := false
	obj := ware.GetEvents(isAdmin)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", utils.UnloadJSON(obj))
}

func getEvent(w http.ResponseWriter, request *http.Request) {
	isAdmin := false
	bodyObject, err := loadBody(request, ware.EventRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	event := ware.GetEvent(bodyObject.EventID, isAdmin)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", utils.UnloadJSON(event))
}

func authenticateSession(w http.ResponseWriter, request *http.Request) {
}

func loginUser(w http.ResponseWriter, request *http.Request) {
}

func userReservations(w http.ResponseWriter, request *http.Request) {
}

func makeReservation(w http.ResponseWriter, request *http.Request) {
}

func registerUser(w http.ResponseWriter, request *http.Request) {
	registerRequest, err := loadBody(request, ware.RegisterRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := ware.Register(registerRequest)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", utils.UnloadJSON(resp))
}

func loginWithReservation(w http.ResponseWriter, request *http.Request) {
}


func sanitize(input string) string {
	var result strings.Builder
	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])
		if unicode.IsSpace(r) || unicode.IsLetter(r) || unicode.IsDigit(r) || r=='-' {
			result.WriteRune(r)
			i += size
		} else {
			i++
		}
	}
	return strings.ToLower(result.String())
}

func loadBody [R ware.Request](request *http.Request, obj R) (R, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error closing request body %v\n", err)
		return obj, err
	}
	defer close(request.Body)
	utils.LoadJSON(body, obj)
	return obj, nil
}

func close(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Printf("Error closing request body %v\n", err)
	}
}
