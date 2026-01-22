package server

import (
	"net/http"

	ware "github.com/JValtteri/qure/server/internal/middleware"
)


func getEvents(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.EventRequest{}, ware.GetEvents)
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
/*
func editReservation(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.ReserveRequest{}, ware.EditReservation)
}
*/

func registerUser(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.RegisterRequest{}, ware.Register)
}

// loginWithReservation (ware.EventLogin{}, ware.ReservationLogin) was removed
// Duplicate of login system creates duplicate of everything in backend and fron
// and is ultimately unnecessary.
//
// Now anonymopus reservation creates an account with the
// _email_			as username and
// _reservation.Id_ as password

func changePassword(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.PasswordChangeRequest{}, ware.ChangePassword)
}

func deleteUser(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.RemovalRequest{}, ware.RemoveUser)
}

func createEvent(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.EventManipulationRequest{}, ware.MakeEvent)
}

func editEvent(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.EventManipulationRequest{}, ware.EditEvent)
}

func deleteEvent(w http.ResponseWriter, request *http.Request) {
	genericHandler(w, request, ware.EventManipulationRequest{}, ware.DeleteEvent)
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
