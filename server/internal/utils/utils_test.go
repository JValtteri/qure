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
    err := LoadJSON(input, &obj)
    if err != nil {
        t.Errorf("JSON unmarshal error: %v" , err)
    }
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
