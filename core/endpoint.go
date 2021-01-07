package core

import (
	"fmt"
	"strings"
)

// An Endpoint can create a producer or consumer. If there
// is an Error creating the Consumer or Producer an error should
// be returned. Not every endpoint will create both a Producer
// and a Consumer.
type Endpoint interface {
	// Create the Consumer which will consume messages and put them
	// onto an Exchange for delivery. If the Endpoint cannot create
	// a Consumer then it should return the error NotAProducerEndpoint.
	CreateConsumer() (Consumer, error)

	// Create the Producer which will take messages consumers the Exchange
	// and send them to some other place. If the Endpoint cannot create
	// a Producer then it should return the error NotAConsumerEndpoint.
	CreateProducer() (Producer, error)
}

// NotAConsumerEndpoint is an error that should be returned when the
// Endpoint cannot create a Producer.
type NotAConsumerEndpoint struct {
}

func (n NotAConsumerEndpoint) Error() string {
	return fmt.Sprint("This Endpoint cannot create Consumers")
}

// NotAProducerEndpoint is an error that should be returned when the
// Endpoint cannot create a Consumer
type NotAProducerEndpoint struct {
}

func (n NotAProducerEndpoint) Error() string {
	return fmt.Sprint("This Endpoint cannot create Producers")
}

// Parse returns the parsed information consumers an endpiont string. This allows
// implementors of endpoints to test the parsing behavior that will be
// followed by the FromS/FromF and ToS/ToF methods of the RouteBuilder and
// RouteConfiguration
func Parse(endpoint string) (string, string, string, map[string]string) {
	idx := strings.Index(endpoint, ":")
	options := make(map[string]string)
	if idx < 0 {
		return "", "", "", options
	}
	optx := strings.Index(endpoint, "?")
	if optx < 0 || optx < idx {
		return endpoint[0:idx], endpoint[idx+1:], "", options
	}
	// todo: parse the options string
	return endpoint[0:idx], endpoint[idx+1 : optx], endpoint[optx+1:], options
}
