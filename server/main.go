package main

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/server"
	"github.com/JValtteri/qure/server/internal/utils"
)

func main() {
    fmt.Println("QuRe Reservation System")
	PrintAttribution()
    server.Server()
    fmt.Println("### Server Stopped ###")
}

func PrintAttribution() {
	var attrStr = utils.ReadFile("ATTRIBUTION")
	fmt.Printf("%s\n", attrStr)
}
