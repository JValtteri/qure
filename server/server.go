package main

import (
    "fmt"
    "net/http"
    "log"
    "strings"
    "unicode"
    "unicode/utf8"
)

func server() {
    log.Println("Server UP")
    LoadConfig()
    http.HandleFunc("/", defaultRequest)
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
    http.Handle("/js/",  http.StripPrefix("/js/",  http.FileServer(http.Dir("./static/js"))))
    http.Handle("/img/", http.StripPrefix("/img/",  http.FileServer(http.Dir("./static/img"))))
    if CONFIG.ENABLE_TLS {
        log.Fatal(http.ListenAndServeTLS(
            fmt.Sprintf( ":%s", CONFIG.SERVER_PORT),
            CONFIG.CERT_FILE,
            CONFIG.PRIVATE_KEY_FILE, nil))
    } else {
        log.Fatal(http.ListenAndServe( fmt.Sprintf(":%s", CONFIG.SERVER_PORT), nil))
    }
}

func defaultRequest(w http.ResponseWriter, request *http.Request) {
    http.ServeFile(w, request, "./static/index.html")
}

func setCorrs(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", CONFIG.ORIGIN_URL)
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func sanitize(input string) string {
    var result strings.Builder
    for i := 0; i < len(input); {
        r, size := utf8.DecodeRuneInString(input[i:])
        if unicode.IsSpace(r) || unicode.IsLetter(r) || unicode.IsDigit(r) {
            result.WriteRune(r)
            i += size
        } else {
            i++
        }
    }
    return strings.ToLower(result.String())
}
