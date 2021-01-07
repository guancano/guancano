package core

import (
	"fmt"
	"strings"
)

type RouteCreator func(builder RouteBuilder)

type RouteBuilder interface {
	From(endpoint Endpoint) RouteConfiguration
	FromS(endpoint string) RouteConfiguration
	FromF(endpoint string, args ...interface{}) RouteConfiguration
}

type routeBuilder struct {
	components          map[string]Component
	routeConfigurations []*routeConfiguration
}

func (r *routeBuilder) From(endpoint Endpoint) RouteConfiguration {
	routeConfiguration := &routeConfiguration{
		components: r.components,
		route: route{
			consumers:  make([]Consumer, 0),
			processors: make([]Processor, 0),
		},
	}
	r.routeConfigurations = append(r.routeConfigurations, routeConfiguration)
	return routeConfiguration.From(endpoint)
}

func (r *routeBuilder) FromS(endpoint string) RouteConfiguration {
	routeConfiguration := &routeConfiguration{
		components: r.components,
		route: route{
			consumers:  make([]Consumer, 0),
			processors: make([]Processor, 0),
		},
	}
	r.routeConfigurations = append(r.routeConfigurations, routeConfiguration)
	return routeConfiguration.FromS(endpoint)
}

func (r *routeBuilder) FromF(endpoint string, args ...interface{}) RouteConfiguration {
	return r.FromS(fmt.Sprintf(endpoint, args...))
}

type RouteConfiguration interface {
	RouteBuilder

	RequestReply() RouteConfiguration
	RequestOnly() RouteConfiguration

	To(endpoint Endpoint) RouteConfiguration
	ToS(endpoint string) RouteConfiguration
	ToF(endpoint string, args ...interface{}) RouteConfiguration

	Process(processor Processor) RouteConfiguration
	ProcessFunction(processorFunc ProcessingFunction) RouteConfiguration

	build() Route
}

type routeConfiguration struct {
	components map[string]Component
	route      route
}

func (r *routeConfiguration) From(endpoint Endpoint) RouteConfiguration {
	producer, err := endpoint.CreateConsumer()
	if err != nil {
		// todo: throw error or log? (waiting on choosing a log framework)
		return r
	}
	r.route.consumers = append(r.route.consumers, producer)
	return r
}

func (r *routeConfiguration) FromS(endpoint string) RouteConfiguration {
	// get prefix
	idx := strings.Index(endpoint, ":")
	if value, found := r.components[endpoint[0:idx]]; found {
		return r.From(value.CreateEndpoint(endpoint, map[string]string{}))
	}
	return r
}

func (r *routeConfiguration) FromF(endpoint string, args ...interface{}) RouteConfiguration {
	return r.FromS(fmt.Sprintf(endpoint, args...))
}

func (r *routeConfiguration) To(endpoint Endpoint) RouteConfiguration {
	consumer, err := endpoint.CreateProducer()
	if err != nil {
		// todo: throw error or log? (waiting on choosing a log framework)
		return r
	}
	r.route.processors = append(r.route.processors, consumer)
	return r
}

func (r *routeConfiguration) ToS(endpoint string) RouteConfiguration {
	// get prefix
	idx := strings.Index(endpoint, ":")
	if value, found := r.components[endpoint[0:idx]]; found {
		return r.To(value.CreateEndpoint(endpoint, map[string]string{}))
	}
	return r
}

func (r *routeConfiguration) ToF(endpoint string, args ...interface{}) RouteConfiguration {
	return r.ToS(fmt.Sprintf(endpoint, args...))
}

func (r *routeConfiguration) Process(processor Processor) RouteConfiguration {
	r.route.processors = append(r.route.processors, processor)
	return r
}

func (r *routeConfiguration) ProcessFunction(processingFunction ProcessingFunction) RouteConfiguration {
	return r.Process(processHolder{
		processingFunction: processingFunction,
	})
}

func (r *routeConfiguration) RequestReply() RouteConfiguration {
	if r.route.pattern == "" {
		r.route.pattern = RequestReplyExchange
	}
	return r
}

func (r *routeConfiguration) RequestOnly() RouteConfiguration {
	if r.route.pattern == "" {
		r.route.pattern = RequestOnlyExchange
	}
	return r
}

func (r *routeConfiguration) build() Route {
	r.route.id = generator.Hex128()
	return &r.route
}

type Route interface {
	ConsumingService
}

type route struct {
	id        string
	pattern   string
	initiator Initiator

	consumers  []Consumer
	processors []Processor
}

func (r *route) Init() {
	for _, f := range r.consumers {
		f.Init()
	}
	for _, s := range r.processors {
		if c, ok := s.(Producer); ok {
			c.Init()
		}
	}
}

func (r *route) Start() {
	r.initiator = &routeInitiator{
		route: r,
	}

	for _, producer := range r.consumers {
		producer.Start(r.initiator)
	}

	for _, s := range r.processors {
		if c, ok := s.(Producer); ok {
			c.Start()
		}
	}
}

func (r *route) Stop() {
	for _, f := range r.consumers {
		f.Stop()
	}
	for _, s := range r.processors {
		if c, ok := s.(Producer); ok {
			c.Stop()
		}
	}
}

func (r *route) Close() {
	for _, f := range r.consumers {
		f.Close()
	}
	for _, s := range r.processors {
		if c, ok := s.(Producer); ok {
			c.Close()
		}
	}
}

type routeInitiator struct {
	route *route
}

func (r *routeInitiator) Exchange(in Message) Exchange {
	// create initial exchange
	exchange := NewExchangeWithPattern(r.route.pattern)
	exchange.Out(in)
	exchange.rotate()

	// for each step handle the in/out at each step, essentially
	// rotating the out message to be the in message for the
	// next step
	for idx := 0; idx < len(r.route.processors); idx++ {
		if r.route.processors[idx] != nil {
			r.route.processors[idx].Process(exchange)
			exchange.rotate()
		}
	}

	// rotate and return the exchange
	exchange.rotate()
	return exchange
}

func (r *routeInitiator) Pattern() string {
	return r.route.pattern
}
