package utils

import(
    "crypto/rand"
    "fmt"
)

// Returns a string containing random chars from [A..Z,a..z,0..9]
func RandomChars(length int) (string, error) {
    ints, err := RandomInts(length, 62)
    for i, v := range ints {
        ints[i] = asciiOffset(v)
    }
    bytes := Itob(ints)
    return string(bytes), err
}

func asciiOffset(v int) int {
        if v < 26 {
            return v+65
        } else if v < 52 {
            return v-26+97
        } else {
            return v-52+48
        }
}

func RandomInts(length int, base int) ([]int, error) {
    ints := make([]int, length)
    bytes, err := RandomBytes(length)
    for i, b := range bytes {
        ints[i] = int(b) % base
    }
    return ints, err
}

func RandomBytes(length int) ([]byte, error) {
    buffer := make([]byte, length)
    _, err := rand.Read(buffer)
    if err != nil {
        return buffer, fmt.Errorf("error generating random bytes: %v", err) // Should not be possible (random byte generation)
    }
    return buffer, nil
}
