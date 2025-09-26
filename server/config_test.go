package main

import (
    "testing"
)

func TestTLSConfig(t *testing.T) {
    expect_url  := ""
    expect_port := ""
    expect_tls  := true
    expect_cer  := "cert.pem"
    expect_pem  := "privkey.pem"
    LoadConfig("test_tls_config.json")
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
    expect_url  := "localhost"
    expect_port := "3000"
    expect_tls  := false
    LoadConfig("config.json.example")
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
