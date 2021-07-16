package transport

import "encoding/json"

// NotifyEvent -
type NotifyEvent struct {
	// Send - which node send event
	Send string `json:"send"`

	// Receives -
	Receives []string `json:"receives"`

	// Session -
	Session string `json:"session"`

	// Message -
	Message Message `json:"message"`
}

// Message -
type Message struct {
	// Header -
	Header map[string]interface{} `json:"heade"`

	// Body -
	Body []byte `json:"body"`
}

// EncodedBody -
func (n *NotifyEvent) EncodedBody(target interface{}) error {
	err := json.Unmarshal(n.Message.Body, target)
	if err != nil {
		return err
	}

	return nil
}
