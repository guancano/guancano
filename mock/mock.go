package mock

import (
	"github.com/guanaco/guancano/core"
)

const Prefix = "mock"

func ComponentCreator(context core.Context) (core.Component, error) {
	component := MockComponent{
		consumers: make(map[string]*mockConsumer),
		producers: make(map[string]*mockProducer),
	}
	component.SetPrefix(Prefix)
	component.SetContext(context)
	return component, nil
}

type MockComponent struct {
	core.BaseComponent
	consumers map[string]*mockConsumer
	producers map[string]*mockProducer
}

func (m MockComponent) CreateEndpoint(path string, options map[string]string) core.Endpoint {
	return &mockEndpoint{
		name:      path,
		component: &m,
	}
}

func (m *MockComponent) Send(route string, message core.Message) {
	if consumer, found := m.consumers[route]; found {
		for idx := 0; idx < len(consumer.initiators); idx++ {
			exchange := consumer.initiators[idx].Exchange(message)
			if exchange.Pattern() == core.RequestReplyExchange {
				consumer.responses = append(consumer.responses, exchange.In())
			}
		}
	}
}

func (m *MockComponent) ConsumerStats(path string) (int, []core.Message) {
	if consumer, ok := m.consumers[path]; ok {
		if consumer == nil {
			return 0, nil
		}
		return len(consumer.responses), consumer.responses
	}
	return 0, nil
}

func (m *MockComponent) ProducerStats(path string) (int, []core.Message) {
	if producer, ok := m.producers[path]; ok {
		if producer == nil {
			return 0, nil
		}
		return producer.invocations, producer.messages
	}
	return 0, nil
}

type mockEndpoint struct {
	name      string
	component *MockComponent
}

func (m *mockEndpoint) CreateConsumer() (core.Consumer, error) {
	if value, found := m.component.consumers[m.name]; found {
		return value, nil
	}
	consumer := &mockConsumer{
		name:      m.name,
		component: m.component,
	}
	m.component.consumers[m.name] = consumer
	return consumer, nil
}

type mockConsumer struct {
	name       string
	component  *MockComponent
	initiators []core.Initiator
	responses  []core.Message
}

func (m mockConsumer) Init() {

}

func (m mockConsumer) Stop() {

}

func (m mockConsumer) Close() {

}

func (m *mockConsumer) Start(initiator core.Initiator) {
	m.initiators = append(m.initiators, initiator)
}

func (m *mockEndpoint) CreateProducer() (core.Producer, error) {
	if value, found := m.component.producers[m.name]; found {
		return value, nil
	}

	producer := &mockProducer{
		name:        m.name,
		component:   m.component,
		invocations: 0,
		messages:    make([]core.Message, 0),
	}
	m.component.producers[m.name] = producer

	return producer, nil
}

type mockProducer struct {
	name        string
	component   *MockComponent
	invocations int
	messages    []core.Message
}

func (m *mockProducer) Init() {

}

func (m *mockProducer) Stop() {

}

func (m *mockProducer) Close() {

}

func (m *mockProducer) Start() {

}

func (m *mockProducer) Process(exchange core.Exchange) {
	m.messages = append(m.messages, exchange.In())
	m.invocations++
}
