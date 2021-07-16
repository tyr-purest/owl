package broker

// Broker -
type Broker interface {
	// Pub -
	Pub(topic string, i interface{}) error

	// Sub -
	Sub(topic string, cb interface{}) error

	// Close -
	Close()
}
