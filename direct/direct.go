package direct

import (
	"github.com/guanaco/guancano/core"
)

const Prefix = "direct"

func ComponentCreator(ctx core.Context) (core.Component, error) {
	component := DirectComponent{
		directs: make(map[string]core.Initiator),
	}
	component.SetPrefix(Prefix)
	component.SetContext(ctx)
	return component, nil
}

// Implementation of a DirectComponent. A DirectComponent
// moves messages directly from a named Producer to a named
// Consumer.
type DirectComponent struct {
	core.BaseComponent
	directs map[string]core.Initiator
}

func (d DirectComponent) CreateEndpoint(path string, options map[string]string) core.Endpoint {
	return &directEndpoint{
		name:      path,
		component: d,
	}
}

// Implementation of the endpoint that maps back to the direct links inside of
// the DirectComponent
type directEndpoint struct {
	name      string
	component DirectComponent
}

func (d *directEndpoint) CreateConsumer() (core.Consumer, error) {
	return &directConsumer{
		endpoint: d,
	}, nil
}

type directConsumer struct {
	endpoint *directEndpoint
}

func (d *directConsumer) Name() string {
	return d.endpoint.name
}

func (d *directConsumer) Start(initiator core.Initiator) {
	d.endpoint.component.directs[d.endpoint.name] = initiator
}

func (d *directConsumer) Init() {

}

func (d *directConsumer) Stop() {
	delete(d.endpoint.component.directs, d.endpoint.name)
}

func (d *directConsumer) Close() {

}

func (d *directEndpoint) CreateProducer() (core.Producer, error) {
	return &directProducer{
		endpoint: d,
	}, nil
}

type directProducer struct {
	endpoint *directEndpoint
}

func (d *directProducer) Name() string {
	return d.endpoint.name
}

func (d *directProducer) Init() {

}

func (d *directProducer) Start() {

}

func (d *directProducer) Stop() {

}

func (d *directProducer) Close() {

}

func (d *directProducer) Process(exchange core.Exchange) {
	d.endpoint.component.directs[d.endpoint.name].Exchange(exchange.In())
}
