package broker

// Pub -
func (n *ns) Pub(topic string, i interface{}) error {
	if err := n.Publish(topic, i); err != nil {
		return err
	}

	return nil
}

// Sub -
func (n *ns) Sub(topic string, cb interface{}) error {
	sub, err := n.Subscribe(topic, cb)
	if err != nil {
		return err
	}

	n.subs = append(n.subs, sub)

	return nil
}

// Close -
func (n *ns) Close() {
	for _, s := range n.subs {
		s.Unsubscribe()
	}
}
