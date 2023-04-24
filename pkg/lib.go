package pkg

import (
	"github.com/wagoodman/go-partybus"

	"github.com/nextlinux/go-logger"
	"github.com/nextlinux/vulnersdb/internal/bus"
	"github.com/nextlinux/vulnersdb/internal/log"
)

func SetLogger(l logger.Logger) {
	log.Set(l)
}

func SetBus(b *partybus.Bus) {
	bus.SetPublisher(b)
}
