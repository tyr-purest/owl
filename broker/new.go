package broker

import (
	"github.com/nats-io/nats.go"
)

// nc -
var nc *nats.EncodedConn

// ns -
type ns struct {
	*nats.EncodedConn
	subs []*nats.Subscription
}

// Init -
func Init() {
	n, _ := nats.Connect("nats:4222")

	nc, _ = nats.NewEncodedConn(n, nats.JSON_ENCODER)
}

// NewBroker -
func NewBroker() Broker {
	ns := &ns{
		nc,
		make([]*nats.Subscription, 0),
	}

	return ns
}
