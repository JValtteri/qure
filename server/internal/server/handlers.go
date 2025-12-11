package server

import (
	"fmt"
	"net/http"

	c "github.com/JValtteri/qure/server/internal/config"
	ware "github.com/JValtteri/qure/server/internal/middleware"
)


func defaultRequest(w http.ResponseWriter, request *http.Request) {
	http.ServeFile(w, request, fmt.Sprintf("%s/index.html", c.CONFIG.SOURCE_DIR))
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

func logoutUser(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.AuthenticateRequest{}, ware.Logout)
}

func userReservations(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.UserReservationsRequest{}, ware.GetUserReservatoions)
}

func makeReservation(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.ReserveRequest{}, ware.MakeReservation)
}

func registerUser(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.RegisterRequest{}, ware.Register)
}

func loginWithReservation(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.EventLogin{}, ware.ReservationLogin)
}

func changePassword(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.PasswordChangeRequest{}, ware.ChangePassword)
}

func deleteUser(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.RemovalRequest{}, ware.RemoveUser)
}

func createEvent(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.EventCreationRequest{}, ware.MakeEvent)
}

func editEvent(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.EventCreationRequest{}, ware.EditEvent)
}

func genericHandler [R ware.Request, P ware.Response](
	w http.ResponseWriter,	request *http.Request,
	requestType R,			middlewareFunction func(R)P,
) {
	req, err := loadRequestBody(request, ware.UniversalRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := Fingerprint(request)
	sessionKey := getCookie(request, "sessionKey")
	appendFields(&req, fingerprint, sessionKey)
	convertTo(&requestType, req)
	response := middlewareFunction(requestType)
	sendJsonResponse(w, response)
}
