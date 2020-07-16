# Primitive for your Web Server

## Add a route to your existing web server

The following will bind a new route to your server where users can
POST files to `/primitive` with the `file` form key and process them with
primitive. The response returns a raw PNG buffer with Content-Type `image/png`.

```go
import (
	"net/http"

	pr "github.com/fogleman/primitive/server/route"
)

handler := pr.PrimitiveRoute(pr.Config{
	MaxUploadMb: 10,
	FileKey: "file",
})
http.HandleFunc("/primitive", handler)
```

## Example web server

An examples web server you can start that has a `/primitive` route
bound that accepts file uploads to be processed with primitive
can be started with the following instructions.

```
$ go get github.com/fogleman/primitive/server

$ # run the web server on port :8080
$ go run $(go env GOPATH)/src/github.com/fogleman/primitive/server

$ # in another terminal window, make sure the server is running
$ curl localhost:8080
We're running!
```

## TODOS

- Make all parameters configurable (Background, Alpha, Mode, etc.)
- Dockerfile for example web server