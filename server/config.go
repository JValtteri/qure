package main

import (
    "log"
    "os"
    "encoding/json"
)

type Config struct {
    ORIGIN_URL       string
    SERVER_PORT      string
    ENABLE_TLS       bool
    CERT_FILE        string
    PRIVATE_KEY_FILE string
}

const CONFIG_FILE string = "config.json"
var CONFIG Config

func LoadConfig() {
    raw_config := readConfig(CONFIG_FILE)
    unmarshal(raw_config, &CONFIG)
    log.Printf("Server url/port: %s:%s\n", CONFIG.ORIGIN_URL, CONFIG.SERVER_PORT)
    if CONFIG.ENABLE_TLS {
        log.Println("TLS is Enabled")
    } else {
        log.Println("HTTP-Only mode")
    }
}

func readConfig(fileName string) []byte {
    raw_config, err := os.ReadFile(fileName)
    if err != nil {
        log.Fatal(err)
    }
    return raw_config
}

func unmarshal(data []byte, config any) {
    err := json.Unmarshal(data, &config)
    if err != nil {
        log.Fatal(err)
    }
}
