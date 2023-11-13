module github.com/dkrizic/todo/client

go 1.19

replace github.com/dkrizic/todo/api/todo => ../api/todo

require (
	github.com/dkrizic/todo/api/todo v0.0.0-20230209100053-e18c0151a032
	github.com/google/uuid v1.4.0
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.59.0
)

require (
	github.com/go-chi/chi/v5 v5.0.10 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.14.0 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20231030173426-d783a09b4405 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)
