package core

type Service interface {
	Init()
	Stop()
	Close()
}

type ConsumingService interface {
	Service

	Start()
}
