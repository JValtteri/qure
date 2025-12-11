package server

import (
	"log"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state"
	c "github.com/JValtteri/qure/server/internal/config"
)

func setupFirstAdminUser(adminName string, keyGenerator func(*crypt.Key, int)(crypt.Key, error) ) {
	var keyType crypt.Key
	if state.AdminClientExists() {
		return
	}
	adminPassword, err := keyGenerator(&keyType, c.CONFIG.FIRST_PASSWORD_LENGTH)
	if err != nil {
		log.Fatalf("Creating password for first admin account failed!\nThis may be a hardware issue.")
	}
	state.NewClient("admin", adminName, adminPassword, false)
	log.Printf("An admin user was created\n Username: %v, Password: %v\n Please change the password on first login",
	adminName, adminPassword)
}
