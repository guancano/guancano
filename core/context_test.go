package core

import (
	"fmt"
	"testing"
)

// testApiUsability a standalone test to establish that the APi is at
// least somewhat usable consumers the standpoint of an integrator. This
// should _not_ be started or do anything. This is _not_ a functional
// test.
func TestApiUsability(t *testing.T) {
	context := Create()
	context.Add(func(builder RouteBuilder) {
		// simple string-based routeing
		builder.FromS("direct:route1").ToS("direct:route2")

		// with a processor function
		builder.FromS("direct:route1").ProcessFunction(func(exchange Exchange) {
			fmt.Printf("Exchange: %s", exchange.Id())
		}).ToS("direct:route3")

		// pattern-based routing
		builder.FromF("direct:route%d", 3).ToF("direct:%s", "12")
	})
}
