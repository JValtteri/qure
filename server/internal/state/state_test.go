package state

import (
	"testing"

	c "github.com/JValtteri/qure/server/internal/config"

)

func TestInitialize(t *testing.T) {
	c.CONFIG.DB_FILE_NAME = "test.gob"
	Initialize()
}
