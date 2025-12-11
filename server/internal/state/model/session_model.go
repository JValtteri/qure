package model

import (
	"github.com/JValtteri/qure/server/internal/crypt"
	"github.com/JValtteri/qure/server/internal/utils"
)


type Session struct {
	Key			crypt.Key		// Session cookie
	ExpiresDt	utils.Epoch
	Fingerprint	crypt.Hash
}
