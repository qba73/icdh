[![Go](https://github.com/qba73/icdh/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/qba73/icdh/actions/workflows/go.yml)
![GitHub](https://img.shields.io/github/license/qba73/icdh)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/qba73/icdh)
[![Go Report Card](https://goreportcard.com/badge/github.com/qba73/icdh)](https://goreportcard.com/report/github.com/qba73/icdh)
[![Go Reference](https://pkg.go.dev/badge/github.com/qba73/icdh@v0.1.0.svg)](https://pkg.go.dev/github.com/qba73/icdh@v0.1.0)

# icdh

```icdh``` is a Go client library for [NGINX Ingress Controller Deep Service Insight](https://docs.nginx.com/nginx-ingress-controller/logging-and-monitoring/service-insight/) API.

## Using the Go library

Import the library using:

```go
import "github.com/qba73/icdh"
```

## Creating a client

Create a new ```client``` object by calling ```icdh.NewClient(baseURL)```

```go
client, err := icdh.NewClient("http://localhost:9114")
if err != nil {
    // handle err
}
```

Or create a client with a specific http Client:

```go
myHTTPClient := &http.Client{}

client, err := icdh.NewClient(
    "http://localhost:9114",
    icdh.WithHTTPClient(myHTTPClient),
)
if err != nil {
    // handle error
}
```

## Retrieve statistics for host `my.service.com`

```go
stats, err := client.GetStats(ctx, "my.service.com")
if err != nil {
    // handle err
}
```

## Retrieve statistics for name (transport) `service`

```go
stats, err := client.GetTSStats(ctx, "service")
if err != nil {
    // handle err
}
```
