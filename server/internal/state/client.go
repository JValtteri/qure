package state

import (
	"fmt"
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
	c "github.com/JValtteri/qure/server/internal/config"
	"github.com/JValtteri/qure/server/internal/state/model"
)


func GetClientByEmail(email string) (*model.Client, bool) {
	clients.RLock()
	defer clients.RUnlock()
	client, found := clients.ByEmail[email]
	return client, found
}

func GetClientByID(clientID crypt.ID) (*model.Client, bool) {
	clients.RLock()
	defer clients.RUnlock()
	client, found := clients.ByID[clientID]
	return client, found
}

func GetClientBySession(sessionKey crypt.Key) (*model.Client, bool) {
	clients.RLock()
	defer clients.RUnlock()
	client, found := clients.BySession[sessionKey]
	return client, found
}

func NewClient(role string, email string, password crypt.Key, temp bool) (*model.Client, error) {
	var client *model.Client
	var err error
	var expire utils.Epoch = calculateExpiration(temp)
	sessionKey, err := model.CreateUniqueKey(c.CONFIG.SESSION_KEY_LENGTH, clients.BySession)
    if err != nil {
        return client, fmt.Errorf("Error creating key: %v", err)
    }
	client, err = createClient(email, expire, password, role, temp)
    if err != nil {
        return client, fmt.Errorf("error creating client: %v", err)
    }
	registerClient(client, sessionKey)						// Thread safe action is done here
    return client, err
}

func ChangeClientPassword(client *model.Client, password crypt.Key) {
	clients.Lock()
	defer clients.Unlock()
	client.Password = crypt.GenerateHash(password)
}

func RemoveClient(client *model.Client) {
	for session := range(client.Sessions) {
		delete(clients.BySession, session)
	}
	delete(clients.ByEmail, client.Email)
	delete(clients.ByID, client.Id)
}

func AdminClientExists() bool {
	for _, v := range(clients.ByEmail) {
		if v.Role == "admin" {
            return true
        }
    }
    return false
}

func createClient(email string, expire utils.Epoch, password crypt.Key, role string, isTemp bool) (*model.Client, error) {
	var client *model.Client
	var username = email
	if !uniqueEmail(email) {
		return nil, fmt.Errorf("error: client email not unique")
	}
	id, err := model.CreateUniqueID(16, clients.ByID)
	if err != nil {
		return nil, fmt.Errorf("error: Creating a new client\n%v", err) // Should not be possible (random byte generation)
	}
	if isTemp {
		username, err = model.CreateUniqueID(16, clients.ByEmail)
		if err != nil {
			return nil, fmt.Errorf("error: Creating a new ID\n%v", err) // Should not be possible (random byte generation)
		}
		password = crypt.Key(id)
	}
	client = model.CreateClient(id, expire, username, password, role)
	client.IsTemporary = isTemp
	if isTemp {
		client.Email = email
	}
	return client, nil
}

func uniqueEmail(email string) bool {
    return model.Unique(email, clients.ByEmail)
}

func registerClient(client *model.Client, sessionKey crypt.Key) {
	clients.Lock()
	defer clients.Unlock()
	clients.ByID[client.Id] = client;
	clients.BySession[sessionKey] = client;
	clients.ByEmail[client.Email] = client;
}

func calculateExpiration(temp bool) utils.Epoch {
	var expire utils.Epoch = 0
    if temp {
        expire = utils.EpochNow() + c.CONFIG.TEMP_CLIENT_AGE
    } else {
        expire-- // Set expire to maximum
    }
    return expire
}
