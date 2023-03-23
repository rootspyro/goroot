# GoRoot

GoRoot it's a golang module or micro-framework for the backend development.

## Features

GoRoot is fully developed with golang standard libraries and supports functionalities such as:

- Responses:
    - Plain Text
    - Json
    - Html rendering
    - Static files

- Dynamic Routing
- Middlewares management
- Cors policy configuration
- Server configuration
- Http request logger
## Quick Start

Start your project with GoRoot

```bash
go mod init ${project_name}
go get -u "github.com/rootspyro/goroot"
```

Now open the main.go file in your code editor and type the following code.
```golang
// main.go
package main

import "github.com/rootspyro/goroot"

func main() {

    // Create a new server
	server := goroot.New(goroot.Config{})

    // Create 
	server.Get("/", func(root *goroot.Root) {
		root.OK().Send("Hello World")
	})
    
    // default port 3000
	server.Listen()
}
```

Finally run your main.go and you're done!

```bash
go run main.go
```

You have created a web server with GoRoot.

Follow the [Documentation](docs/README.md) to keep learning!.



## License

[MIT](https://choosealicense.com/licenses/mit/)


## Authors

- [@rootspyro](https://www.github.com/rootspyro)

