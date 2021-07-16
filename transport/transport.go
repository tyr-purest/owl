package transport

import (
	"sync"

	"github.com/tyr-purest/owl/broker"

	"github.com/davecgh/go-spew/spew"
)

var transport *TransPort

// Transport -
type TransPort struct {
	mux sync.Mutex

	// nodes - 管理自身Server的所有ws
	nodes map[string][]*Node

	// b - broker
	b broker.Broker
}

// NewTransport -
func NewTransport() {
	transport = &TransPort{
		nodes: make(map[string][]*Node),
		b:     broker.NewBroker(),
	}

	// Notify -
	transport.b.Sub(Notify, func(e *NotifyEvent) {
		ns, ok := transport.GetSessionNodes(e.Session)
		if !ok {
			return
		}

		for _, n := range ns {
			n.OnNotify(e)
		}
	})

	// ToOther -
	transport.b.Sub(ToOther, func(e *NotifyEvent) {
		ns, ok := transport.GetSessionNodes(e.Session)
		if !ok {
			return
		}

		for _, n := range ns {
			if n.ID != e.Send {
				n.OnNotify(e)
			}
		}
	})

	// ToSomePeople -
	transport.b.Sub(ToSomePeople, func(e *NotifyEvent) {
		ns, ok := transport.GetSessionNodes(e.Session)
		if !ok {
			return
		}

		receivers := make(map[string]struct{})
		for _, v := range e.Receives {
			receivers[v] = struct{}{}
		}

		for _, n := range ns {
			if _, ok := receivers[n.ID]; ok {
				n.OnNotify(e)
			}
		}
	})
}

// Join -
func Join(node *Node, session string) {
	nodes, ok := transport.GetSessionNodes(session)
	if !ok {
		transport.SetSessionNodes(session, []*Node{node})
		return
	}

	nodes = append(nodes, node)
	transport.SetSessionNodes(session, nodes)
}

// Leave -
func Leave(session, nodeID string) {
	nodes, ok := transport.GetSessionNodes(session)
	if !ok {
		return
	}

	for i, v := range nodes {
		if v.ID == nodeID {
			nodes = append(nodes[:i], nodes[i+1:]...)
		}
	}

	transport.SetSessionNodes(session, nodes)

	spew.Dump(len(transport.nodes))
}

// GetSessionNodes -
func (t *TransPort) GetSessionNodes(session string) ([]*Node, bool) {
	t.mux.Lock()
	defer t.mux.Unlock()

	nodes, ok := t.nodes[session]
	return nodes, ok
}

// SetSessionNodes -
func (t *TransPort) SetSessionNodes(session string, nodes []*Node) {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.nodes[session] = nodes
}
