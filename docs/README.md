# GoRoot Documentation

## Contents
 - [1. Server](#server)
	 - [1.1 Configuration](#configuration)
		 - [1.1.1 Global Middlewares](#global-middlewares)
		 - [1.1.2 Cors](#cors-policy)
		 - [1.1.3 Pages](#pages)
	- [1.2 Endpoints](#endpoints)
	- [1.3 Static Files](#static-files)

- [2. Handlers](#handlers)
	- [2.1 Root]()
		- [2.1.1 Status]()
		- [2.1.2 Responses]()
	- [2.2 Middlewares]()

- [3. HTML Rendering]()
	- [3.1 File structures]()


## Server

The server object it's the main interaction interface between the programmer and the module. It is used to configure the application, create rules, routes and run the server.

To create a new server, simply create an object with the `goroot.New()` function using the configuration structure as a parameter.


The next step would be to run the server, for this our new server object has the Listen() function, this function executes the server listening on the default port 3000.

```go
//main.go
package main

import "github.com/rootspyro/goroot"

func main() {
	server := goroot.New(goroot.Config{})
	server.Listen()
}
```
If you create the PORT environment variable (Example: PORT=5000) GoRoot will detect it and set the default port with the value of the PORT variable.

```bash
# .env file
PORT=5000
```
```go
func main() {
	server := goroot.New(goroot.Config{})
	server.Listen()
	// output: Server listening on 0.0.0.0:5000!
}
```


### Configuration
The function `New()` receives as parameter the `Config{}` object. In this object we can define different rules and parameters for the correct operation of our server.

#### Global Middlewares

If you have functions or middlewares that you want to be executed before any handler, you can configure a list of middlewares in the server configuration.

We are going to create two example middlewares:
```go
package middlewares

import (
	"fmt"
	"github.com/rootspyro/goroot"
)

func Test1(handler goroot.Handler) goroot.Handler {
	return func(root *goroot.Root) {
		fmt.Println("Middleware 1")
		handler(root)
	}
}

func Test2(handler goroot.Handler) goroot.Handler {
	return func(root *goroot.Root) {
		fmt.Println("Middleware 2")
		handler(root)
	}
}
```
Then, we will pass it to the configuration as a list of middlewares.

```go
func main() {
	server := goroot.New(goroot.Config{
		Middlewares: []goroot.Middleware{
			middlewares.Test1,
			middlewares.Test2,
		},
	})
	
	server.Get("/", func(root *goroot.Root) {
		fmt.Println("My Handler")
		root.OK().Send("Hello World")
	})
	server.Listen()
}
```
```bash
# Console Output for the http request to localhost:3000/

Server listening on 0.0.0.0:3000!
2023/03/23 00:18:44 GET:/ - localhost:3000
Middleware 2
Middleware 1
My Handler
```

#### Cors Policy

It is highly recommended to configure the Cors Policy on our server, because if we leave the Cors configuration empty then GoRoot will set the Allowed Origins as `"*"`.
```go
// main.go
package main

import (
	"github.com/rootspyro/goroot"
	"github.com/rootspyro/goroot/cors"
)

func main() {

	corsConfig := cors.New(cors.Config{
		Methods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		Origins: []string{"http://example.domain.com", "http://example.domain2.com"},
	})

	server := goroot.New(goroot.Config{
		Cors: corsConfig,
	})

	server.Listen()
}

```

#### Pages

If we want to render html files with GoRoot, we will need to indicate a specific configuration on our files.

To configure the html rendering we are going to use the `pages.NewPages()` function, this function requires three important parameters:

- __baseName:__ The name of the main layout template. Example: "layout".
- __tmpPath:__ The path of the .html files. Example: "./templates".
- __templates:__ The list of the required templates and components excluding the content files. Example: ["index", "navbar", "footer"].

```go
/*
Files Structure:
 - main.go
 - templates/
	 - index.html
	 - content/
		 - home.html
	 - components/
		 - navbar.html
		 - footer.html
	
*/
package main

import (
	"github.com/rootspyro/goroot"
	"github.com/rootspyro/goroot/pages"
)

func main() {

	server := goroot.New(goroot.Config{
		Pages: pages.NewPages("layout", "./templates", []string{
			"index",
			"components/navbar",
			"components/footer",
		}),
	})

	server.Get("/", func(root *goroot.Root) {
		root.RenderTempate("content/home", PageData{
			Title: "Home page",
			Content: "Hello world from the home!",
		})
	})

	server.Listen()
}

type PageData struct {
	Title string
	Content string
}

```
We will return to this example later in the section on [HTML Rendering](#html-rendering).

### Endpoints
The server object has a number of functions to create the routes of your application that get as parameters:
- __path:__ the route. Example: "/users"
- __handler:__ the handler function 

```go
func main() {

	server := goroot.New(goroot.Config{})

	server.Get("/", func(root *goroot.Root) {
		root.Send("This is the GET Method of the Route: /")
	})

	server.Post("/", func(root *goroot.Root) {
		root.Send("This is the POST Method of the Route: /")
	})

	server.Put("/", func(root *goroot.Root) {
		root.Send("This is the PUT Method of the Route: /")
	})

	server.Patch("/", func(root *goroot.Root) {
		root.Send("This is the PATCH Method of the Route: /")
	})

	server.Delete("/", func(root *goroot.Root) {
		root.Send("This is the DELETE Method of the Route: /")
	})

	server.Listen()
}
```

### Static Files

GoRoot has the `StaticFiles()` function that allows you to serve static files and resources easily. It receives two parameters:
- __path:__ the route from which the resources will be served. Example: "/static"
- __src:__ the local directory of the files to serve. Example: "./frontend/static"

```go
/*
Files Structure:
 - main.go
 - frontend/
	 - templates/
	 - static/
		 - css/
		 - js/
		 - pics/
			 - icon.png
*/
func main() {

	server := goroot.New(goroot.Config{})

	server.StaticFiles("/static", "./frontend/static")
	// The resources will be served in http://localhost:3000/static
	// Example: http://localhost:3000/static/pics/icon.png
	
	server.Listen()
}

```

## Handlers
