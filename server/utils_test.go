package main

import (
    "testing"
    "bytes"
)

var badJson []byte = []byte(`
{
    "spam": "eggs
}`)

/*
func TestUnloadBadJSON(t *testing.T) {
    var input []byte = badJson
    const expect string = ""
    var got Event
    loadJSON(input, &got)
    if got != expect {
        t.Errorf("Expected: %v, Got: %v\n", expect, got)
    }
}
*/

func TestLoadUnloadJSON(t *testing.T) {
    var input []byte = eventJSON
    //var expect string = string(eventJSON)
    var obj Event
    loadJSON(input, &obj)
    var got []byte = []byte(unloadJSON(obj))
    // Test a section only
    if bytes.Contains(got, []byte(`"shortDescription": "Lorem ipsum dolor sit amet, meis illud at his"`)) {
        t.Errorf("Expected: shortDescription to contain lorem ipsum. Got: %v\n", string(got))
        //t.Errorf("Expected: %v, Got: %v\n", string(expect), string(got))
    }
}

func TestReadFile(t *testing.T) {
    const input string = "utils.go"
    var expect []byte = []byte("main")
    var got[]byte = readFile(input)
    if !bytes.Equal(bytes.Fields(got)[1], expect) {
        t.Errorf("Expected: '%v', Got: '%v'\n", string(expect), string(bytes.Fields(got)[0]))
    }
}
