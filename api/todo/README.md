# Go API Server for todo

A simple Todo API

## Overview
This server was generated by the [openapi-generator]
(https://openapi-generator.tech) project.
By using the [OpenAPI-Spec](https://github.com/OAI/OpenAPI-Specification) from a remote server, you can easily generate a server stub.
-

To see how to make this your own, look here:

[README](https://openapi-generator.tech)

- API version: 1.0.0
- Build date: 2023-01-27T22:58:41.854806+01:00[Europe/Berlin]
For more information, please visit [https://todo.krizic.net](https://todo.krizic.net)


### Running the server
To run the server, follow these simple steps:

```
go run main.go
```

To run the server in a docker container
```
docker build --network=host -t todo .
```

Once image is built use
```
docker run --rm -it todo
```
