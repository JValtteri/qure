package model

import (
	"testing"

	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
)


var password	= crypt.Hash("test")
var email 		= "test@mail"
var role		= "test"

func TestClientModel(t *testing.T) {
	client := getTestClient()
	got := client.GetPasswordHash()
	if got != password {
		t.Errorf("Expected: %v, Got: %v\n", password, got)
	}
	gots := client.GetEmail()
	if gots != email {
		t.Errorf("Expected: %v, Got: %v\n", email, gots)
	}
	gots = client.GetRole()
	if gots != role {
		t.Errorf("Expected: %v, Got: %v\n", role, gots)
	}
	gotb := client.IsAdmin()
	if gotb != false {
		t.Errorf("Expected: %v, Got: %v\n", false, gotb)
	}
}


func TestClientsReservations(t *testing.T) {
	clients := getTestClients()
	client := CreateClient(crypt.ID("testID"), utils.EpochNow(), email, crypt.Key(password), role)
	if got := client.GetEmail() ; got != email {
		t.Errorf("Expected: %v, Got: %v\n", email, got)
	}
	clients.ByEmail[email] = client
	clients.ByID["testID"] = client
	res		:= getTestReservation()
	err := clients.AddReservation(client.Id, &res)
	if err != nil {
		t.Fatalf("Adding reservation failed: %s\n", err)
	}
	ress := client.GetReservations()
	if len(ress) != 1 {
		t.Errorf("Expected: %v, Got: %v\n", 1, len(ress))
	}
	err = clients.AddReservation("wrong id", &res)
	if err == nil {
		t.Fatalf("Adding reservation with wrong id should have failed: \n")
	}
}

func TestClientsSessions(t *testing.T) {
	clients := getTestClients()
	client := CreateClient(crypt.ID("testID"), utils.EpochNow(), email, crypt.Key(password), role)
	key, err := client.AddSession(role, email, false, crypt.Hash("123"), &clients)
	if err != nil {
		t.Fatalf("Adding session failed: %s\n", err)
	}
	client2, ok := clients.GetClientBySession(key)
	if !ok {
		t.Fatalf("Client not found \n")
	}
	if client2.Id != client.Id {
		t.Fatalf("Client doesn't match \n")
	}
}


func getTestClient() Client {
	return Client {
		Id:				crypt.ID("1"),
		Password:		password,
		Email:			email,
		Role:			role,
		Sessions:		make(map[crypt.Key]Session),
	}
}

func getTestClients() Clients {
	return  Clients{
		ByID:		make(map[crypt.ID]*Client),
		BySession:	make(map[crypt.Key]*Client),
		ByEmail:	make(map[string]*Client),
	}
}

