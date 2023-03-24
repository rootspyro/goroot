# GoRoot Documentation 📖

## Contents
 - [1. Server](#server)
	 - [1.1 Configuration](#configuration)
		 - [1.1.1 Global Middlewares](#global-middlewares)
		 - [1.1.2 Cors](#cors-policy)
		 - [1.1.3 Pages](#pages)
	- [1.2 Endpoints](#endpoints)
	- [1.3 Static Files](#static-files)

- [2. Handlers](#handlers)
	- [2.1 Root](#root)
		- [2.1.1 Status](#status)
		- [2.1.2 Responses](#responses)
	- [2.2 Middlewares](#middlewares)

- [3. HTML Rendering](#html-rendering)


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

You can also define dynamic routes by enclosing between `{}` the name of the parameters you want to receive from the request. Example:

```go
func main() {

	server := goroot.New(goroot.Config{})

	server.Get("/users", func(root *goroot.Root) {
		users := services.GetUsers()
		root.Json(users)
	})
	
	server.Get("/users/{userID}", func(root *goroot.Root){
		userID := root.RequestParams["userID"]
		user := services.GetUser(userID)
		root.Json(users)
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

Handlers are the functions related to the application endpoints, they manage what happens when an http request arrives at the server.

A GoRoot Handler consists of a simple function that receives as a single parameter an object of type Root.

```go
// handlers/examplefile.go
package handlers

func ExampleHandler(root *goroot.Root){
	root.OK().Send("Hello world from an example handler!")
}
```

### Root 
The Root object is the hearth of the handlers, it has the methods to send and receive data between the server and the client as:

	- HTTP Statuses
	- HTTP Requests body
	- HTTP Requests parameters and queries
	- Responses format: Plain Text, Json, Html rendering.

Some relevant functions are the different statuses and response formats that can be send it by the handler.

#### Status

The `Status()` function receives as a single parameter an integer with the HTTP code that we want to write in the response header.

```go
// handlers/examplefile.go
package handlers

func ExampleHandler(root *goroot.Root){
	// HTTP - 404 NOT FOUND
	root.Status(404)
}
```

For this purpose there is also a list of functions for the most regularly used HTTP codes.

```go
root.OK() // HTTP - 200 OK
root.Created() // HTTP - 201 Created

root.BadRequest() // HTTP - 400 Bad Request
root.Unauthorized() // HTTP - 401 Unauthorized
root.Forbidden() // HTTP - 403 Forbidden
root.NotFound() // HTTP - 404 Not Found
root.MethodNotAllowed() // HTTP - 405 Method Not Allowed

root.InternalServerError() // HTTP - 500 Internal Server Error 
root.NotImplemendted() // HTTP - 501 Not Implemented
```
Very usefully, but not enough, most of the time we will need to send a message  to the client.


#### Responses
Currently we can send three types of responses with GoRoot:
 - Plaint Text
 - Json
 - Html

We going to talk more about html in the [HTML Rendering](#html-rendering) section. Let's focus on the other two functions, `Send()` and `Json()`.

```go
server.Get("/", func(root *goroot.Root) {
	// A plain text response
	root.Send("This is the GET Method of the Route: /")
})
	
server.GetUser("/users", func(root *goroot.Root){
	// A json response
	exampleUser := User{ID: 1, Username: "rootspyro"}
	root.Json(exampleUser)
})

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
}
```
In the same way, we can combine status methods with response methods.

```go
// Simulation of a POST request
server.Post("/users", func(root *goroot.Root) {
	
	body, _ := root.Body()
	var reqUser User
	json.Unmarshal(body, &reqUser)
			
	newUser := services.NewUSer(reqUser)
	root.Created().Json(newUser) // json response with 201 Code
})
	
type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
}
```

### Middlewares
We already talked about [Global Midlewares](#global-middleware), but if we want to create a middleware exclusively for a handler or several handlers we have to use the `AddMiddleware()` function.  

Returning to the last example concerning middlewares:

```go
func main() {
	server := goroot.New(goroot.Config{})
	
	// In this case, the middleware is assigned only to MyHandler.
	server.Get("/", server.AddMiddleware(
		MyHandler,
		middlewares.Example1,
	))
	
	server.Listen()
}

func MyHandler(root goroot.Root){
		fmt.Println("My Handler")
		root.OK().Send("Hello World")
}
```

## Html Rendering
Continuing with the example seen in the [Pages](#pages) section, it is time to explain the configuration and structure that GoRoot uses to serve html templates.

First we gonna define an example code for the HTML templates.

### index.html

In this file the base template is defined as 'layout' and is the main structure for the components and content.
```html
<!-- Here we define the name of the base template as 'layout' -->
{{define "layout"}}
<!DOCTYPE html>
<html>
	<head>
		<title>{{.Title}}</title>
	</head>
	<body>
		<header>
			{{template "navbar" .}} <!-- template defined in navbar.html -->
		</header>

		<div class="main-container">
			{{template "content" .}} 
			<!-- With the "content" tag we are going to define
			the html code of our different pages -->
		</div>
		
		<footer>
			{{template "footer" .}} <!-- template defined in footer.html -->
		</footer>
	</body>
</html>
{{end}}
```
### navbar.html

In this file the navbar template is defined.
```html
{{define "navbar"}}
<div>
	<nav>
		<ul>
			<li><a href="/">Home</a></li>
		</ul>
	</nav>
</div>
{{end}}
```

### footer.html

In this file the footer template is defined.
```html
{{define "footer"}}
<div>
	<p>My footer</p>
</div>
{{end}}
```


In this example we have a main template (index.html) and a number of components requested in the main template (navbar.html and footer.html). 

These are the files that we must indicate in the configuration of the pages, since we need them to render and serve all the content.

```go
package main

import (
	"github.com/rootspyro/goroot"
	"github.com/rootspyro/goroot/pages"
)

func main() {

	pagesConfig := pages.NewPages(
		"layout", // the main template that was defined in index.html
		"./templates", // the local dir of the html templates in the code
		[]string{ // the name of the essential files without the extension .html
			"index",
			"navbar",
			"footer",
		},
	)

	server := goroot.New(goroot.Config{
		Pages: pagesConfig,
	})

	server.Listen()
}
```
Once our templates have been configured on the server, it is time to create the routes that were defined in the navbar component.

### home.html
```html
{{define "content" .}}
<div class="home-container">
	<h1>GoRoot</h1>
	<p>Golang Backend Module</p>
</div>
{{end}}
```
### about.html
```html
{{define "content" .}}
<div class="home-container">
	<h1>About</h1>
	<p>{{.Content}}</p>
</div>
{{end}}
```

In the main.go file we create the paths for the home and about pages.

### main.go
```go
// main.go
package main

import (
	"github.com/rootspyro/goroot"
	"github.com/rootspyro/goroot/pages"
)

func main() {

	pagesConfig := pages.NewPages(
		"layout", // the main template that was defined in index.html
		"./templates", // the local dir of the html templates in the code
		[]string{ // the name of the essential files without the extension .html
			"index",
			"navbar",
			"footer",
		},
	)

	server := goroot.New(goroot.Config{
		Pages: pagesConfig,
	})

	server.Get("/", func(root *goroot.Root) {
	
		// The RenderTemplate function recieves the name
		// of the file with the content and the data
		// root.RenderTemplate(filename, data)
		root.RenderTempate("home", Page{Title:"GoRoot"})
	})

	server.Get("/about", func(root *goroot.Root) {
		root.RenderTempate("about", Page{
			Title: "About",
			Content: "This is the content of the about page",
		})
	})

	server.Listen()
}

type Page struct {
	Title string
	Content string
}
```


## The End 
And that's all about the GoRoot Documentation. Thanks for reading it.

I give you a duck. 🦆
