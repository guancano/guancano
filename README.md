# Guancano

> The guanaco (Lama guanicoe) is a camelid native to South America, closely related to the llama. Its name comes from the Quechua word huanaco[2] 
> (modern spelling wanaku). Young guanacos are called chulengos.[3] Guanacos are one of two wild South American camelids, the other being the Vicuña, 
> which lives at higher elevations.`
> 
> \- [Wikipedia](https://en.wikipedia.org/wiki/Guanaco)

Guancano is an integration framework loosely based on the API implemented by 
[Apache Camel](https://camel.apache.org/). The API should be familiar
to Camel users but Guancano is not a port of Camel.

## Goals
* Create a Camel-inspired API
* Provide [Enterprise Integration Patterns](https://www.enterpriseintegrationpatterns.com/)
* Specific focus on routing, splitting, choices, and limiting

## Non-Goals
* Replicate all of Camel in Go
* Reimplement all the Camel components in Go

## Example

The following code example highlights the one way to save http requests to a
directory.

```go
package cmd

import (
    "github.com/guanaco/guancano/core"
    "github.com/guanaco/guancano/http"
    "github.com/guanaco/guancano/file"
)

func main(args []string) {
    // create a new context from the core 
    context := core.Create()

    // register components that will be used in the route
    context.Register(http.ComponentCreator)
    context.Register(file.ComponentCreator)
    
    // construct routes for the context
    context.Add(func(builder core.RouteBuilder) {
        // build a single route that takes traffic to localhost:9090/upload
        // and saves it to a file in /tmp/files with the name upload_ followed
        // by the "date" value from the exchange headers
        builder.FromS("http://localhost:9090/upload").
                ToS("file:/tmp/files/upload_{{date}}.raw")
    })

    // start the context
    context.Start()
}
```
 