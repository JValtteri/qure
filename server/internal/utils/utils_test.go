package utils

import (
    "testing"
    "bytes"
)

/*
func UnloadBadJSON() int {
    log.SetOutput(io.Discard)
    var input []byte = badJson
    var expect Event
    var got Event
    loadJSON(input, &got)
    if got != expect {
        //t.Errorf("Expected: %v, Got: %v\n", expect, got)
        return 0
    }
    return 1
}
*/

func TestLoadUnloadJSON(t *testing.T) {
    var input []byte = exampleJson
    var obj ExampleStruct
    LoadJSON(input, &obj)
    var got []byte = []byte(UnloadJSON(obj))
    // Test a section only
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
