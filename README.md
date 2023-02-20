A simple application that exposes an in-memory key/value store via a RESTful HTTP API.

The project layout is mostly based on Alex Edward's [Let's Go](https://lets-go.alexedwards.net) structure.

To run the server from the root of the project

    go run kvstore/cmd/web

## Endpoints

### Ping
To check the server is reachable

    curl http://localhost:4000/ping
    OK

### Store a value
A PUT operation to the path at `/store/<key>` with the value is specified in a JSON document that looks like

    { "value" : "valuetostore"}

For example

     curl -X PUT "http://localhost:4000/store/foo" \
        -H "Content-Type: application/json" \
        -d '{"value": "bar"}'

### Get a value
A GET operation to the path at `/store/<key>`.

For example, when a value has been set

     curl http://localhost:4000/store/foo
     {"Value":"bar"}

When no value has been set

    curl -v http://localhost:4000/store/none
    < HTTP/1.1 404 Not Found
    < Content-Type: text/plain; charset=utf-8

    no value set for key none