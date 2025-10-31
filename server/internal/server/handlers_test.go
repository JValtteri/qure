package server

import (
	"fmt"
	"bytes"
	"regexp"
	"strings"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
	"github.com/JValtteri/qure/server/internal/utils"
	ware "github.com/JValtteri/qure/server/internal/middleware"
)


func TestGetEvents(t *testing.T) {
	if _, err := testGetEvents(); err != nil {
		t.Logf("Error in response handler: %v\n", err)
	}
}

func TestGetEvent(t *testing.T) {
	if _, err := testGetEvent(crypt.ID("nothing")); err != nil {
		t.Errorf("Expected: '%s', Got: '%v'\n", "error", err)
	}
}

func TestRegisterUser(t *testing.T) {
	if _, err := testRegisterUser("test"); err != nil {
		t.Errorf("Error in response handler:\n %v\n", err)
	}
}

func TestRegisterUserTwice(t *testing.T) {
	testRegisterDuplicateUser("double")
	if _, err := testRegisterDuplicateUser("double"); err != nil {
		t.Errorf("Error in response handler:\n %v\n", err)
	}
}

func TestResumeSession(t *testing.T) {
	if _, err := testResumeSession("test"); err != nil {
		t.Errorf("Error in response handler:\n %v\n", err)
	}
}

func TestLogin(t *testing.T) {
	if _, err := testLoginUser("login"); err != nil {
		t.Errorf("Error in response handler:\n %v\n", err)
	}
}

func TestEventLifesycle(t *testing.T) {
	setupFirstAdminUser("admin", deterministicKeyGenerator)
	client, err := state.NewClient("admin", "test-admin", "adminpasswordexample", false)
	if err != nil {
		t.Fatalf("Error generating test-admin account:\n%v", err)
	}
	sessionKey, err := testLoginAdmin("test-admin")
	if err != nil {
		t.Errorf("Response handler:\n%v\n", err)
	}
	eventID, err := testMakeEvent(sessionKey)
	if err != nil {
		t.Errorf("Response handler:\n%v\n", err)
	}
	if len(eventID) < 9 {
		t.Errorf("Unexpected EventID: %v\n", eventID)
	}
	_, err = testReserve(sessionKey, "test-admin", 1, state.ID(eventID), client.Id)
	if err != nil {
		t.Errorf("Response handler:\n%v\n", err)
	}
}


func testGetEvents() (string, error) {
	data := TestData[ware.UniversalRequest] {
		handler: getEvents,
		expected: TExpected{
            status: http.StatusOK,
			body: `null`,
        },
		request: TRequest[ware.UniversalRequest] {
			rtype: "GET",
			path: "/api/events",
			body: ware.UniversalRequest{},
		},
	}
	key, err := eventTester(data, "SessionKey")
	if err != nil {
		return key, fmt.Errorf("getEvents(): %v", err)
	}
	return key, nil
}

func testGetEvent(eventID crypt.ID) (string, error) {
	data := TestData[ware.EventRequest] {
		handler: getEvent,
		expected: TExpected{
            status: http.StatusOK,
			body: `{"ID":"","Name":"","ShortDescription":"","LongDescription":"","Draft":false,"DtStart":0,"DtEnd":0,"StaffSlots":0,"Staff":0,"Timeslots":null}`,
        },
		request: TRequest[ware.EventRequest] {
			rtype: "POST",
			path: "/api/event",
			body: ware.EventRequest{eventID, false},
		},
	}
	key, err := eventTester(data, "SessionKey")
	if err != nil {
		return key, fmt.Errorf("getEvent(): %v", err)
	}
	return key, nil
}

func testRegisterUser(name string) (string, error) {
	data := TestData[ware.RegisterRequest] {
		handler: registerUser,
		expected: TExpected{
            status: http.StatusOK,
			body: `{"SessionKey":"<key>","Error":""}`,
        },
		request: TRequest[ware.RegisterRequest] {
			rtype: "POST",
			path: "/api/user/register",
			body: ware.RegisterRequest{User: name, Password: "password", Ip: state.IP("0.0.0.0")},
		},
	}
	key, err := eventTester(data, "SessionKey")
	if err != nil {
		return key, fmt.Errorf("registerUser(): %v", err)
	}
	return key, nil
}

func testRegisterDuplicateUser(name string) (string, error) {
	data := TestData[ware.RegisterRequest] {
		handler: registerUser,
		expected: TExpected{
            status: http.StatusOK,
			body: `{"SessionKey":"<key>","Error":"error creating client: error: client email not unique"}`,
        },
		request: TRequest[ware.RegisterRequest] {
			rtype: "POST",
			path: "/api/user/register",
			body: ware.RegisterRequest{User: name, Password: "password", Ip: state.IP("0.0.0.0")},
		},
	}
	key, err := eventTester(data, "SessionKey")
	if err != nil {
		return key, fmt.Errorf("registerUser(): %v", err)
	}
	return key, nil
}

func testResumeSession(sessionKey crypt.Key) (string, error) {
		data := TestData[ware.AuthenticateRequest] {
		handler: authenticateSession,
		expected: TExpected{
            status: http.StatusOK,
			body: `{"Authenticated":false,"IsAdmin":false,"SessionKey":"<key>","Error":""}`,
        },
		request: TRequest[ware.AuthenticateRequest] {
			rtype: "POST",
			path: "/api/session/auth",
			body: ware.AuthenticateRequest{SessionKey: sessionKey, Ip: state.IP("0.0.0.0")},
		},
	}
	key, err := eventTester(data, "SessionKey")
	if err != nil {
		return key, fmt.Errorf("authenticateSession(): %v", err)
	}
	return key, nil
}

func testLoginUser(name string) (string, error) {
	data := TestData[ware.LoginRequest] {
		handler: loginUser,
		expected: TExpected{
            status: http.StatusOK,
			body: `{"Authenticated":false,"IsAdmin":false,"SessionKey":"<key>","Error":""}`,
        },
		request: TRequest[ware.LoginRequest] {
			rtype: "POST",
			path: "/api/user/login",
			body: ware.LoginRequest{User: name, Password: crypt.Key("password"), Ip: state.IP("0.0.0.0")},
		},
	}
	key, err := eventTester(data, "SessionKey")
	if err != nil {
		return key, fmt.Errorf("loginUser(): %v", err)
	}
	return key, nil
}

func testLoginAdmin(name string) (string, error) {
	data := TestData[ware.LoginRequest] {
		handler: loginUser,
		expected: TExpected{
            status: http.StatusOK,
			body: `{"Authenticated":true,"IsAdmin":true,"SessionKey":"<key>","Error":""}`,
        },
		request: TRequest[ware.LoginRequest] {
			rtype: "POST",
			path: "/api/user/login",
			body: ware.LoginRequest{User: name, Password: crypt.Key("adminpasswordexample"), Ip: state.IP("0.0.0.0")},
		},
	}
	key, err := eventTester(data, "SessionKey")
	if err != nil {
		return key, fmt.Errorf("loginUser(): %v", err)
	}
	return key, nil
}

func testMakeEvent(sessionKey string) (string, error) {
	event := state.EventFromJson(state.EventJson)
	data := TestData[ware.EventCreationRequest] {
		handler: createEvent,
		expected: TExpected{
            status: http.StatusOK,
			body: `{"EventID":"<key>"}`,
        },
		request: TRequest[ware.EventCreationRequest] {
			rtype: "POST",
			path: "/api/admin/create",
			body: ware.EventCreationRequest{crypt.Key(sessionKey), state.IP("0.0.0.0"), event},
		},
	}
	key, err := eventTester(data, "EventID")
	if err != nil {
		return key, fmt.Errorf("createEvent(): %v", err)
	}
	return key, nil
}

func testReserve(sessionKey string, name string, size int, eventID state.ID, clientID state.ID) (string, error) {
	data := TestData[ware.ReserveRequest] {
		handler: makeReservation,
		expected: TExpected{
            status: http.StatusOK,
			body: fmt.Sprintf(`{"Id":"<key>","EventID":"%v","ClientID":"%v","Size":1,"Confirmed":1,"Timeslot":1100,"Expiration":4700,"Error":""}`, eventID, clientID),
        },
		request: TRequest[ware.ReserveRequest] {
			rtype: "POST",
			path: "/api/user/reserve",
			body: ware.ReserveRequest{crypt.Key(sessionKey), name, state.IP("0.0.0.0"), size, eventID, state.Epoch(1100)},
		},
	}
	key, err := eventTester(data, "Id")
	if err != nil {
		return key, fmt.Errorf("makeReservation(): %v", err)
	}
	return key, nil
}

func eventTester[R ware.Request](d TestData[R], keyName ...string) (string, error) {
	requestBodyWriter := makeWriter(d.request.body)
	req, err := http.NewRequest(d.request.rtype, d.request.path, requestBodyWriter)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	rr := httptest.NewRecorder()
	d.handler(rr, req)
	if status := rr.Code; status != d.expected.status {
		return "", fmt.Errorf("handler returned wrong status code:\n got:  %v\n want: %v",
			status, d.expected.status)
	}
	sessionKey := extractRandom(rr.Body.String(), keyName[0])
	strippedBody := stripRandom(rr.Body.String(), keyName)
	if strippedBody != d.expected.body {
		return "", fmt.Errorf("handler returned unexpected body:\n got:  %v\n want: %v",
			strippedBody, d.expected.body)
	}
	return sessionKey, nil
}

type TestData [R ware.Request]struct {
	handler func(http.ResponseWriter, *http.Request)
	expected TExpected
	request TRequest[R]
}

type TExpected struct {
	status	int
	body	string
}

type TRequest [R ware.Request]struct {
	rtype string
	path string
	body R
}

func makeWriter [R ware.Request](r R) *bytes.Buffer {
	strJson := utils.UnloadJSON(r)
	if strJson == "{}" {
		a := bytes.NewBufferString("")
		fmt.Printf("%v", a)
		return a
	}
	return bytes.NewBufferString(strJson)
}

func extractRandom(input string, keyName string) string {
	// Replaces the random session key with "<key>"
	regexpSpell := fmt.Sprintf(`("%s"\s*:\s*")[^"]*(")`, keyName)
	re := regexp.MustCompile(regexpSpell)
	key := re.FindString(input)
	key = strings.Replace(key, fmt.Sprintf(`"%s":`, keyName), "", 1)
	key = strings.ReplaceAll(key, `"`, "")
	return key
}

func stripRandom(str string, keyNames []string) string {
	// Replaces the random session key with "<key>"
	for _, key := range(keyNames) {
		regexpSpell := fmt.Sprintf(`("%s"\s*:\s*")[^"]*(")`, key)
		re := regexp.MustCompile(regexpSpell)
		str = re.ReplaceAllString(str, fmt.Sprintf(`"%s":"<key>"`, key))
	}
	return str
}

func deterministicKeyGenerator(keyType *crypt.Key, length int) (crypt.Key, error) {
	return crypt.Key("adminpasswordexample"), nil
}
