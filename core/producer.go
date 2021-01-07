package core

// A Producer processes Exchanges in a Guancano Route
type Producer interface {
	ConsumingService
	Processor
}
