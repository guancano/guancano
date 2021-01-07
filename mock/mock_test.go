package mock

import (
	"github.com/stretchr/testify/assert"
	"github.com/guanaco/guancano/core"
	"testing"
)

func TestMock(t *testing.T) {

	context := core.Create()

	component := context.Register(ComponentCreator)
	mocker := component.(MockComponent)

	context.Add(func(builder core.RouteBuilder) {
		builder.FromS("mock:none").ToS("mock:test1")
	})

	context.Start()

	producer := mocker.producers["mock:test1"]
	assert.NotNil(t, producer)

	exchange := core.NewExchange()
	producer.Process(exchange)

	mocker.Send("mock:none", core.NewTextMessage("test"))

	invocations, messages := mocker.ProducerStats("mock:test1")
	assert.Equal(t, 2, invocations)
	assert.Equal(t, 2, len(messages))
}

func TestRequestOnlyMock(t *testing.T) {

	context := core.Create()

	component := context.Register(ComponentCreator)
	mocker := component.(MockComponent)

	context.Add(func(builder core.RouteBuilder) {
		builder.FromS("mock:start").RequestOnly().ToS("mock:test1").ToS("mock:test2").ToS("mock:test1").ToS("mock:test1")
	})

	context.Start()

	mocker.Send("mock:start", core.NewTextMessage("test"))

	invocations, messages := mocker.ProducerStats("mock:test1")
	assert.Equal(t, 3, invocations)
	assert.Equal(t, 3, len(messages))

	// request only would not send messages back to the start
	invocations, messages = mocker.ConsumerStats("mock:start")
	assert.Equal(t, 0, invocations)
	assert.Equal(t, 0, len(messages))
}

func TestInOutMock(t *testing.T) {

	context := core.Create()

	component := context.Register(ComponentCreator)
	mocker := component.(MockComponent)

	context.Add(func(builder core.RouteBuilder) {
		builder.FromS("mock:start").RequestReply().ToS("mock:test1")
	})

	context.Start()

	mocker.Send("mock:start", core.NewTextMessage("test"))

	invocations, messages := mocker.ProducerStats("mock:test1")
	assert.Equal(t, 1, invocations)
	assert.Equal(t, 1, len(messages))

	invocations, messages = mocker.ConsumerStats("mock:start")
	assert.Equal(t, 1, invocations)
	assert.Equal(t, 1, len(messages))
}
