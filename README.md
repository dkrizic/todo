# ToDo List

Sample project that includes the following features:

* Go programming language
* CRUD service offering a simple ToDo list
* Interface is defined via protobuf
* Protoc generates
  * gRPC Server
  * REST Gateway
  * Swagger documentation

So in summary the service offers a REST API and a gRPC API from a single definition.
The code is compiled for

* amd64
* arm64

and available as a Multiaarch Docker image under the following tags:

```
ghcr.io/dkrizic/todo-server:latest
```

## API

The directory /api contains the protobuf definition of the API

```
todo.proto
```

All other files are generated or downloaded. The generated can be updated running

```
./update
```

The generated files are commited

## Server

The server is implemented in Go and can be found in the directory /server
It can be locally run using

```
$ go run server.go
INFO[0000] Creating new server                          
INFO[0000] Serving gRPC on :8080                        
INFO[0000] Service HTTP on :8090                        
INFO[0035] Creating new todo                             id=d93f5341-071b-462b-af6e-d397aafbe206 title="Another todo"
INFO[0035] Getting all todos                             count=1
```

### Ports

The following ports are used

* 8080 - gRPC
* 8090 - REST

## Client

The client user the gRPC API to create and list todos

```
$ go run client.go 
INFO[0000] Starting app                                 
INFO[0000] Got todo                                      id=d93f5341-071b-462b-af6e-d397aafbe206 title="Another todo"
```