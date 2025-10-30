package server

import (
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
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
    adminName := "admin"
    setupFirstAdminUser(adminName, crypt.CreateHumanReadableKey)
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

func setupFirstAdminUser(adminName string, keyGenerator func(*crypt.Key, int)(crypt.Key, error) ) {
    var keyType crypt.Key
    if state.AdminClientExists() {
        return
    }
    adminKeyLength := 25
    adminPassword, err := keyGenerator(&keyType, adminKeyLength)
    if err != nil {
        log.Fatalf("Creating password for first admin account failed!\nThis may be a hardware issue.")
    }
    state.NewClient("admin", adminName, adminPassword, false)
    log.Printf("An admin user was created\n Username: %v, Password: %v\n Please change the password on first login",
    adminName, adminPassword)
}
