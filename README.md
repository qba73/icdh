[![Go](https://github.com/qba73/icdh/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/qba73/icdh/actions/workflows/go.yml)

# icdh
```icdh``` is a Go client library for NGINX Ingress Controller Service Insight API.

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
