package core

func Create() Context {
	return &context{
		components: make(map[string]Component),
		routes:     make([]Route, 0),
	}
}

// Context is the container for routes that is used
// to collect and control the lifecycle of the routes.
// The Context also maintains the dependency injection
// context for providing configuration and other
// elements.
type Context interface {
	ConsumingService

	Register(creator ComponentCreator) Component
	RegisterWithPrefix(prefix string, creator ComponentCreator) Component

	Add(creator RouteCreator)
}

// context is the implementation of thc *context interface
type context struct {
	components map[string]Component
	routes     []Route
}

// Init calls each Route's Init() in turn. There is no specific
// order to the initialization.
func (c *context) Init() {
	for _, route := range c.routes {
		route.Init()
	}
}

func (c *context) Start() {
	for _, route := range c.routes {
		route.Start()
	}
}

func (c *context) Stop() {
	for _, route := range c.routes {
		route.Stop()
	}
}

func (c *context) Close() {
	for _, route := range c.routes {
		route.Close()
	}
}

func (c *context) Add(creator RouteCreator) {
	builder := &routeBuilder{
		components:          c.components,
		routeConfigurations: make([]*routeConfiguration, 0),
	}
	creator(builder)
	for idx := 0; idx < len(builder.routeConfigurations); idx++ {
		c.routes = append(c.routes, builder.routeConfigurations[idx].build())
	}
}

func (c *context) Register(creator ComponentCreator) Component {
	component, _ := creator(c)
	c.register(component.Prefix(), component)
	return component
}

func (c *context) RegisterWithPrefix(prefix string, creator ComponentCreator) Component {
	component, _ := creator(c)
	c.register(prefix, component)
	return component
}

func (c *context) register(prefix string, component Component) {
	if _, found := c.components[prefix]; found {
		// todo: log out component overwrite
		return
	}
	c.components[prefix] = component
}
