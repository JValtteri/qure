package server

import (
    "log"
    "github.com/JValtteri/qure/server/internal/utils"
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

func LoadConfig(configName string) {
    raw_config := readConfig(configName)
    unmarshal(raw_config, &CONFIG)
    log.Printf("Server url/port: %s:%s\n", CONFIG.ORIGIN_URL, CONFIG.SERVER_PORT)
    if CONFIG.ENABLE_TLS {
        log.Println("TLS is Enabled")
    } else {
        log.Println("HTTP-Only mode")
    }
}

func readConfig(fileName string) []byte {
    raw_config := utils.ReadFile(fileName)
    return raw_config
}

func unmarshal(data []byte, config any) {
    err := utils.LoadJSON(data, config)
    if err != nil {
        log.Fatalf("JSON unmarshal error: %v" , err)
    }
}
