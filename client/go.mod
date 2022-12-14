module github.com/dkrizic/todo/client

go 1.19

replace github.com/dkrizic/todo/api/todo => ../api/todo

require (
	github.com/dkrizic/todo/api/todo v0.0.0-20230110140236-29735055cd79
	github.com/google/uuid v1.3.0
	github.com/sirupsen/logrus v1.9.0
	google.golang.org/grpc v1.51.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.14.0 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
