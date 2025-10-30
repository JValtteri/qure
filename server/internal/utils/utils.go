package utils

import (
    "log"
    "os"
    "time"
    "encoding/json"
)

type Epoch uint

func EpochNow() Epoch {
    return Epoch(uint(time.Now().Unix()))
}

func LoadJSON(data []byte, obj any) error {
    return json.Unmarshal(data, &obj)
}

func UnloadJSON(object any) string {
    body, err := json.Marshal(object)
    if err != nil {
        log.Println("JSON marshalling error:" , err)
    }
    return string(body)
}

func ReadFile(fileName string) []byte {
    raw_file, err := os.ReadFile(fileName)
    if err != nil {
        log.Fatal("File error:" , err)
    }
    return raw_file
}

func ItoB(ints []int) []byte {
    length := len(ints)
    bytes := make([]byte, length)
    for i, v := range ints {
        bytes[i] = byte(v)
    }
    return bytes
}
