package main

import (
    "log"
    "os"
    "encoding/json"
)

func loadJSON(data []byte, obj any) {
    err := json.Unmarshal(data, &obj)
    if err != nil {
        log.Fatal("JSON unmarshal error:" , err)
    }
}

func unloadJSON(object any) string {
    body, err := json.Marshal(object)
    if err != nil {
        log.Println("JSON response marshalling error:" , err)
    }
    return string(body)
}

func readFile(fileName string) []byte {
    raw_file, err := os.ReadFile(fileName)
    if err != nil {
        log.Fatal("File error:" , err)
    }
    return raw_file
}

func itob(ints []int) []byte {
    length := len(ints)
    bytes := make([]byte, length, length)
    for i, v := range ints {
        bytes[i] = byte(v)
    }
    return bytes
}

