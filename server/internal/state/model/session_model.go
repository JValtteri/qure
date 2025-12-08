package model

import (
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
)


var SESSION_KEY_LENGTH = 20

type Session struct {
	Key			crypt.Key		// Session cookie
	ExpiresDt	utils.Epoch
	Fingerprint	crypt.Hash
}
