package utils

import (
	"bytes"
	"testing"

	"github.com/JValtteri/qure/server/internal/testjson"
)

type ExampleStruct struct {     // Needed to avoid an import cycle
	ID					string
	Name				string;
	ShortDescription	string;
	LongDescription		string;
	DtStart				Epoch;
	DtEnd				Epoch;
	StaffSlots			int;
	Staff				int;
}

func TestUnloadBadJSON(t *testing.T) {
	var input []byte = testjson.BadJson
	var expect ExampleStruct
	var got ExampleStruct
	err := LoadJSON(input, &got)
	if err == nil {
		t.Errorf("Expected error:\n")
	}
	if got != expect {
		t.Errorf("Expected: %v, Got: %v\n", expect, got)
	}
}

func TestLoadUnloadJSON(t *testing.T) {
	var input []byte = testjson.EventJson
    var obj ExampleStruct
    err := LoadJSON(input, &obj)
    if err != nil {
        t.Errorf("JSON unmarshal error: %v" , err)
    }
    var got []byte = []byte(UnloadJSON(obj))
    if bytes.Contains(got, []byte(`"shortDescription": "Lorem ipsum dolor sit amet, meis illud at his"`)) {
        t.Errorf("Expected: shortDescription to contain lorem ipsum. Got: %v\n", string(got))
    }
}

func TestReadFile(t *testing.T) {
    const input string = "utils.go"
    var expect []byte = []byte("package")
    var got[]byte = ReadFile(input)
    if !bytes.Equal(bytes.Fields(got)[0], expect) {
        t.Errorf("Expected: '%v', Got: '%v'\n", string(expect), string(bytes.Fields(got)[0]))
    }
}

func TestEpoch(t *testing.T) {
    expected := Epoch(1761221589)
    got := EpochNow()
    t.Logf("Epoch now: %v", expected)
    if got < expected {
        t.Errorf("Expected: '%v' < '%v'\n", expected, got)
    }
    if got > expected*10 {
        t.Errorf("Expected: '%v' > '%v'\n", expected*10, got)
    }
}

func TestItoB(t *testing.T) {
    expected := "ABC"
    input := []int{65, 66, 67}
    got := ItoB(input)
    if string(got) < expected {
        t.Errorf("Expected: '%v', Got: '%v'\n", expected, string(got))
    }
}
