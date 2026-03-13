package main

import (
    "fmt"
    "github.com/JValtteri/qure/server/internal/server"
	"github.com/JValtteri/qure/server/internal/utils"
)

func main() {
    printStylizedLogo()
	PrintAttribution()
    server.Server()
    fmt.Println("### Server Stopped ###")
}

func printStylizedLogo() {
	fmt.Printf("%s%s%s%s%s%s%s%s%s%s\n",
		"\n",
		"      ██████               ███████████\n",
		"    ███░░░░███            ░░███░░░░░███\n",
		"   ███    ░░███ █████ ████ ░███    ░███   ██████\n",
		"  ░███     ░███░░███ ░███  ░██████████   ███░░███\n",
		"  ░███   ██░███ ░███ ░███  ░███░░░░░███ ░███████\n",
		"  ░░███ ░░████  ░███ ░███  ░███    ░███ ░███░░░\n",
		"   ░░░██████░██ ░░████████ █████   █████░░██████\n",
		"     ░░░░░░ ░░   ░░░░░░░░ ░░░░░   ░░░░░  ░░░░░░\n\n",
		"    QuRe Reservation System\n",
	)
}

func PrintAttribution() {
	var attrStr = utils.ReadFile("ATTRIBUTION")
	fmt.Printf("%s\n", attrStr)
}
