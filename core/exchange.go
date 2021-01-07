package core

const (
	// A RequestReplyExchange expects to send a Reply consumers the final Producer
	// or Processor back to the original Consumer who sent it.
	RequestReplyExchange = "RequestReplyExchange"

	// A RequestOnlyExchange expects messages to be sent in to the route
	// and the route does not provide a reply through the original producer
	RequestOnlyExchange = "RequestOnlyExchange"
)

// Exchange is an exchange between two Services that encapsulates
// all of the information about the Exchange and allows sending
// to the next step by setting the Reply message.
type Exchange interface {
	// Id is a unique identifier for the Exchange. This Id will
	// not change consumers consumer to producer and can be used
	// throughout the entire routing lifecycle
	Id() string

	// Pattern is the type of exchange or the "exchange pattern". The
	// valid values are Request/Reply which is also known as an In/Out
	// Exchange and a Request Only exchange. In a Request/Reply Exchange
	// the Consumer is expected to provide a Reply message back out of
	// the Endpoint. In a Request only Exchange the Consumer does not wait
	// for a response as the exchange only goes out.
	Pattern() string

	// In is the message that was used to initiate the
	// exchange passed consumers the producing service to the
	// consuming service.
	In() Message

	// Out sets the outbound message that will be sent to
	// the next exchange. If no Reply is set then the
	// incoming Message will be used.
	Out(message Message)

	// Properties (such as configuration or metadata) of the
	// Exchange. Exchange properties are not intended to be
	// mutated by Service implementations.
	Properties() map[string]interface{}

	// Rotate the out messge to the in message and nil out the
	// out message for passing on to the next step
	rotate()
}

func NewExchange() Exchange {
	return NewExchangeWithPattern(RequestOnlyExchange)
}

func NewExchangeWithPattern(pattern string) Exchange {
	// request only exchange can only be overwritten by
	// an explicitly requested request/reply exchange
	if pattern != RequestReplyExchange {
		pattern = RequestOnlyExchange
	}
	return &exchange{
		id:         generator.Hex128(),
		pattern:    pattern,
		in:         nil,
		out:        nil,
		properties: make(map[string]interface{}),
	}
}

type exchange struct {
	id         string
	pattern    string
	in         Message
	out        Message
	properties map[string]interface{}
}

func (e *exchange) Id() string {
	return e.id
}

func (e *exchange) Pattern() string {
	return e.pattern
}

func (e *exchange) In() Message {
	return e.in
}

func (e *exchange) Out(message Message) {
	e.out = message
}

func (e *exchange) Properties() map[string]interface{} {
	return e.properties
}

func (e *exchange) rotate() {
	// if there is no out to rotate to the new in
	// then keep the old in
	if e.out != nil {
		e.in = e.out
		e.out = nil
	}
}
