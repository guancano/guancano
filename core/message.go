package core

import "fmt"

// A Message is sent and/or consumed by each service through the
// use of an Exchange. Each Exchange has an incoming Message and
// expects an outgoing Message as well.
type Message interface {
	// Each Message has a body that is the contents of the message
	Body() interface{}

	// And the ability to update the contents of the message
	Update(body interface{})

	// Each Message has Headers which can be used to send Message
	// metadata. Different Producer and Consumer implementations will
	// set these Headers in their own way.
	Headers() *map[string]interface{}
}

type coreMessage struct {
	headers *map[string]interface{}
	body    interface{}
}

func newCoreMessage(body interface{}) coreMessage {
	headers := make(map[string]interface{})
	return coreMessage{
		headers: &headers,
		body:    body,
	}
}

func (c coreMessage) Headers() *map[string]interface{} {
	return c.headers
}

func (c coreMessage) Body() interface{} {
	return c.body
}

func (c coreMessage) Update(body interface{}) {
	c.body = body
}

func NewTextMessage(text string) TextMessage {
	return TextMessage{
		coreMessage: newCoreMessage(text),
	}
}

type TextMessage struct {
	coreMessage
}

func (t TextMessage) Text() string {
	return fmt.Sprintf("%v", t.body)
}
