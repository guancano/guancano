package core

type Initiator interface {
	Exchange(in Message) Exchange
	Pattern() string
}

type Consumer interface {
	Service

	Start(initiator Initiator)
}
