package main

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/server"
)

func main() {
    fmt.Println("### Server Started ###")
    server.Server()
    fmt.Println("### Server Stopped ###")
}
