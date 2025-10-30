package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
	"github.com/JValtteri/qure/server/internal/utils"
	ware "github.com/JValtteri/qure/server/internal/middleware"
)


func TestSanitizeAllow(t *testing.T) {
    const input string = "asd 123-5"
    var got string = sanitize(input)
    if got != input {
        t.Errorf("Expected: %s, Got: %s\n", input, got)
    }
}

func TestSanitizeBlock(t *testing.T) {
    const input string = "<asd=1,23>"
    const expect string = "asd123"
    var got string = sanitize(input)
    if got != expect {
        t.Errorf("Expected: '%s', Got: '%s'\n", expect, got)
    }
}

func TestGetEvents(t *testing.T) {
	if err := testGetEvents(); err != nil {
		t.Logf("Error in response handler: %v\n", err)
	}
}

func TestGetEvent(t *testing.T) {
	if err := testGetEvent(crypt.ID("nothing")); err == nil {
		t.Errorf("Expected: '%s', Got: '%v'\n", "error", err)
	}
}

func TestRegisterUser(t *testing.T) {
	if err := testRegisterUser("test"); err != nil {
		t.Errorf("Error in response handler: %v\n", err)
	}
}

func testGetEvents() error {
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
	err := eventTester(data)
	if err != nil {
		return fmt.Errorf("getEvents(): %v", err)
	}
	return nil
}

func testGetEvent(eventID crypt.ID) error {
	data := TestData[ware.EventRequest] {
		handler: getEvent,
		expected: TExpected{
            status: http.StatusOK,
			body: `null`,
        },
		request: TRequest[ware.EventRequest] {
			rtype: "POST",
			path: "/api/event",
			body: ware.EventRequest{eventID},
		},
	}
	err := eventTester(data)
	if err != nil {
		return fmt.Errorf("getEvent(): %v", err)
	}
	return nil
}

func testRegisterUser(name string) error {
	data := TestData[ware.RegisterRequest] {
		handler: registerUser,
		expected: TExpected{
            status: http.StatusOK,
			body: `{SessionKey:"asd", Error:""}`,
        },
		request: TRequest[ware.RegisterRequest] {
			rtype: "POST",
			path: "/api/user/register",
			body: ware.RegisterRequest{User: name, Password: "passwd", Ip: state.IP("0.0.0.0")},
		},
	}
	err := eventTester(data)
	if err != nil {
		return fmt.Errorf("registerUser(): %v", err)
	}
	return nil
}

func eventTester[R ware.Request](d TestData[R]) error {
	requestBodyWriter := makeWriter(d.request.body)
	req, err := http.NewRequest(d.request.rtype, d.request.path, requestBodyWriter)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	rr := httptest.NewRecorder()
	d.handler(rr, req)
	if status := rr.Code; status != d.expected.status {
		return fmt.Errorf("handler returned wrong status code: got %v want %v",
			status, d.expected.status)
	}
	if rr.Body.String() != d.expected.body {
		return fmt.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), d.expected.body)
	}
	return nil
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
