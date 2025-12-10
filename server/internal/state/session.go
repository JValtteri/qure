package state

import (
	"fmt"
	"github.com/JValtteri/qure/server/internal/utils"
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/state/model"
)


func ResumeSession(sessionKey crypt.Key, resumeFingerprint string) (*model.Client, error) {
	client, found := getClient(clients.BySession, sessionKey)
    if !found {
        return client, fmt.Errorf("no session matching key found: %v", sessionKey)
    }
    if !fingerprintMatch(resumeFingerprint, client, sessionKey) {
        RemoveSession(sessionKey)
        return client, fmt.Errorf("fingerprint doesn't match stored fingerprint")
    }
	err := cullExpired(&client.Sessions)
	return client, err
}

func AddSession(client *model.Client, role string, email string, temp bool, fingerprint crypt.Hash) (crypt.Key, error) {
	return client.AddSession(role, email, temp, fingerprint, &clients)
}

func fingerprintMatch(resumeFingerprint string, client *model.Client, sessionKey crypt.Key) bool {
	storedFingerprint := client.Sessions[sessionKey].Fingerprint
    return crypt.CompareToHash(resumeFingerprint, storedFingerprint)
}

func getClient(structure map[crypt.Key]*model.Client, key crypt.Key) (*model.Client, bool) {
	clients.RLock()
	defer clients.RUnlock()
	client, found := structure[key]
	return client, found
}

func cullExpired(sessions *map[crypt.Key]model.Session) error {
    var err error
    for key, session := range *sessions {
        now := utils.EpochNow()
		if now < session.ExpiresDt {
            continue
        }
        err = RemoveSession(key)
    }
    return err
}

func RemoveSession(sessionKey crypt.Key) error {
	clients.Lock()
	defer clients.Unlock()
	client, found := clients.BySession[sessionKey]
	if !found {
		return fmt.Errorf("session remove error: session not found")
	}
	// We trust that client.sessions[sessionKey] matches clients.BySession
	delete(client.Sessions, sessionKey)		// Remove from client's sessions
	delete(clients.BySession, sessionKey)	// Remove from globas sessions
	return nil
}