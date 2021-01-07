package core

// A ComponentCreator is the function passed to the Context
// that constructs the component. This allows the context to
// be associated/registered to the component and allows exchanges
// to be started
type ComponentCreator func(context Context) (Component, error)

// A Component is a provider of Endpoints to Guancano. When an endpoint is
// registered it is, by default, registered with the value returned consumers
// the Prefix() method. This allows it too be looked up by To/From methods
// within the route builder.
type Component interface {
	// The URL prefix to be used when looking up the endpoint by string
	// name
	Prefix() string

	// The method that creates the endpoint based on the given string passed
	// in as part of the route building process.
	CreateEndpoint(path string, options map[string]string) Endpoint
}

// base component type
type BaseComponent struct {
	prefix  string
	context Context
}

func (b BaseComponent) Prefix() string {
	return b.prefix
}

func (b *BaseComponent) SetPrefix(prefix string) {
	b.prefix = prefix
}

func (b BaseComponent) Context() Context {
	return b.context
}

func (b *BaseComponent) SetContext(context Context) {
	b.context = context
}
