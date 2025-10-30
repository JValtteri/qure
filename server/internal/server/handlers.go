package server

import (
	"net/http"

	ware "github.com/JValtteri/qure/server/internal/middleware"
)


func defaultRequest(w http.ResponseWriter, request *http.Request) {
	http.ServeFile(w, request, "./static/index.html")
}

func getEvents(w http.ResponseWriter, request *http.Request) {
	isAdmin := false
	events := ware.GetEvents(isAdmin)
	sendJsonResponse(w, events)
}

func getEvent(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.EventRequest{}, ware.GetEvent)
}

func authenticateSession(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.AuthenticateRequest{}, ware.AuthenticateSession)
}

func loginUser(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.LoginRequest{}, ware.Login)
}

func userReservations(w http.ResponseWriter, request *http.Request) {
}

func makeReservation(w http.ResponseWriter, request *http.Request) {
}

func registerUser(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.RegisterRequest{}, ware.Register)
}

func loginWithReservation(w http.ResponseWriter, request *http.Request) {
}

func createEvent(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.EventCreationRequest{}, ware.MakeEvent)
}

func genericHandler [R ware.Request, P ware.Response](w http.ResponseWriter, request *http.Request, requestType R, middlewareFunction func(R)P) {
	req, err := loadRequestBody(request, requestType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := middlewareFunction(req)
	sendJsonResponse(w, response)
}
