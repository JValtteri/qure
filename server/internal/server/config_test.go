package server

import (
    "testing"
    "log"
    "io"
    "os"
)

func TestNull1(t *testing.T) {
    log.SetOutput(io.Discard)
    log.SetOutput(os.Stdout)
    if false {
        t.Errorf("Error\n")
    }
}

func TestTLSConfig(t *testing.T) {
    //log.SetOutput(io.Discard)
    expect_url  := ""
    expect_port := ""
    expect_tls  := true
    expect_cer  := "cert.pem"
    expect_pem  := "privkey.pem"
    LoadConfig("test_tls_config.json")
    //log.SetOutput(os.Stdout)
    if CONFIG.ORIGIN_URL != expect_url {
        t.Errorf("Expected: %v, Got: %v\n", expect_url, CONFIG.ORIGIN_URL)
    }
    if CONFIG.SERVER_PORT != expect_port {
        t.Errorf("Expected: %v, Got: %v\n", expect_port, CONFIG.SERVER_PORT)
    }
    if CONFIG.ENABLE_TLS != expect_tls {
        t.Errorf("Expected: %v, Got: %v\n", expect_tls, CONFIG.ENABLE_TLS)
    }
    if CONFIG.CERT_FILE != expect_cer {
        t.Errorf("Expected: %v, Got: %v\n", expect_tls, CONFIG.ENABLE_TLS)
    }
    if CONFIG.PRIVATE_KEY_FILE != expect_pem {
        t.Errorf("Expected: %v, Got: %v\n", expect_tls, CONFIG.ENABLE_TLS)
    }
}

func TestLoadConfig(t *testing.T) {
    log.SetOutput(io.Discard)
    expect_url  := "localhost"
    expect_port := "3000"
    expect_tls  := false
    LoadConfig("../../config.json.example")
    log.SetOutput(os.Stdout)
    if CONFIG.ORIGIN_URL != expect_url {
        t.Errorf("Expected: %v, Got: %v\n", expect_url, CONFIG.ORIGIN_URL)
    }
    if CONFIG.SERVER_PORT != expect_port {
        t.Errorf("Expected: %v, Got: %v\n", expect_port, CONFIG.SERVER_PORT)
    }
    if CONFIG.ENABLE_TLS != expect_tls {
        t.Errorf("Expected: %v, Got: %v\n", expect_tls, CONFIG.ENABLE_TLS)
    }
}
