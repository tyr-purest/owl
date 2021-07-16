package transport

import (
	"net/http"

	"github.com/tyr-purest/owl/broker"

	uuid "github.com/satori/go.uuid"
	"github.com/tyr-tech-team/hawk/status"
)

// Node -
type Node struct {
	// b - broker
	b broker.Broker

	// ws -  ws connection
	Ws *webSocket

	// ID - node
	ID string

	// OnNotify -
	OnNotify func(e *NotifyEvent)
}

// NewNode -
func NewNode(w http.ResponseWriter, r *http.Request) (*Node, error) {
	ws, err := NewWs(w, r)
	if err != nil {
		return nil, err
	}

	id := uuid.NewV4().String()

	node := &Node{
		Ws: ws,
		b:  broker.NewBroker(),
		ID: id,
	}

	return node, nil
}

// WsReply -
func (n *Node) WsReply(data []byte) {
	n.Ws.out <- data
}

// Notify -
func (n *Node) Notify(topic string, e *NotifyEvent) error {
	if e.Session == "" {
		return status.InvalidParameter.Err()
	}
	err := n.b.Pub(topic, e)
	if err != nil {
		return err
	}

	return nil
}
