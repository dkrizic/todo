module github.com/dkrizic/proto-demo/server

go 1.19

require (
	github.com/dkrizic/proto-demo/api v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.14.0
	github.com/sirupsen/logrus v1.9.0
	golang.org/x/exp v0.0.0-20221126150942-6ab00d035af9
	golang.org/x/net v0.2.0
	google.golang.org/grpc v1.51.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/dkrizic/proto-demo/api => ../api
