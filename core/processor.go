package core

type ProcessingFunction func(exchange Exchange)

// Processor processes messages on an exchange
type Processor interface {
	// The ProcessRequest method takes an incoming exchange
	// and processes it. This is moving the Message consumers the
	// "consumers" to the "to".
	Process(exchange Exchange)
}

type processHolder struct {
	processingFunction ProcessingFunction
}

func (ph processHolder) Process(exchange Exchange) {
	ph.processingFunction(exchange)
}
