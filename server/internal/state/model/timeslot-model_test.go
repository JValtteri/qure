package model

import (
	"fmt"
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
)

func TestFilter(t *testing.T) {
	var list = makeIDlist([]string{"a", "b", "a", "b", "a", "b", "3"})
	var expected = makeIDlist([]string{"a", "a", "a", "3"})
	var filtered = filterFrom(list, crypt.ID("b"))
	if !areEqual(filtered, expected) {
		t.Fatalf("Expected: %v, Got: %v\n", expected, filtered)
	}
}

func TestFilterN(t *testing.T) {
	var list = makeIDlist([]string{"a", "b", "a", "b", "a", "b", "3"})
	var expected = makeIDlist([]string{"a", "a", "a", "b", "3"})
	var filtered = filterNfrom(2, list, crypt.ID("b"))
	if !areEqual(filtered, expected) {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expected, filtered)
	}
}

func TestPop(t *testing.T) {
	var list = makeIDlist([]string{"a", "b", "a", "b", "a", "b", "3"})
	var expectedPopped = makeIDlist([]string{"a", "b", "a"})
	var expectedList = makeIDlist([]string{"b", "a", "b", "3"})
	var popped, _ = pop(&list, 3)
	if !areEqual(popped, expectedPopped) {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expectedPopped, popped)
	}
	if !areEqual(list, expectedList) {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", list, expectedList)
	}
}

func TestOverIndexPop(t *testing.T) {
	var list = makeIDlist([]string{"a", "b"})
	var expectedPopped = makeIDlist([]string{"a", "b"})
	var popped, err = pop(&list, 3)
	if !areEqual(popped, expectedPopped) {
		t.Fatalf("Expected: %v, Got: %v (oldQueue)\n", expectedPopped, popped)
	}
	if err == nil {
		t.Fatalf("Over index should warn!\n")
	}
}

func makeIDlist(list []string) []crypt.ID {
	var idList = []crypt.ID{}
	for _, value := range list {
		idList = append(idList, crypt.ID(fmt.Sprintf("%v", value)))
	}
	return idList
}

func areEqual(a []crypt.ID, b []crypt.ID) bool {
	if len(a) != len(b) {
		fmt.Println("length mismatch")
		return false
	}
	for index, value := range a {
		if value != b[index] {
			fmt.Printf("A: %v, B: %v\n", value, b[index])
			return false
		}
	}
	return true
}
