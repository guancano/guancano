package direct

import (
	"github.com/stretchr/testify/assert"
	"github.com/guanaco/guancano/core"
	"github.com/guanaco/guancano/mock"
	"testing"
)

func TestBasicDirect(t *testing.T) {

	context := core.Create()
	context.Register(ComponentCreator)
	component := context.Register(mock.ComponentCreator)
	mocker := component.(mock.MockComponent)

	// collector function to verify test results
	process := func(exchange core.Exchange) {
		in := exchange.In()
		if text, ok := in.(core.TextMessage); ok {
			t.Logf("got text message in exchange %s, \"%s\"", exchange.Id(), text.Text())
		} else {
			t.Logf("got message in exchange %s", exchange.Id())
		}
	}

	context.Add(func(builder core.RouteBuilder) {
		builder.FromS("mock:start1").ToS("direct:route1a")
		builder.FromS("mock:start2").ToS("direct:route1b")
		builder.FromS("direct:route1a").ToS("direct:route2a")
		builder.FromS("direct:route1b").ToS("direct:route2b")
		builder.FromS("direct:route2a").FromS("direct:route2b").ToS("direct:route3")
		builder.FromS("direct:route3").ProcessFunction(process).ToS("mock:out")
	})

	context.Start()

	message := core.NewTextMessage("hello")
	mocker.Send("mock:start1", message)
	mocker.Send("mock:start2", message)

	count, messages := mocker.ProducerStats("mock:out")
	assert.Equal(t, 2, count)
	assert.Equal(t, 2, len(messages))
}

func TestMultiPrefix(t *testing.T) {

	context := core.Create()
	context.Register(ComponentCreator)
	context.RegisterWithPrefix("alternate", ComponentCreator)
	component := context.Register(mock.ComponentCreator)
	mocker := component.(mock.MockComponent)

	// collector function to verify test results
	process := func(exchange core.Exchange) {
		t.Logf("got message in exchange %s", exchange.Id())
	}

	context.Add(func(builder core.RouteBuilder) {
		builder.FromS("mock:start").ToS("direct:route1a").ToS("direct:route1b")
		builder.FromS("direct:route1a").ToS("direct:route2a")
		builder.FromS("direct:route1b").ToS("direct:route2b")
		builder.FromS("direct:route2a").FromS("direct:route2b").ToS("direct:route3")
		builder.FromS("direct:route3").ProcessFunction(process).ToF("alternate:route%d", 1).ToF("alternate:route%s", "2")
		builder.FromS("alternate:route1").ProcessFunction(process).ToS("mock:aggregate")
		builder.FromS("alternate:route2").ProcessFunction(process)
	})

	context.Start()

	mocker.Send("mock:start", core.NewTextMessage("test"))
	counts, _ := mocker.ProducerStats("mock:aggregate")
	assert.Equal(t, 2, counts)
}
