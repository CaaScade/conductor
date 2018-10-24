package manager

import (
	"conductor/pkg/server"
)

func Spawn() error {
	return server.Run()
}
